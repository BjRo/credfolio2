# Devcontainer Network Isolation and Permission Model

**Date**: 2026-01-20
**Bean**: credfolio2-pl9u

## Context

This project uses Claude Code as an AI-assisted development tool. When running AI agents with broad file system and command execution capabilities, there's a tension between productivity (allowing the agent to work autonomously) and security (limiting potential damage from mistakes or unexpected behavior).

We needed a way to:
1. Allow Claude Code to work efficiently during development
2. Limit network access to prevent unintended external communication
3. Maintain human oversight when running outside the controlled environment

## Decision

We implemented a two-layer security model:

### 1. Devcontainer with Network Firewall

The devcontainer ([.devcontainer/init-firewall.sh](.devcontainer/init-firewall.sh)) uses iptables to restrict outbound network access to an explicit allowlist:

- **GitHub** (api, web, git) - for version control and gh CLI
- **registry.npmjs.org** - for npm/pnpm packages
- **proxy.golang.org** - for Go modules
- **api.anthropic.com** - for Claude API
- **VS Code marketplace domains** - for extension updates
- **Sentry/Statsig** - for telemetry

All other outbound connections are blocked. The devcontainer has `dangerouslySkipPermissions: true` set in VS Code settings, allowing Claude Code to operate without per-action prompts.

### 2. Permission-Prompted Mode Outside Devcontainer

The shared `.claude/settings.json` has an empty `allow` array:

```json
{
  "permissions": {
    "allow": []
  }
}
```

This means when running Claude Code outside the devcontainer (e.g., on the host machine), every action requires explicit user approval.

## Reasoning

**Why network isolation instead of permission prompts in devcontainer?**
- Network isolation is a stronger boundary - even approved commands can't exfiltrate data
- Reduces prompt fatigue during intensive development sessions
- The allowlist approach makes the security model auditable and explicit

**Why require prompts outside devcontainer?**
- The host machine has access to credentials, other projects, and broader network
- Forces conscious decision-making when not in the sandboxed environment
- Prevents accidental use of autonomous mode in sensitive contexts

**Why not use a VPN or proxy instead?**
- iptables is simpler, doesn't require external infrastructure
- Works offline once container is built
- Easy to audit and modify the allowlist

## Consequences

1. **New Go dependencies** require adding their module hosts to the firewall (we just added `proxy.golang.org`)
2. **Developers must rebuild** the devcontainer when firewall rules change
3. **External API integrations** will need their domains added to the allowlist
4. **Running outside devcontainer** is intentionally slower due to permission prompts
5. **Container startup** may fail if allowed domains can't be resolved (DNS dependency)
