# dnsmasq for Dynamic Firewall Whitelisting

**Date**: 2026-01-22
**Bean**: credfolio2-nfjx

## Context

The devcontainer uses a firewall (iptables + ipset) to sandbox network access, allowing only specific domains. The original implementation resolved domain IPs once at container creation and added them to an ipset whitelist.

This approach worked for domains with stable IPs (GitHub, npm) but caused problems with Go dependencies. Go's module proxy (`proxy.golang.org`, `sum.golang.org`) and Google Cloud Storage (`storage.googleapis.com`) use dynamic IP pools that change frequently. When IPs changed after container creation, `go get` commands would fail, forcing container rebuilds.

## Decision

Implemented dynamic DNS-based firewall whitelisting using dnsmasq with ipset integration:

1. **dnsmasq** runs as a local DNS resolver on 127.0.0.1
2. Domains configured in `dynamic-domains.conf` trigger automatic ipset updates when resolved
3. When any process (like `go get`) resolves a dynamic domain, dnsmasq:
   - Resolves the DNS query via upstream servers
   - Automatically adds the returned IP to the `allowed-domains` ipset
   - Returns the IP to the requesting process
4. A warmup script pre-resolves dynamic domains at container start

### Files Changed

- `.devcontainer/Dockerfile`: Added dnsmasq package
- `.devcontainer/init-firewall.sh`: Integrated dnsmasq setup and configuration
- `.devcontainer/dynamic-domains.conf`: Configurable list of domains for dynamic resolution
- `.devcontainer/warmup-dns.sh`: Pre-resolves dynamic domains at startup

## Reasoning

### Why dnsmasq?

- **Mature and lightweight**: dnsmasq is a well-established DNS/DHCP server with native ipset support
- **No additional infrastructure**: Runs locally, no external proxy or service needed
- **Transparent**: Applications don't need modification; DNS resolution "just works"
- **Extensible**: Adding new dynamic domains is a one-line config change

### Alternatives Considered

1. **Transparent proxy (Squid)**: More complex, higher resource usage, overkill for this use case
2. **Periodic IP refresh**: Would require background jobs, potential gaps in connectivity
3. **Whitelist entire Google IP ranges**: Too broad, reduces sandbox effectiveness
4. **nftables with DNS**: Modern but less portable, dnsmasq is more widely available

## Consequences

### Positive

- Go dependencies install without container rebuilds
- Dynamic domains can be added via config file without code changes
- Sandbox security maintained (only explicitly listed domains are allowed)
- Minimal performance impact (dnsmasq is fast, adds ~1ms to DNS queries)

### Negative

- Additional dependency (dnsmasq) in the container
- Slight increase in container image size (~2MB)
- dnsmasq must be running for dynamic domain access (handled by init script)

### Future Considerations

- If other domains exhibit similar IP instability, add them to `dynamic-domains.conf`
- The same pattern could be applied to npm or other package registries if needed
