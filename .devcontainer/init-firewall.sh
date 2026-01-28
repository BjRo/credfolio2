#!/bin/bash
set -euo pipefail  # Exit on error, undefined vars, and pipeline failures
IFS=$'\n\t'       # Stricter word splitting

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DYNAMIC_DOMAINS_FILE="/etc/dnsmasq.d/dynamic-domains.conf"
DNSMASQ_CONF="/etc/dnsmasq.d/firewall.conf"
UPSTREAM_DNS_CACHE="/etc/dnsmasq.d/upstream-dns.cache"

# 1. Extract Docker DNS info BEFORE any flushing
# Capture Docker DNS chains, their rules, AND the jump rules from OUTPUT/POSTROUTING chains
DOCKER_DNS_RULES=$(iptables-save -t nat | grep -E "(127\.0\.0\.11|DOCKER_OUTPUT|DOCKER_POSTROUTING)" || true)

# Get upstream DNS servers - check cache first (from previous run), then resolv.conf
# Use || true to handle case where resolv.conf only has 127.0.0.1 (from previous run)
if [ -f "$UPSTREAM_DNS_CACHE" ]; then
    echo "Restoring upstream DNS from cache..."
    UPSTREAM_DNS=$(cat "$UPSTREAM_DNS_CACHE")
else
    # Filter out localhost DNS (127.0.0.1 and 127.0.0.11) to find real upstream servers
    # Note: 127.0.0.11 is Docker's embedded DNS, not a real upstream
    UPSTREAM_DNS=$(grep -E "^nameserver" /etc/resolv.conf | grep -vE "127\.0\.0\.(1|11)$" | awk '{print $2}' | head -3 || true)
fi

# Check if Docker's embedded DNS (127.0.0.11) is present - we'll need to restore its NAT rules
DOCKER_DNS_PRESENT=$(grep -E "^nameserver 127\.0\.0\.11$" /etc/resolv.conf || true)

if [ -z "$UPSTREAM_DNS" ]; then
    # Try Docker's internal DNS gateway first - Docker Desktop typically uses 192.168.65.1
    HOST_GW=$(ip route | grep default | awk '{print $3}')
    if [ -n "$HOST_GW" ]; then
        echo "Using host gateway as DNS fallback: $HOST_GW"
        UPSTREAM_DNS="$HOST_GW"
    else
        # Absolute fallback
        UPSTREAM_DNS="8.8.8.8
8.8.4.4"
    fi
fi

# Ensure DNS is working before network operations
# This handles container restart where resolv.conf points to 127.0.0.1 but dnsmasq isn't running
echo "Ensuring DNS is available..."
pkill dnsmasq 2>/dev/null || true

# If Docker's embedded DNS was present, keep using it (along with upstream fallback)
# This ensures DNS works during iptables flush/restore window via Docker's NAT rules
if [ -n "$DOCKER_DNS_PRESENT" ]; then
    echo "Docker embedded DNS detected, preserving 127.0.0.11"
    {
        echo "# Temporary DNS configuration for firewall setup"
        echo "nameserver 127.0.0.11"
        for dns in $UPSTREAM_DNS; do
            echo "nameserver $dns"
        done
    } > /etc/resolv.conf
else
    {
        echo "# Temporary DNS configuration for firewall setup"
        for dns in $UPSTREAM_DNS; do
            echo "nameserver $dns"
        done
    } > /etc/resolv.conf
fi
echo "Using DNS: $(grep nameserver /etc/resolv.conf | awk '{print $2}' | tr '\n' ' ')"

# Save upstream DNS for future restarts (before we modify resolv.conf to use dnsmasq)
mkdir -p "$(dirname "$UPSTREAM_DNS_CACHE")"
echo "$UPSTREAM_DNS" > "$UPSTREAM_DNS_CACHE"

# Flush existing rules and delete existing ipsets
iptables -F
iptables -X
iptables -t nat -F
iptables -t nat -X
iptables -t mangle -F
iptables -t mangle -X
ipset destroy allowed-domains 2>/dev/null || true

