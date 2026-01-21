---
# credfolio2-5l5u
title: 'Backend: Implement GraphQL resolvers'
status: todo
type: feature
priority: normal
created_at: 2026-01-21T14:27:51Z
updated_at: 2026-01-21T14:28:41Z
parent: credfolio2-zbqk
---

Implement GraphQL resolvers that connect to Bun repositories.

## Dependencies
- Requires credfolio2-ozxj (gqlgen setup) to be completed first

## Scope
Implement resolvers for all queries and types:

### Queries to implement
- `user(id: ID!): User` - Get user by ID
- `referenceLetters(userId: ID!): [ReferenceLetter!]!` - Get user's reference letters
- `referenceLetter(id: ID!): ReferenceLetter` - Get single reference letter
- `files(userId: ID!): [File!]!` - Get user's files
- `file(id: ID!): File` - Get single file

### Type resolvers
- ReferenceLetter.user - Resolve user relation
- ReferenceLetter.file - Resolve file relation
- File.user - Resolve user relation

## Checklist
- [ ] Inject repository dependencies into resolver struct
- [ ] Implement user query resolver
- [ ] Implement referenceLetters query resolver
- [ ] Implement referenceLetter query resolver
- [ ] Implement files query resolver
- [ ] Implement file query resolver
- [ ] Implement ReferenceLetter.user field resolver
- [ ] Implement ReferenceLetter.file field resolver
- [ ] Implement File.user field resolver
- [ ] Add resolver tests
- [ ] Verify queries work via GraphQL playground

## Notes
- Use TDD approach per project workflow
- Resolvers should handle errors gracefully, returning nil for not found