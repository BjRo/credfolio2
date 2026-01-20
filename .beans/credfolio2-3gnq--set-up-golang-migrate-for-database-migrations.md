---
# credfolio2-3gnq
title: Set up golang-migrate for database migrations
status: completed
type: task
priority: normal
created_at: 2026-01-20T11:26:06Z
updated_at: 2026-01-20T15:18:53Z
parent: credfolio2-jpin
blocking:
    - credfolio2-x9d6
---

Configure database migration tooling using golang-migrate.

## Requirements

- Install golang-migrate CLI
- Create migrations directory structure
- Initial migration for schema setup
- Make target or script for running migrations
- Make sure that everything is compatible with bun which we will use later for db access
- Make sure that we have a separate test and development database(s)

## Acceptance Criteria

- Can create new migrations with `make migration name=...`
- Can run migrations up with `make migrate-up`
- Can rollback with `make migrate-down`
- Migration state is tracked in database

## Technical Notes

- Use file:// source for local development
- Store migrations in src/backend/migrations/
- Use timestamp versioning (like Rails: 20260120143000_name.up.sql)
- Test database: credfolio_test on same postgres instance

## Checklist

- [x] Add golang-migrate to devcontainer Dockerfile
- [x] Update docker-compose.yml to create test database
- [x] Add test database config to .env.example
- [x] Update config.go for test database support
- [x] Create Makefile with migration commands
- [x] Create initial migration (schema_migrations table is auto-created by golang-migrate)
- [x] Test migration commands work
- [x] Update CLAUDE.md with new commands
