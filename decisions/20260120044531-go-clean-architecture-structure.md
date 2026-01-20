# Go Backend Clean Architecture Structure

**Date**: 2026-01-20
**Bean**: credfolio2-85vv

## Context

The Go backend started as a simple `main.go` with inline HTTP handlers. As the project grows to include database access, external service integrations, and business logic, we need a clear structure that separates concerns and enables testability.

## Decision

Adopted a Clean Architecture-inspired directory structure:

```
src/backend/
├── cmd/server/main.go      # Application entry point, wiring
├── internal/
│   ├── config/             # Configuration loading
│   ├── domain/             # Core business entities (no dependencies)
│   ├── repository/         # Data access interfaces & implementations
│   ├── service/            # Business logic / use cases
│   ├── handler/            # HTTP handlers (presentation layer)
│   └── infrastructure/     # External integrations (LLM, storage, etc.)
├── pkg/                    # Shared utilities (if needed)
├── migrations/             # Database migrations
└── go.mod
```

## Reasoning

- **Separation of concerns**: Each layer has a single responsibility
- **Testability**: Business logic in `service/` can be tested without HTTP or database
- **Dependency direction**: Dependencies point inward (handlers → services → domain)
- **Go idioms**: Uses `internal/` to prevent external imports, follows standard Go project layout conventions
- **Flexibility**: Infrastructure can be swapped (e.g., different database, different LLM provider) without touching business logic

Alternatives considered:
- Flat structure: Rejected as it doesn't scale and makes testing harder
- DDD with aggregates: Too heavyweight for current project scope

## Consequences

- All new HTTP handlers go in `internal/handler/`
- Business logic belongs in `internal/service/`, not in handlers
- Domain entities in `internal/domain/` should have no external dependencies
- External service clients (LLM, S3, etc.) live in `internal/infrastructure/`
- Configuration loading and environment variables handled in `internal/config/`
