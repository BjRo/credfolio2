---
# credfolio2-zioz
title: Data Access Layer
status: completed
type: task
priority: normal
created_at: 2026-01-21T13:48:53Z
updated_at: 2026-01-21T14:10:44Z
parent: credfolio2-tikg
blocking:
    - credfolio2-zbqk
---

Set up the Go data access layer connecting the database schema to application code. This provides the foundation for all data operations.

## Goals

- Establish database connection with Bun ORM
- Define domain entities matching the database schema
- Create repository interfaces and implementations
- Enable creating and retrieving users, files, and reference letters

## Checklist

- [x] Add Bun ORM dependency and configure database connection
- [x] Create domain entities (User, File, ReferenceLetter) with Bun struct tags
- [x] Define repository interfaces in domain layer
- [x] Implement PostgreSQL repositories using Bun
- [x] Add database health check to existing /health endpoint
- [x] Write integration tests for repositories
- [x] Document Bun ORM decision in /decisions

## Technical Notes

- Use [Bun](https://github.com/uptrace/bun) for database access (SQL-first ORM for Go)
- Bun uses pgx as the underlying driver
- Connection pool configured via environment variables
- Repositories return domain entities, not database rows
- Integration tests use credfolio_test database

## Decision: Bun ORM

**Choice**: Use Bun ORM instead of raw pgx or other ORMs (GORM, sqlx, etc.)

**Rationale**:

- SQL-first approach: write real SQL, not a DSL - easier debugging and optimization
- Lightweight: less magic than GORM, more ergonomic than raw pgx
- Built on pgx: gets pgx performance benefits
- Struct mapping with tags for clean domain entities
- Good support for PostgreSQL-specific features (JSONB, arrays, etc.)
- Active maintenance by uptrace team
