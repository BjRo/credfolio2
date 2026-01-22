---
# credfolio2-nfjx
title: Fix Go dependency installation without container rebuild
status: in-progress
type: bug
created_at: 2026-01-22T08:08:30Z
updated_at: 2026-01-22T09:20:00Z
---

## Problem
Go domains (proxy.golang.org, sum.golang.org, storage.googleapis.com) use dynamic IPs that change frequently. The current firewall script resolves IPs only at container creation, causing go get to fail when IPs change.

## Solution
Use dnsmasq with ipset integration to dynamically add IPs to the whitelist as DNS queries happen.

## Checklist
- [x] Install dnsmasq in Dockerfile
- [x] Create dynamic domains config file (extensible for future use)
- [x] Create dnsmasq configuration
- [x] Modify init-firewall.sh to start dnsmasq and configure DNS
- [x] Create warmup script to pre-resolve dynamic domains
- [x] Document the change in decisions/
- [x] Fix container restart failures (DNS caching for postStartCommand re-runs)
- [x] Fix VS Code devcontainer startup failure (exit code 6)
  - Root cause 1: Docker DNS NAT rules not fully captured (missing jump rules)
  - Root cause 2: IFS=$'\n\t' broke iptables arg parsing in restore loop
  - Root cause 3: dnsmasq upstream DNS was host gateway (no DNS server), not Docker's DNS
  - Fixed grep pattern to capture DOCKER_OUTPUT and DOCKER_POSTROUTING chains
  - Added detection and preservation of Docker's embedded DNS (127.0.0.11)
  - Fixed iptables restore to use subshell with reset IFS for proper word splitting
  - Added 127.0.0.11 as dnsmasq upstream when Docker DNS is present
  - Added DNS verification with retry logic before GitHub API call
  - Added retry logic for curl command
- [ ] Test the solution in VS Code devcontainer