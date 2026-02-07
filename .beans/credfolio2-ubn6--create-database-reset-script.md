---
# credfolio2-ubn6
title: Create database reset script
status: todo
type: task
created_at: 2026-02-07T16:29:34Z
updated_at: 2026-02-07T16:29:34Z
parent: credfolio2-ynmd
---

Create a one-command script to wipe the dev database and re-run all migrations.

## Why

Resetting the database currently requires knowing the psql connection details, dropping/recreating the database, and running migrations via the migrate CLI. Claude often fumbles this sequence — wrong host, wrong password, forgetting the test database, or trying to use `make` (which is not available in the devcontainer). A single script eliminates this detour.

## What

Create `.claude/scripts/reset-db.sh` that:
1. Accepts an optional environment argument (defaults to `dev`):
   - `dev` → `credfolio_dev`
   - `test` → `credfolio_test`
   - `all` → both databases
2. Drops and recreates the target database(s)
3. Runs all migrations via the `migrate` CLI
4. Reports success with migration status

### Connection details (from docker-compose)
- Host: `credfolio2-postgres`
- Port: `5432`
- User: `credfolio`
- Password: `credfolio_dev`

### Example usage
```bash
.claude/scripts/reset-db.sh          # Reset dev database
.claude/scripts/reset-db.sh test     # Reset test database
.claude/scripts/reset-db.sh all      # Reset both
```

## Notes

- `make` is NOT available in the devcontainer — the script must call `migrate` directly
- The script should use `PGPASSWORD` env var to avoid password prompts
- Migrations are at `src/backend/migrations/`
- The migrate CLI DSN format: `postgres://credfolio:credfolio_dev@credfolio2-postgres:5432/credfolio_dev?sslmode=disable`

## Integration

- Add to CLAUDE.md under Common Commands
- Consider adding a permission rule in allow list

## Definition of Done
- [ ] Script created at `.claude/scripts/reset-db.sh`
- [ ] Script supports dev, test, and all environments
- [ ] Script drops, recreates, and migrates the database
- [ ] Script works in the devcontainer (no make dependency)
- [ ] CLAUDE.md updated to reference the script
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Branch pushed and PR created for human review