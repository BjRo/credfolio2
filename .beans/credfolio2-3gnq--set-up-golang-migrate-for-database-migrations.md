---
# credfolio2-3gnq
title: Set up golang-migrate for database migrations
status: in-progress
type: task
priority: normal
created_at: 2026-01-20T11:26:06Z
updated_at: 2026-01-20T14:34:27Z
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

## Acceptance Criteria
- Can create new migrations with `make migration name=...`
- Can run migrations up with `make migrate-up`
- Can rollback with `make migrate-down`
- Migration state is tracked in database

## Technical Notes
- Use file:// source for local development
- Store migrations in src/backend/migrations/
- Use sequential versioning (not timestamps)