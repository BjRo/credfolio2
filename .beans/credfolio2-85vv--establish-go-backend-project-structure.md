---
# credfolio2-85vv
title: Establish Go backend project structure
status: in-progress
type: task
priority: normal
created_at: 2026-01-20T11:26:26Z
updated_at: 2026-01-20T12:33:59Z
parent: credfolio2-jpin
blocking:
    - credfolio2-n907
---

Reorganize the Go backend to follow Clean Architecture principles.

## Directory Structure
```
src/backend/
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── internal/
│   ├── config/              # Configuration loading
│   ├── domain/              # Domain models (entities)
│   ├── repository/          # Data access interfaces & implementations
│   ├── service/             # Business logic
│   ├── handler/             # HTTP/GraphQL handlers
│   └── infrastructure/      # External integrations (LLM, storage)
├── pkg/                     # Shared utilities (if needed)
├── migrations/              # Database migrations
└── go.mod
```

## Requirements
- Clear separation of concerns
- Dependency injection friendly
- Testable architecture

## Acceptance Criteria
- Existing code migrated to new structure
- All imports updated
- Project builds and tests pass
- README updated with structure explanation

## Checklist
- [x] Create internal/config directory with package doc
- [x] Create internal/domain directory with package doc
- [x] Create internal/repository directory with package doc
- [x] Create internal/service directory with package doc
- [x] Create internal/handler directory with package doc
- [x] Create internal/infrastructure directory with package doc
- [x] Create pkg directory with package doc
- [x] Create migrations directory
- [x] Move handlers out of main.go into handler package
- [x] Update main.go to use new handler package
- [x] Add tests for handlers
- [x] Verify build passes
- [x] Run full pnpm test to ensure CI passes