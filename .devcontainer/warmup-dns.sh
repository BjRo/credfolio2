#!/bin/bash
set -euo pipefail

# Warmup script for dynamic DNS domains
# This script pre-resolves domains configured for dynamic DNS resolution,
# causing dnsmasq to add their IPs to the firewall whitelist.
#
# Run this after dnsmasq starts to ensure immediate connectivity.

DYNAMIC_DOMAINS_FILE="${1:-/etc/dnsmasq.d/dynamic-domains.conf}"

if [ ! -f "$DYNAMIC_DOMAINS_FILE" ]; then
    echo "Dynamic domains file not found: $DYNAMIC_DOMAINS_FILE"
    exit 1
fi

echo "Warming up DNS cache for dynamic domains..."

# Read domains from config file, skip comments and empty lines
while IFS= read -r domain || [ -n "$domain" ]; do
    # Skip empty lines and comments
    [[ -z "$domain" || "$domain" =~ ^[[:space:]]*# ]] && continue

    # Trim whitespace
    domain=$(echo "$domain" | xargs)

    if [ -n "$domain" ]; then
        echo "Resolving $domain..."
        # Use dig to trigger DNS resolution through dnsmasq
        # This causes dnsmasq to add the IPs to the ipset
        if dig +short "$domain" @127.0.0.1 >/dev/null 2>&1; then
            echo "  ✓ $domain resolved successfully"
        else
            echo "  ⚠ $domain resolution failed (may work later)"
        fi
    fi
done < "$DYNAMIC_DOMAINS_FILE"

echo "DNS warmup complete"
