---
# credfolio2-zbqk
title: GraphQL API
status: todo
type: epic
priority: normal
created_at: 2026-01-20T11:24:59Z
updated_at: 2026-01-21T14:24:37Z
parent: credfolio2-tikg
---

Set up GraphQL API with gqlgen for serving profile and reference letter data.

## Goals
- Configure gqlgen with schema-first approach
- Define GraphQL schema for profiles, positions, skills
- Implement resolvers connecting to Bun ORM
- Set up URQL client on frontend

## Checklist
- [ ] Add gqlgen dependency and initialize
- [ ] Define GraphQL schema (profile, position, skill, referenceFile types)
- [ ] Generate resolver stubs
- [ ] Implement profile queries
- [ ] Implement reference letter queries
- [ ] Set up URQL client in Next.js
- [ ] Create typed query hooks
- [ ] Add GraphQL playground for development