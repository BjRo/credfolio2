---
# credfolio2-zbqk
title: GraphQL API
status: draft
type: epic
created_at: 2026-01-20T11:24:59Z
updated_at: 2026-01-20T11:24:59Z
parent: credfolio2-1fwu
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