---
# credfolio2-zbqk
title: GraphQL API
status: in-progress
type: epic
priority: normal
created_at: 2026-01-20T11:24:59Z
updated_at: 2026-01-21T14:27:28Z
parent: credfolio2-tikg
---

Set up GraphQL API with gqlgen for serving data to the frontend.

## Goals
- Configure gqlgen with schema-first approach
- Define GraphQL schema for existing domain entities (User, File, ReferenceLetter)
- Implement resolvers connecting to Bun ORM repositories
- Set up URQL client on frontend with typed hooks

## Scope Refinement
Initial implementation focuses on existing domain entities to support the "First Vertical Slice" milestone. Profile/Position/Skill types are tracked separately for future expansion.

## Child Beans
- `credfolio2-ozxj` - Backend: gqlgen setup and schema
- `credfolio2-5l5u` - Backend: Implement GraphQL resolvers (blocked by ozxj)
- `credfolio2-ui1x` - Frontend: URQL client setup (blocked by ozxj)
- `credfolio2-7h7b` - Future: Profile, Position, Skill types (draft)