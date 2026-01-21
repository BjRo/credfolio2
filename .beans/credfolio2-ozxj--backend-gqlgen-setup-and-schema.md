---
# credfolio2-ozxj
title: 'Backend: gqlgen setup and schema'
status: completed
type: feature
priority: normal
created_at: 2026-01-21T14:27:41Z
updated_at: 2026-01-21T15:20:32Z
parent: credfolio2-zbqk
blocking:
    - credfolio2-5l5u
    - credfolio2-ui1x
---

Set up gqlgen for GraphQL API generation with schema-first approach.

## Scope
- Add gqlgen dependency
- Initialize gqlgen configuration
- Define GraphQL schema for existing entities (User, File, ReferenceLetter)
- Configure GraphQL playground for development
- Generate resolver stubs

## Schema Types
Based on existing domain entities in `internal/domain/entities.go`:

- **User**: id, email, name, createdAt, updatedAt
- **File**: id, filename, contentType, sizeBytes, storageKey, createdAt, user
- **ReferenceLetter**: id, title, authorName, authorTitle, organization, dateWritten, rawText, extractedData, status, createdAt, updatedAt, user, file

## Checklist
- [x] Add gqlgen dependency (`go get github.com/99designs/gqlgen`)
- [x] Initialize gqlgen (`go run github.com/99designs/gqlgen init`)
- [x] Configure gqlgen.yml for project structure
- [x] Define schema.graphqls with User, File, ReferenceLetter types
- [x] Add Query type with appropriate queries
- [x] Run code generation (`go generate ./...`)
- [x] Wire up GraphQL handler in chi router at `/graphql`
- [x] Enable GraphQL playground at `/playground` (dev only)

## Notes
- Follow schema-first approach per gqlgen best practices
- Keep domain entities separate from GraphQL models (use resolvers to map)