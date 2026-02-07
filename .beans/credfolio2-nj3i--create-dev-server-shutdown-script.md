---
# credfolio2-nj3i
title: Create dev server shutdown script
status: todo
type: task
created_at: 2026-02-07T16:29:25Z
updated_at: 2026-02-07T16:29:25Z
parent: credfolio2-ynmd
---

Create a reliable one-command script to shut down all dev server processes.

## Why

Shutting down the dev server currently requires a multi-step sequence documented in CLAUDE.md (pkill turbo, pkill go, pkill next, fuser -k ports, sleep, verify). Claude often gets this wrong or incomplete, leaving orphan processes that block the next startup. A single script eliminates this detour.

## What

Create `.claude/scripts/stop-dev.sh` that:
1. Kills turbo orchestrator: `pkill -f "turbo run dev"`
2. Kills Go backend: `pkill -f "go run cmd/server"`
3. Kills Next.js frontend: `pkill -f "next dev"`
4. Kills anything on ports 3000 and 8080: `fuser -k 8080/tcp 3000/tcp`
5. Waits briefly for processes to terminate
6. Verifies ports are free with `lsof`
7. Reports success/failure

The script should be idempotent — safe to run even if nothing is running.

## Integration

- Add to CLAUDE.md under Common Commands: `.claude/scripts/stop-dev.sh`
- Update the "Starting Dev Servers" section to reference the script
- Consider adding a permission rule: `Bash(.claude/scripts/stop-dev.sh)` in allow list

## Definition of Done
- [ ] Script created at `.claude/scripts/stop-dev.sh`
- [ ] Script handles all process types (turbo, go, next)
- [ ] Script is idempotent (no errors if nothing is running)
- [ ] Script verifies ports are free after cleanup
- [ ] CLAUDE.md updated to reference the script
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Branch pushed and PR created for human review