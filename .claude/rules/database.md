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
