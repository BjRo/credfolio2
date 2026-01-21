# Bun ORM for Database Access

**Date**: 2026-01-21
**Bean**: credfolio2-zioz

## Context

The credfolio2 backend needs a database access layer to interact with PostgreSQL. We needed to choose between several options:

- Raw `database/sql` with `pgx` driver
- GORM (full-featured ORM)
- sqlx (SQL-focused extension to database/sql)
- Bun (SQL-first ORM built on pgx)

## Decision

Adopted [Bun ORM](https://github.com/uptrace/bun) for all database operations.

Key implementation details:
- Domain entities defined in `internal/domain/entities.go` with Bun struct tags
- Repository interfaces defined in `internal/domain/repository.go`
- PostgreSQL implementations in `internal/repository/postgres/`
- Database connection via `internal/infrastructure/database/database.go`
- Health check integrated into `/health` endpoint

## Reasoning

Bun was selected for several reasons:

1. **SQL-first approach**: Write actual SQL queries rather than a DSL. This makes debugging easier and allows PostgreSQL-specific optimizations.

2. **Lightweight**: Less magic than GORM - you see what SQL is being executed. More ergonomic than raw pgx for common operations.

3. **Built on pgx**: Uses pgx as the underlying driver, getting its performance benefits including connection pooling and prepared statement caching.

4. **Struct mapping**: Clean mapping between Go structs and database rows using struct tags, without heavy reflection overhead.

5. **PostgreSQL features**: Good support for PostgreSQL-specific features like JSONB, arrays, and custom types that credfolio2 uses.

6. **Active maintenance**: Actively maintained by the uptrace team with regular releases.

Alternatives considered:
- **GORM**: Too much magic, harder to debug, generates unpredictable SQL
- **Raw pgx**: Very low-level, requires writing boilerplate for common operations
- **sqlx**: Good but Bun provides better struct mapping and query building

## Consequences

1. **Dependencies added**:
   - `github.com/uptrace/bun`
   - `github.com/uptrace/bun/dialect/pgdialect`
   - `github.com/uptrace/bun/driver/pgdriver`
   - `github.com/google/uuid`

2. **Entity definitions**: All domain entities must include `bun.BaseModel` and use Bun struct tags for column mapping.

3. **Repository pattern**: Repositories accept `*bun.DB` and implement domain interfaces, keeping the domain layer independent of Bun.

4. **Testing**: Integration tests run against `credfolio_test` database (set via `CREDFOLIO_ENV=test`).

5. **Migrations**: Database schema managed separately via golang-migrate; Bun does not handle migrations.
