---
# credfolio2-fwef
title: Enable auto-memory for cross-session learning
status: todo
type: task
created_at: 2026-02-07T16:00:40Z
updated_at: 2026-02-07T16:00:40Z
parent: credfolio2-ynmd
---

Enable Claude Code's auto-memory feature so Claude builds up project knowledge across sessions.

## Why

Auto-memory lets Claude write notes for itself as it works — project patterns, debugging insights, architecture decisions, conventions discovered during sessions. These persist and load automatically in future sessions, reducing repeated exploration and building institutional knowledge over time.

This complements CLAUDE.md (what we write for Claude) with what Claude learns on its own.

## What

- Set `CLAUDE_CODE_DISABLE_AUTO_MEMORY=0` in the devcontainer environment
- Decide where to configure it: `devcontainer.json` environment variables or Dockerfile
- Verify auto-memory activates on session start (check `~/.claude/projects/<project>/memory/`)
- Optionally seed initial memory by telling Claude to remember key patterns

## How auto-memory works

- Notes saved to `~/.claude/projects/<project>/memory/`
- `MEMORY.md` entrypoint (first 200 lines loaded at session start)
- Topic files (e.g., `debugging.md`, `api-conventions.md`) loaded on-demand
- Claude curates MEMORY.md as a concise index, moves details to topic files
- Memory is user-local (not committed to git)

## Definition of Done
- [ ] `CLAUDE_CODE_DISABLE_AUTO_MEMORY=0` configured in devcontainer
- [ ] Verified auto-memory directory is created on session start
- [ ] Documented in devcontainer README that auto-memory is enabled
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Branch pushed and PR created for human review