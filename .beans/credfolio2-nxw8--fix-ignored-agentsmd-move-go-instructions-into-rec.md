---
# credfolio2-nxw8
title: Fix ignored Agents.md — move Go instructions into recognized memory
status: in-progress
type: task
priority: normal
created_at: 2026-02-07T16:05:55Z
updated_at: 2026-02-07T19:07:13Z
parent: credfolio2-ynmd
---

The root `Agents.md` file contains Go-specific development tips but is completely ignored because Claude Code's memory system only recognizes specific file names and locations.

## The Problem

Claude Code loads memory from these sources (and only these):
- `CLAUDE.md` / `.claude/CLAUDE.md` (project memory)
- `CLAUDE.local.md` (local project memory)
- `.claude/rules/*.md` (modular rules)
- `~/.claude/CLAUDE.md` (user memory)
- Auto-memory at `~/.claude/projects/<project>/memory/`

A file called `Agents.md` is **not recognized** — it's never loaded into context automatically. Claude only sees it if it happens to read the file during codebase exploration, which explains why the instructions are consistently ignored.

## Current content of Agents.md

```markdown
# For developing in Golang:
- Use `go mod download -json MODULE` to see source files from a dependency
- Use `go doc foo.Bar` or `go doc -all foo` to read documentation
- Use `go run .` or `go run ./cmd/foo` instead of `go build` to avoid build artifacts
```

## Solution

Move this content into `.claude/rules/golang.md` with path-specific scoping so it loads when working on Go files:

```yaml
---
paths:
  - "src/backend/**"
  - "**/*.go"
---
```

This ensures the instructions are:
1. Actually loaded by Claude Code's memory system
2. Only loaded when working on Go/backend code (not frontend work)
3. Discoverable via `/memory` command

Delete the original `Agents.md` after migration.

## Note

This overlaps with credfolio2-zsm2 (Refactor CLAUDE.md into modular rules). Can be done as part of that task or independently — the fix is small and self-contained.

## Definition of Done
- [x] `.claude/rules/golang.md` created with path-scoped Go instructions
- [x] Content from `Agents.md` migrated
- [x] Original `Agents.md` deleted
- [x] Verified rules load when working on `src/backend/` files
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [ ] Branch pushed and PR created for human review