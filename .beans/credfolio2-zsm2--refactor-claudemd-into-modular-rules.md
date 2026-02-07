---
# credfolio2-zsm2
title: Refactor CLAUDE.md into modular rules
status: todo
type: task
created_at: 2026-02-07T15:59:13Z
updated_at: 2026-02-07T15:59:13Z
parent: credfolio2-ynmd
---

Break up the 327-line root CLAUDE.md into focused, modular files under `.claude/rules/`.

## Why

The entire CLAUDE.md loads into context at session start. Much of its content is situational — database debugging commands aren't needed when working on frontend components, and visual verification instructions aren't needed when writing migrations. Modular rules with path-specific scoping load content only when relevant, keeping the base context lean.

## What

Audit the current CLAUDE.md and extract concerns into `.claude/rules/`:

### Keep in CLAUDE.md (always-loaded essentials)
- Project overview & directory structure
- Key technical decisions (package manager, build system, stacks)
- Common commands (build, dev, clean)
- Git details & workflow rules
- The "STOP — Before Marking Any Work Complete" checklist
- Decision documentation reference

### Move to `.claude/rules/` (situational content)

| File | Content | Path scope |
|------|---------|------------|
| `database.md` | PostgreSQL connection details, psql commands, migration commands, database debugging | `src/backend/**` |
| `dev-servers.md` | Starting dev servers, killing stale processes, port cleanup, common issues | None (unconditional but rarely needed) |
| `visual-verification.md` | Agent-browser workflow, fixture resume, upload instructions | `src/frontend/**` |
| `backend.md` | Backend-specific patterns, Go conventions, file locations | `src/backend/**` |
| `frontend.md` | Frontend-specific patterns, Next.js/React/Tailwind conventions, file locations | `src/frontend/**` |
| `devcontainer.md` | Devcontainer notes, tool management, rebuild info | `.devcontainer/**` |

### Path-specific scoping example
```yaml
---
paths:
  - "src/backend/**"
---

# Database Debugging with psql
...
```

## Approach
- Read the current CLAUDE.md carefully
- Identify each section and whether it's always-needed or situational
- Create the rules files with appropriate `paths` frontmatter
- Slim down CLAUDE.md to the essentials
- Verify nothing is lost by comparing before/after

## Definition of Done
- [ ] `.claude/rules/` directory created with modular rule files
- [ ] Path-specific scoping applied where appropriate
- [ ] Root CLAUDE.md reduced to essential always-needed content
- [ ] No information lost — all content preserved across files
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Branch pushed and PR created for human review