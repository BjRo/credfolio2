---
# credfolio2-ubn6
title: Create database reset script
status: in-progress
type: task
priority: normal
created_at: 2026-02-07T16:29:34Z
updated_at: 2026-02-08T10:21:04Z
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
.claude/scripts/reset-db.sh          # Reset dev database (with confirmation prompt)
.claude/scripts/reset-db.sh test     # Reset test database (no prompt)
.claude/scripts/reset-db.sh all      # Reset both (with confirmation prompt)
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
- [x] Script created at `.claude/scripts/reset-db.sh`
- [x] Script supports dev, test, and all environments
- [x] Script drops, recreates, and migrates the database
- [x] Script works in the devcontainer (no make dependency)
- [x] CLAUDE.md updated to reference the script
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] Branch pushed and PR created for human review

---

## Implementation Plan

### Approach

Create a bash script at `.claude/scripts/reset-db.sh` that provides a safe, one-command way to reset databases. The script will:
- Parse environment argument (dev/test/all), defaulting to dev
- Prompt for confirmation when resetting dev database (but never for test)
- Use `psql` to drop and recreate databases
- Use `migrate` CLI to run all migrations
- Provide clear output with color-coded status messages
- Follow existing script patterns from `.claude/scripts/start-work.sh`

This eliminates the common failure modes: wrong connection details, using unavailable `make` commands, or forgetting to run migrations after reset.

### Files to Create/Modify

- `.claude/scripts/reset-db.sh` (create) — Main reset script with confirmation prompts
- `/workspace/CLAUDE.md` (modify) — Add reset-db.sh to Common Commands section
- `.claude/rules/database.md` (modify) — Add reset script documentation

### Steps

1. **Create the reset script at `.claude/scripts/reset-db.sh`**
   - Add shebang and `set -e` for error handling
   - Define color variables (RED, GREEN, YELLOW, NC) matching `start-work.sh` pattern
   - Define connection constants:
     - `POSTGRES_HOST=credfolio2-postgres`
     - `POSTGRES_PORT=5432`
     - `POSTGRES_USER=credfolio`
     - `POSTGRES_PASSWORD=credfolio_dev`
     - `MIGRATIONS_DIR=/workspace/src/backend/migrations`
   - Parse command-line argument (default to "dev")
   - Validate argument is one of: dev, test, all
   - Show usage message if invalid argument

2. **Implement confirmation prompt logic**
   - Create `confirm_reset()` function that:
     - Takes database name as parameter
     - Prints warning message in YELLOW
     - Prompts: "Reset {database_name}? This will DELETE ALL DATA. [y/N]: "
     - Reads user input (case-insensitive)
     - Returns 0 for yes/y, 1 otherwise
   - Only call `confirm_reset()` for dev database or "all" mode
   - Never prompt for test database (auto-proceed)

3. **Implement database reset logic**
   - Create `reset_database()` function that:
     - Takes database name as parameter (credfolio_dev or credfolio_test)
     - Exports `PGPASSWORD=$POSTGRES_PASSWORD` to avoid password prompts
     - Step 1: Drop database using `psql`:
       ```bash
       psql -h $POSTGRES_HOST -U $POSTGRES_USER -d postgres -c "DROP DATABASE IF EXISTS $db_name"
       ```
     - Step 2: Create database using `psql`:
       ```bash
       psql -h $POSTGRES_HOST -U $POSTGRES_USER -d postgres -c "CREATE DATABASE $db_name"
       ```
     - Step 3: Run migrations using `migrate`:
       ```bash
       migrate -path $MIGRATIONS_DIR -database "postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST:$POSTGRES_PORT/$db_name?sslmode=disable" up
       ```
     - Step 4: Show migration version status:
       ```bash
       migrate -path $MIGRATIONS_DIR -database "postgres://..." version
       ```
     - Print success message in GREEN with migration count/version

4. **Implement main script flow**
   - Parse environment argument: `ENV=${1:-dev}`
   - Validate environment is one of: dev, test, all
   - Branch based on environment:
     - **dev**: Call `confirm_reset credfolio_dev`, if confirmed call `reset_database credfolio_dev`
     - **test**: Call `reset_database credfolio_test` (no confirmation)
     - **all**: Call `confirm_reset "BOTH dev and test databases"`, if confirmed call `reset_database credfolio_dev` then `reset_database credfolio_test`
   - Exit with appropriate messages if user cancels

5. **Make script executable**
   - Run: `chmod +x .claude/scripts/reset-db.sh`

6. **Update `/workspace/CLAUDE.md`**
   - Find the "Common Commands" section (around line 89)
   - Add database reset commands after the existing commands:
     ```bash
     # Database reset (requires confirmation for dev)
     .claude/scripts/reset-db.sh          # Reset dev database
     .claude/scripts/reset-db.sh test     # Reset test database
     .claude/scripts/reset-db.sh all      # Reset both databases
     ```

7. **Update `.claude/rules/database.md`**
   - Add new section "Database Reset Script" after "Database Debugging with psql"
   - Document the script usage, arguments, and confirmation behavior
   - Include examples for all three modes (dev, test, all)
   - Note the confirmation prompt for dev database protection

8. **Test the script**
   - Test with invalid argument: `.claude/scripts/reset-db.sh invalid` (should show usage and exit)
   - Test test mode: `.claude/scripts/reset-db.sh test` (should reset without prompt)
   - Test dev mode cancellation: `.claude/scripts/reset-db.sh` then enter "n" (should exit without changes)
   - Test dev mode confirmation: `.claude/scripts/reset-db.sh` then enter "y" (should reset and run migrations)
   - Test all mode: `.claude/scripts/reset-db.sh all` then enter "y" (should reset both)
   - Verify migrations are applied by checking migration version after reset

### Testing Strategy

**Unit/Integration Tests**
- No Go/TypeScript code changes, so existing `pnpm test` should pass unchanged
- The script itself is bash, so manual testing is required

**Manual Testing**
- Test all three modes (dev, test, all) as described in step 8 above
- Verify confirmation prompts work correctly
- Verify database is actually reset (check tables are empty)
- Verify migrations are applied (check migration version)
- Verify error handling (e.g., if postgres is not running)

**Validation**
- Run `pnpm lint` to ensure no linting issues (should pass as no code changes)
- Run `pnpm test` to ensure all tests pass
- Manually verify CLAUDE.md and database.md formatting is correct

### Edge Cases and Error Handling

- **Postgres not running**: Script will fail with clear psql error
- **Invalid migrate path**: Script will fail with migrate error
- **User cancellation**: Script exits with "Cancelled" message, no changes made
- **Invalid argument**: Script shows usage message and exits with code 1
- **Empty database**: Script should work fine (drop/create succeeds even if database is empty)

### Open Questions

None — the requirements are clear and the approach is straightforward.
