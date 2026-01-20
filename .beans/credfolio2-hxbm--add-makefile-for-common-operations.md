---
# credfolio2-hxbm
title: Add Makefile for common operations
status: scrapped
type: task
priority: normal
created_at: 2026-01-20T11:26:30Z
updated_at: 2026-01-20T15:23:04Z
parent: credfolio2-jpin
---

Create a Makefile with targets for common development tasks.

## Targets to Include
- `make dev` - Start all services for development
- `make build` - Build all artifacts
- `make test` - Run all tests
- `make lint` - Run linters
- `make migrate-up` - Run database migrations
- `make migrate-down` - Rollback migrations
- `make migrate-create name=...` - Create new migration
- `make docker-up` - Start Docker Compose services
- `make docker-down` - Stop Docker Compose services
- `make clean` - Clean build artifacts

## Acceptance Criteria
- All targets work correctly
- Help text via `make help`
- Works on macOS and Linux