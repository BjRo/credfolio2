---
paths:
  - "src/backend/**"
---

# Database Commands

## Common Database Commands

```bash
# Database (run from host, not devcontainer)
docker-compose up -d           # Start PostgreSQL and MinIO
docker-compose down            # Stop services

# Migrations (run from src/backend/)
cd src/backend
make help                      # Show all migration commands
make migration name=create_users  # Create new migration
make migrate-up                # Run pending migrations (on credfolio_dev)
make migrate-down              # Rollback last migration
CREDFOLIO_ENV=test make migrate-up   # Run migrations on test database
make migrate-status            # Show migration status for all environments

# Database debugging with psql (from devcontainer)
psql -h credfolio2-postgres -U credfolio -d credfolio_dev   # Connect to dev database
psql -h credfolio2-postgres -U credfolio -d credfolio_test  # Connect to test database
# Password: credfolio_dev
```

## Database Debugging with psql

The devcontainer includes `postgresql-client` for direct database access. This is useful for debugging data issues, inspecting table contents, and running ad-hoc queries.

### Connection Details

| Setting | Value |
|---------|-------|
| Host | `credfolio2-postgres` |
| Port | `5432` |
| User | `credfolio` |
| Password | `credfolio_dev` |
| Dev Database | `credfolio_dev` |
| Test Database | `credfolio_test` |

### Quick Commands

```bash
# Connect to dev database (interactive)
psql -h credfolio2-postgres -U credfolio -d credfolio_dev

# Connect to test database
PGPASSWORD=credfolio_dev psql -h credfolio2-postgres -U credfolio -d credfolio_test

# Run a single query
PGPASSWORD=credfolio_dev psql -h credfolio2-postgres -U credfolio -d credfolio_dev -c "SELECT * FROM users;"

# List all tables
PGPASSWORD=credfolio_dev psql -h credfolio2-postgres -U credfolio -d credfolio_dev -c "\dt"

# Describe a table structure
PGPASSWORD=credfolio_dev psql -h credfolio2-postgres -U credfolio -d credfolio_dev -c "\d users"

# Export query results to file
PGPASSWORD=credfolio_dev psql -h credfolio2-postgres -U credfolio -d credfolio_dev -c "SELECT * FROM users;" > /tmp/users.txt
```

### Environment Variable for Password

To avoid typing the password repeatedly, set `PGPASSWORD`:

```bash
export PGPASSWORD=credfolio_dev
psql -h credfolio2-postgres -U credfolio -d credfolio_dev
```

## Database Reset Script

The `.claude/scripts/reset-db.sh` script provides a safe, one-command way to reset databases. It drops the database, recreates it, and runs all migrations.

### Usage

```bash
.claude/scripts/reset-db.sh [environment]
```

### Arguments

| Environment | Behavior | Confirmation Required |
|-------------|----------|----------------------|
| `dev` (default) | Reset `credfolio_dev` | Yes |
| `test` | Reset `credfolio_test` | No |
| `all` | Reset both databases | Yes |

### Examples

```bash
# Reset dev database (prompts for confirmation)
.claude/scripts/reset-db.sh

# Reset test database (no prompt, auto-proceeds)
.claude/scripts/reset-db.sh test

# Reset both databases (prompts for confirmation)
.claude/scripts/reset-db.sh all
```

### What It Does

1. Drops the target database (if exists)
2. Creates a fresh database
3. Runs all migrations from `src/backend/migrations/`
4. Displays the final migration version

### Safety Features

- **Confirmation prompts** protect the dev database from accidental wipes
- Test database auto-proceeds (no prompt) for CI/automation workflows
- Clear color-coded output shows progress and status
- Script uses `set -e` to abort on errors

### When to Use

- Reset development database to clean state
- Clear test database between test runs
- Fix migration issues by starting fresh
- Remove accumulated test data