# 2. Selectively restore ONLY internal Docker DNS resolution
if [ -n "$DOCKER_DNS_RULES" ]; then
    echo "Restoring Docker DNS rules..."
    iptables -t nat -N DOCKER_OUTPUT 2>/dev/null || true
    iptables -t nat -N DOCKER_POSTROUTING 2>/dev/null || true
    # Restore rules one by one, ignoring errors for duplicates or invalid rules
    # Filter out chain definitions (lines starting with :) - only process actual rules (-A ...)
    while IFS= read -r rule; do
        [ -z "$rule" ] && continue
        # Skip chain definitions (e.g., ":DOCKER_OUTPUT - [0:0]")
        [[ "$rule" =~ ^: ]] && continue
        echo "  Restoring: $rule"
        # Use subshell with default IFS for proper word splitting
        # (script uses IFS=$'\n\t' globally which breaks iptables arg parsing)
        # shellcheck disable=SC2086
        if ! (IFS=' '; iptables -t nat $rule) 2>&1; then
            echo "  (failed with error above)"
        fi
    done <<< "$DOCKER_DNS_RULES"
    echo "Docker DNS rules restored"
else
    echo "No Docker DNS rules to restore"
fi

# First allow DNS and localhost before any restrictions
# Allow outbound DNS
iptables -A OUTPUT -p udp --dport 53 -j ACCEPT
# Allow inbound DNS responses
iptables -A INPUT -p udp --sport 53 -j ACCEPT
# Allow outbound SSH
iptables -A OUTPUT -p tcp --dport 22 -j ACCEPT
# Allow inbound SSH responses
iptables -A INPUT -p tcp --sport 22 -m state --state ESTABLISHED -j ACCEPT
# Allow localhost
iptables -A INPUT -i lo -j ACCEPT
iptables -A OUTPUT -o lo -j ACCEPT

# Create ipset with CIDR support
ipset create allowed-domains hash:net

