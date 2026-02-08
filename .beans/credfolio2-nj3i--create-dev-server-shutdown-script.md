---
# credfolio2-nj3i
title: Create dev server shutdown script
status: in-progress
type: task
priority: normal
created_at: 2026-02-07T16:29:25Z
updated_at: 2026-02-08T10:35:04Z
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
- Update `.claude/agents/qa.md` to reference the script instead of inline commands
- Update `.claude/rules/dev-servers.md` to reference the script
- Consider adding a permission rule: `Bash(.claude/scripts/stop-dev.sh)` in allow list

## Definition of Done
- [x] Script created at `.claude/scripts/stop-dev.sh`
- [x] Script handles all process types (turbo, go, next)
- [x] Script is idempotent (no errors if nothing is running)
- [x] Script verifies ports are free after cleanup
- [x] CLAUDE.md updated to reference the script
- [x] `.claude/rules/dev-servers.md` updated to reference the script
- [x] `.claude/agents/qa.md` updated to reference the script
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Branch pushed and PR created for human review

## Implementation Plan

### Approach

Create a bash script at `.claude/scripts/stop-dev.sh` that provides a safe, one-command way to stop all dev server processes. The script will follow the same patterns as the existing `reset-db.sh` script (color-coded output, clear status messages, idempotent execution). After creating the script, update all documentation and agent instructions to reference it instead of inline command sequences.

### Files to Create/Modify

- `.claude/scripts/stop-dev.sh` (create) — Main shutdown script with color output and verification
- `/workspace/CLAUDE.md` (modify) — Add stop-dev.sh to Common Commands section
- `/workspace/.claude/rules/dev-servers.md` (modify) — Replace inline commands with script reference
- `/workspace/.claude/agents/qa.md` (modify) — Replace inline shutdown commands with script reference

### Steps

1. **Create the stop-dev.sh script at `.claude/scripts/stop-dev.sh`**
   - Add shebang and header comment documenting usage
   - Define color constants matching `reset-db.sh` pattern (RED, GREEN, YELLOW, NC)
   - Implement main shutdown logic:
     - Step 1: Kill turbo orchestrator with `pkill -f "turbo run dev" 2>/dev/null || true`
     - Step 2: Kill Go backend with `pkill -f "go run cmd/server" 2>/dev/null || true`
     - Step 3: Kill Next.js frontend with `pkill -f "next dev" 2>/dev/null || true`
     - Step 4: Kill anything on ports 8080 and 3000 with `fuser -k 8080/tcp 3000/tcp 2>/dev/null || true`
     - Step 5: Wait 2 seconds for graceful shutdown with `sleep 2`
     - Step 6: Verify ports are free with `lsof -i :8080 -i :3000 2>/dev/null`
   - Report success or failure with color-coded messages
   - Use `|| true` pattern throughout to ensure idempotency (no errors if processes don't exist)
   - Add clear output showing progress for each step (e.g., "[1/6] Stopping Turbo orchestrator...")
   - Make executable: `chmod +x .claude/scripts/stop-dev.sh`

2. **Update CLAUDE.md Common Commands section**
   - Locate the "Common Commands" section (around line 88-111)
   - Add a new "Dev server control" subsection before or after "Cleanup"
   - Add entry: `.claude/scripts/stop-dev.sh  # Stop all dev servers`
   - Keep existing cleanup commands (pnpm clean, rm -rf .turbo) separate since they serve different purposes

3. **Update .claude/rules/dev-servers.md**
   - Locate the "Starting Dev Servers" section
   - Replace the inline multi-line pkill/fuser/lsof sequence (lines 10-18) with a reference to the script
   - Update the bash code block to show:
     ```bash
     # Stop any running dev servers
     .claude/scripts/stop-dev.sh
     
     # Clear Turbopack cache if Next.js hangs (connects but never responds)
     rm -rf src/frontend/.next
     
     # Then start
     pnpm dev
     ```
   - Keep the "Why this approach" and "Common issues" sections but reference the script implementation
   - Update the explanation text to mention the script handles the three-pronged kill approach

4. **Update .claude/agents/qa.md**
   - Locate the "Ensure Dev Servers Are Running" section (around lines 15-38)
   - Replace the inline pkill/fuser sequence (lines 30-38) with a call to the script
   - Update the bash code block to:
     ```bash
     # Stop any running dev servers
     .claude/scripts/stop-dev.sh
     
     # Start fresh
     pnpm dev &
     sleep 5
     ```
   - Keep the port checking logic (lsof check) before the restart decision

5. **Test the script manually**
   - Run with no servers running: `.claude/scripts/stop-dev.sh` (should show "no processes found" or similar, exit 0)
   - Start dev servers with `pnpm dev`, then run: `.claude/scripts/stop-dev.sh` (should cleanly stop all)
   - Verify ports are actually free after running: `lsof -i :8080 -i :3000` (should show nothing)
   - Test idempotency by running twice in a row (both should succeed)

6. **Run standard completion checks**
   - Run `pnpm lint` and fix any errors
   - Run `pnpm test` and ensure all tests pass
   - Verify all checklist items are complete

### Testing Strategy

- **Manual testing**: Start dev servers, run the script, verify ports are free
- **Idempotency testing**: Run script multiple times in a row without errors
- **Edge case testing**: Run script when no servers are running (should not error)
- **Integration testing**: Verify the script works correctly when called from agent instructions
- **Lint/test**: Standard `pnpm lint` and `pnpm test` checks (no code changes, so these should pass trivially)

### Notes

- The script should use `2>/dev/null || true` pattern throughout to suppress errors when processes don't exist
- Follow the same visual style as `reset-db.sh` for consistency (color output, step numbering, clear status messages)
- No need for confirmation prompts — stopping dev servers is safe and expected
- The script should report which processes were found and killed vs. which were already stopped
- Consider adding an optional `--clear-cache` flag in the future to also clear `.turbo/` and `src/frontend/.next/`, but keep it simple for now