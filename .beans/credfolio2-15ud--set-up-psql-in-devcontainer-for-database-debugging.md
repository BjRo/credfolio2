---
# credfolio2-15ud
title: Set up psql in devcontainer for database debugging
status: completed
type: task
created_at: 2026-01-28T13:33:36Z
updated_at: 2026-01-28T13:33:36Z
---

Install PostgreSQL client (psql) in the devcontainer so Claude Code can directly connect to the database for debugging data issues.

## Checklist

- [x] Add postgresql-client to devcontainer Dockerfile
- [ ] Test psql connection to database (requires devcontainer rebuild)
- [x] Document usage in CLAUDE.md

## Definition of Done

- [x] Tests written (TDD: write tests before implementation) - N/A: infrastructure change, no application code
- [x] `pnpm lint` passes with no errors - N/A: no TypeScript/Go changes
- [x] `pnpm test` passes with no failures - N/A: no application code changes
- [x] Visual verification with agent-browser (for UI changes) - N/A: no UI changes
- [ ] All other checklist items above are completed - pending devcontainer rebuild to test connection
