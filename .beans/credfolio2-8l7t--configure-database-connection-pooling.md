---
# credfolio2-8l7t
title: Configure database connection pooling
status: todo
type: task
priority: normal
created_at: 2026-02-02T08:32:29Z
updated_at: 2026-02-05T17:48:38Z
parent: credfolio2-abtx
---

The backend has no explicit connection pool limits configured on the Bun ORM connection. Under load, this will exhaust PostgreSQL's max_connections (default 100) and cause all requests to fail.

## Problem

In `src/backend/internal/infrastructure/database/database.go`, the `sql.DB` is created without pool limits:

```go
sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(cfg.URL())))
db := bun.NewDB(sqldb, pgdialect.New())
// No SetMaxOpenConns(), SetMaxIdleConns(), or SetConnMaxLifetime()
```

**Dangerous defaults:**
- MaxOpenConns = 0 (unlimited)
- MaxIdleConns = 2
- ConnMaxLifetime = 0 (never expires)

## Impact

With 50+ concurrent users:
1. Each GraphQL mutation opens multiple connections
2. Bun creates connections without limit
3. River's pgxpool (10-16 workers) also needs connections
4. PostgreSQL hits max_connections → all requests fail
5. No automatic recovery

## Solution

Add explicit pool configuration:

```go
sqldb.SetMaxOpenConns(25)              // Max concurrent connections
sqldb.SetMaxIdleConns(5)               // Keep some warm for reuse
sqldb.SetConnMaxLifetime(5 * time.Minute) // Recycle stale connections
```

**Connection budget (100 total):**
- Bun ORM: 25
- River pgxpool: 20
- Admin/migrations: 10
- Headroom: 45

## Checklist

- [ ] Add SetMaxOpenConns(25) to database.go
- [ ] Add SetMaxIdleConns(5) to database.go
- [ ] Add SetConnMaxLifetime(5 * time.Minute) to database.go
- [ ] Consider making values configurable via env vars
- [ ] Optionally configure River's pgxpool explicitly in river.go

## Definition of Done
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Branch pushed and PR created for human review