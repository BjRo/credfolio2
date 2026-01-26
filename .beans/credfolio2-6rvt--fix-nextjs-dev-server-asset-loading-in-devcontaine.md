---
# credfolio2-6rvt
title: Fix Next.js dev server asset loading in devcontainer
status: completed
type: bug
priority: normal
created_at: 2026-01-26T10:57:24Z
updated_at: 2026-01-26T16:31:24Z
---

After running the Next.js dev server for a while in the devcontainer, assets like JS chunks fail to load. Rebuilding the devcontainer fixes it temporarily. This is likely related to WebSocket/HMR connection issues or VSCode port forwarding instability in containerized environments.

## Root Cause Analysis

The issue occurs because:
1. Bind mounts with `consistency=delegated` can cause file system events to be missed over time
2. WebSocket connections for HMR (Hot Module Replacement) can become stale
3. VSCode port forwarding can become unstable with long-running dev servers

## Solution

Applied three-pronged fix:

### 1. Use Turbopack (package.json)
Changed dev script from `next dev` to `next dev --turbopack --hostname 0.0.0.0`
- Turbopack has better HMR stability in containerized environments
- Explicit hostname binding ensures proper network accessibility

### 2. Enable file watching polling (devcontainer.json)
Added environment variables:
- `WATCHPACK_POLLING=true` - Forces webpack to use polling instead of native file events
- `CHOKIDAR_USEPOLLING=true` - Forces chokidar (used by many tools) to use polling

Polling is more reliable than inotify for bind mounts between host and container.

### 3. Configure port attributes (devcontainer.json)
Added `portsAttributes` for ports 3000 and 8080 with labels and `onAutoForward: notify`

### 4. Next.js config (next.config.ts)
Added:
- `devIndicators: false` - Reduces dev server overhead
- `webpackMemoryOptimizations: true` - Reduces memory pressure

### 5. Auto-seed demo user on server startup (cmd/server/main.go)
Added `ensureDemoUser()` function that runs at server startup to ensure the demo user exists.
This prevents the recurring "user not found" error that happens when migrations are rolled back
or the database is reset.