# Load dynamic domains list (domains that use dnsmasq for resolution)
DYNAMIC_DOMAINS=()
if [ -f "$SCRIPT_DIR/dynamic-domains.conf" ]; then
    while IFS= read -r line || [ -n "$line" ]; do
        [[ -z "$line" || "$line" =~ ^[[:space:]]*# ]] && continue
        line=$(echo "$line" | xargs)
        [ -n "$line" ] && DYNAMIC_DOMAINS+=("$line")
    done < "$SCRIPT_DIR/dynamic-domains.conf"
fi

# Function to check if domain is in dynamic list
is_dynamic_domain() {
    local domain="$1"
    # Handle empty array case
    if [ ${#DYNAMIC_DOMAINS[@]} -eq 0 ]; then
        return 1
    fi
    for d in "${DYNAMIC_DOMAINS[@]}"; do
        if [ "$d" = "$domain" ]; then
            return 0
        fi
    done
    return 1
}

# Verify DNS is working before proceeding
echo "Verifying DNS resolution..."
for attempt in 1 2 3; do
    if dig +short api.github.com >/dev/null 2>&1; then
        echo "DNS resolution working"
        break
    fi
    if [ "$attempt" -eq 3 ]; then
        echo "ERROR: DNS resolution not working after 3 attempts"
        echo "resolv.conf contents:"
        cat /etc/resolv.conf
        echo "Docker DNS rules:"
        iptables -t nat -L -n 2>/dev/null || true
        exit 1
    fi
    echo "DNS not ready, waiting... (attempt $attempt/3)"
    sleep 2
done

# Fetch GitHub meta information and aggregate + add their IP ranges
echo "Fetching GitHub IP ranges..."
gh_ranges=""
for attempt in 1 2 3; do
    gh_ranges=$(curl -s --connect-timeout 10 https://api.github.com/meta) && break
    echo "Curl attempt $attempt failed, retrying..."
    sleep 2
done
if [ -z "$gh_ranges" ]; then
    echo "ERROR: Failed to fetch GitHub IP ranges after 3 attempts"
    exit 1
fi

if ! echo "$gh_ranges" | jq -e '.web and .api and .git' >/dev/null; then
    echo "ERROR: GitHub API response missing required fields"
    exit 1
fi

echo "Processing GitHub IPs..."
while read -r cidr; do
    if [[ ! "$cidr" =~ ^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}/[0-9]{1,2}$ ]]; then
        echo "ERROR: Invalid CIDR range from GitHub meta: $cidr"
        exit 1
    fi
    echo "Adding GitHub range $cidr"
    ipset add allowed-domains "$cidr"
done < <(echo "$gh_ranges" | jq -r '(.web + .api + .git)[]' | aggregate -q)

# Static domains - resolved once at startup (excluding dynamic domains)
STATIC_DOMAINS=(
    "registry.npmjs.org"
    "api.anthropic.com"
    "api.openai.com"
    "sentry.io"
    "statsig.anthropic.com"
    "statsig.com"
    "marketplace.visualstudio.com"
    "vscode.blob.core.windows.net"
    "update.code.visualstudio.com"
)

# Resolve and add static domains
for domain in "${STATIC_DOMAINS[@]}"; do
    # Skip if it's a dynamic domain
    if is_dynamic_domain "$domain"; then
        echo "Skipping $domain (handled by dnsmasq)"
        continue
    fi

    echo "Resolving $domain..."
    ips=$(dig +noall +answer A "$domain" | awk '$4 == "A" {print $5}')
    if [ -z "$ips" ]; then
        echo "ERROR: Failed to resolve $domain"
        exit 1
    fi

    while read -r ip; do
        if [[ ! "$ip" =~ ^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$ ]]; then
            echo "ERROR: Invalid IP from DNS for $domain: $ip"
            exit 1
        fi
        echo "Adding $ip for $domain"
        ipset add allowed-domains "$ip"
    done < <(echo "$ips")
done

# Get host IP from default route
HOST_IP=$(ip route | grep default | cut -d" " -f3)
if [ -z "$HOST_IP" ]; then
    echo "ERROR: Failed to detect host IP"
    exit 1
fi

HOST_NETWORK=$(echo "$HOST_IP" | sed "s/\.[0-9]*$/.0\/24/")
echo "Host network detected as: $HOST_NETWORK"

# Set up remaining iptables rules
iptables -A INPUT -s "$HOST_NETWORK" -j ACCEPT
iptables -A OUTPUT -d "$HOST_NETWORK" -j ACCEPT

# Set default policies to DROP first
iptables -P INPUT DROP
iptables -P FORWARD DROP
iptables -P OUTPUT DROP

# First allow established connections for already approved traffic
iptables -A INPUT -m state --state ESTABLISHED,RELATED -j ACCEPT
iptables -A OUTPUT -m state --state ESTABLISHED,RELATED -j ACCEPT

# Then allow only specific outbound traffic to allowed domains
iptables -A OUTPUT -m set --match-set allowed-domains dst -j ACCEPT

# Explicitly REJECT all other outbound traffic for immediate feedback
iptables -A OUTPUT -j REJECT --reject-with icmp-admin-prohibited

echo "Firewall configuration complete"

# --- Configure and start dnsmasq for dynamic domain resolution ---
echo ""
echo "=== Configuring dnsmasq for dynamic DNS resolution ==="

# Create dnsmasq config directory
mkdir -p /etc/dnsmasq.d

# Copy dynamic domains file
cp "$SCRIPT_DIR/dynamic-domains.conf" "$DYNAMIC_DOMAINS_FILE"

# Generate dnsmasq configuration
cat > "$DNSMASQ_CONF" << 'DNSMASQ_BASE'
# dnsmasq configuration for dynamic firewall whitelisting
# Auto-generated by init-firewall.sh

# Listen only on localhost
listen-address=127.0.0.1
bind-interfaces

# Don't read /etc/resolv.conf (we set upstream servers explicitly)
no-resolv

# Don't poll for changes to /etc/resolv.conf
no-poll

# Cache size (number of DNS entries to cache)
cache-size=1000

DNSMASQ_BASE

# Add upstream DNS servers
echo "# Upstream DNS servers" >> "$DNSMASQ_CONF"
# If Docker's embedded DNS was detected, use it as the primary upstream
# (it's more reliable than the host gateway on Docker Desktop)
if [ -n "$DOCKER_DNS_PRESENT" ]; then
    echo "server=127.0.0.11" >> "$DNSMASQ_CONF"
fi
for dns in $UPSTREAM_DNS; do
    echo "server=$dns" >> "$DNSMASQ_CONF"
done

# Add ipset rules for dynamic domains
echo "" >> "$DNSMASQ_CONF"
echo "# Dynamic domain ipset rules" >> "$DNSMASQ_CONF"
echo "# These domains will have their resolved IPs added to the allowed-domains ipset" >> "$DNSMASQ_CONF"
if [ ${#DYNAMIC_DOMAINS[@]} -gt 0 ]; then
    for domain in "${DYNAMIC_DOMAINS[@]}"; do
        echo "ipset=/$domain/allowed-domains" >> "$DNSMASQ_CONF"
        echo "Added dnsmasq ipset rule for: $domain"
    done
else
    echo "No dynamic domains configured"
fi

# Stop any existing dnsmasq
pkill dnsmasq 2>/dev/null || true
sleep 1

# Start dnsmasq
echo "Starting dnsmasq..."
dnsmasq --conf-file="$DNSMASQ_CONF" --pid-file=/var/run/dnsmasq.pid

# Verify dnsmasq is running
if pgrep dnsmasq >/dev/null; then
    echo "dnsmasq started successfully"
else
    echo "ERROR: Failed to start dnsmasq"
    exit 1
fi

# Update /etc/resolv.conf to use local dnsmasq
# Keep a backup of original nameservers as comments
echo "Configuring system to use local dnsmasq..."
{
    echo "# Local dnsmasq resolver for dynamic firewall whitelisting"
    echo "nameserver 127.0.0.1"
    echo ""
    echo "# Original upstream servers (used by dnsmasq):"
    for dns in $UPSTREAM_DNS; do
        echo "# nameserver $dns"
    done
} > /etc/resolv.conf

# Run warmup script to pre-resolve dynamic domains
echo ""
echo "=== Running DNS warmup ==="
if [ -f "$SCRIPT_DIR/warmup-dns.sh" ]; then
    chmod +x "$SCRIPT_DIR/warmup-dns.sh"
    "$SCRIPT_DIR/warmup-dns.sh" "$DYNAMIC_DOMAINS_FILE"
else
    echo "Warmup script not found, skipping"
fi

echo ""
echo "=== Verifying firewall configuration ==="

# Verify blocked access
if curl --connect-timeout 5 https://example.com >/dev/null 2>&1; then
    echo "ERROR: Firewall verification failed - was able to reach https://example.com"
    exit 1
else
    echo "✓ Blocked: example.com (as expected)"
fi

# Verify GitHub API access
if ! curl --connect-timeout 5 https://api.github.com/zen >/dev/null 2>&1; then
    echo "ERROR: Firewall verification failed - unable to reach https://api.github.com"
    exit 1
else
    echo "✓ Allowed: api.github.com"
fi

# Verify Go proxy access (tests dnsmasq integration)
if ! curl --connect-timeout 5 https://proxy.golang.org/cached-only/ >/dev/null 2>&1; then
    echo "WARNING: Unable to reach proxy.golang.org - Go dependencies may not work"
else
    echo "✓ Allowed: proxy.golang.org (via dnsmasq)"
fi

echo ""
echo "Firewall and dnsmasq configuration complete"
