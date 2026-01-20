---
# credfolio2-85vv
title: Establish Go backend project structure
status: todo
type: task
priority: normal
created_at: 2026-01-20T11:26:26Z
updated_at: 2026-01-20T11:26:40Z
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