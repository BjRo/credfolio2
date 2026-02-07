---
# credfolio2-fwef
title: Enable auto-memory for cross-session learning
status: in-progress
type: task
priority: normal
created_at: 2026-02-07T16:00:40Z
updated_at: 2026-02-07T20:40:43Z
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
- [x] `CLAUDE_CODE_DISABLE_AUTO_MEMORY=0` configured in devcontainer
- [x] Verified auto-memory directory is created on session start
- [x] Documented in devcontainer README that auto-memory is enabled
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Branch pushed and PR created for human review

## Implementation Plan

### Context Discovered During Exploration

Auto-memory is **already partially working** in the current environment:
- `~/.claude/projects/-workspace/memory/MEMORY.md` exists and contains useful project knowledge (testing patterns, lint rules, GraphQL upload patterns, codegen, workflow notes)
- `CLAUDE_CODE_DISABLE_AUTO_MEMORY` is **not set** in the environment, which means Claude Code is using its default behavior
- The `~/.claude` directory is volume-mounted (`source=claude-code-config-${devcontainerId},target=/home/node/.claude,type=volume`), so memory already persists across container rebuilds

The bean's request to explicitly set `CLAUDE_CODE_DISABLE_AUTO_MEMORY=0` is a defensive measure -- it makes the intent explicit so future Dockerfile changes or Claude Code version upgrades that might change defaults won't silently disable the feature.

### Approach

Add `CLAUDE_CODE_DISABLE_AUTO_MEMORY=0` to the `containerEnv` section of `devcontainer.json`. This is the right location because:
1. It keeps all environment configuration in one place alongside existing vars like `NODE_OPTIONS`, `CLAUDE_CONFIG_DIR`, etc.
2. The `Dockerfile` already has `ENV` statements for build-time concerns (PATH, GOPATH, etc.); runtime Claude Code behavior belongs in `devcontainer.json`
3. The existing `containerEnv` section is where the `CLAUDE_CONFIG_DIR` env var is already set, creating a natural grouping of Claude-related settings

Then update the devcontainer README to document auto-memory alongside the existing Claude Code permissions documentation.

### Files to Create/Modify

- `/workspace/.devcontainer/devcontainer.json` -- Add `CLAUDE_CODE_DISABLE_AUTO_MEMORY` to `containerEnv`
- `/workspace/.devcontainer/README.md` -- Add auto-memory documentation section

### Steps

1. **Add env var to `devcontainer.json`**
   - Open `/workspace/.devcontainer/devcontainer.json`
   - Add `"CLAUDE_CODE_DISABLE_AUTO_MEMORY": "0"` to the `containerEnv` object (line 71 area, after the existing Claude-related `CLAUDE_CONFIG_DIR` entry)
   - The resulting `containerEnv` section should look like:
     ```json
     "containerEnv": {
       "NODE_OPTIONS": "--max-old-space-size=4096",
       "CLAUDE_CONFIG_DIR": "/home/node/.claude",
       "CLAUDE_CODE_DISABLE_AUTO_MEMORY": "0",
       "POWERLEVEL9K_DISABLE_GITSTATUS": "true",
       "WATCHPACK_POLLING": "true",
       "CHOKIDAR_USEPOLLING": "true"
     }
     ```

2. **Document auto-memory in devcontainer README**
   - Open `/workspace/.devcontainer/README.md`
   - Add a new section after the "Claude Code Permissions (YOLO Mode)" section (after line 77)
   - Title: `## Auto-Memory`
   - Content should cover:
     - What auto-memory is and why it's enabled
     - Where memory is stored (`~/.claude/projects/<project>/memory/`)
     - How MEMORY.md works (first 200 lines loaded at session start, topic files loaded on-demand)
     - That memory persists via the Docker volume mount (`claude-code-config-${devcontainerId}`)
     - That memory is user-local and not committed to git
     - How to disable if needed (set `CLAUDE_CODE_DISABLE_AUTO_MEMORY=1` in `containerEnv`)

3. **Update devcontainer rule file**
   - Open `/workspace/.claude/rules/devcontainer.md`
   - Add a bullet point noting auto-memory is enabled via `CLAUDE_CODE_DISABLE_AUTO_MEMORY=0`

4. **Verify auto-memory is active**
   - After rebuilding the devcontainer (or in the current session), verify:
     - `echo $CLAUDE_CODE_DISABLE_AUTO_MEMORY` outputs `0`
     - `ls ~/.claude/projects/-workspace/memory/` shows `MEMORY.md`
   - Note: Since auto-memory is already working, this is a confirmation step rather than a first-time setup

5. **Run lint and tests**
   - `pnpm lint` -- should pass (changes are only to JSON and Markdown files)
   - `pnpm test` -- should pass (no code changes)

6. **Create feature branch and PR**
   - Branch name: `chore/enable-auto-memory`
   - Commit the `devcontainer.json`, `README.md`, and `devcontainer.md` changes
   - Create PR for human review

### Testing Strategy

- **Automated**: `pnpm lint` and `pnpm test` should pass without issues (only JSON and Markdown changes)
- **Manual**: After container rebuild, verify `echo $CLAUDE_CODE_DISABLE_AUTO_MEMORY` returns `0` and memory directory exists at `~/.claude/projects/-workspace/memory/`
- **No visual verification needed**: This is a devcontainer config change, not a UI change

### Notes

- No ADR needed: This is a simple configuration enablement, not an architectural change or new dependency
- The Docker volume mount for `~/.claude` (already configured) is what provides persistence -- no additional volume configuration is required
- The existing `MEMORY.md` content (testing patterns, lint rules, GraphQL upload patterns, codegen, workflow) is valuable and will be preserved since it lives in the volume, not in the repo
