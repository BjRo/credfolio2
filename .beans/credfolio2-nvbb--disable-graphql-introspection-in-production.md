---
# credfolio2-nvbb
title: Disable GraphQL introspection in production
status: draft
type: task
created_at: 2026-01-21T15:11:15Z
updated_at: 2026-01-21T15:11:15Z
parent: credfolio2-wxn8
---

Disable GraphQL schema introspection queries in production to prevent exposing the type system.

## Context
Introspection queries (`__schema`, `__type`) allow clients to discover the entire API schema. This is useful for tooling but can help attackers understand the API surface in production.

## Acceptance Criteria
- [ ] Introspection queries return an error in production
- [ ] Introspection works normally in development
- [ ] Error message doesn't leak information

## Implementation Notes
- gqlgen supports disabling introspection via server configuration
- Consider `srv.SetRecoverFunc()` and introspection settings