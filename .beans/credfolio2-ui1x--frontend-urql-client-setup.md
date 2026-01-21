---
# credfolio2-ui1x
title: 'Frontend: URQL client setup'
status: completed
type: feature
priority: normal
created_at: 2026-01-21T14:28:04Z
updated_at: 2026-01-21T16:20:38Z
parent: credfolio2-zbqk
---

Set up URQL GraphQL client in the Next.js frontend.

## Dependencies
- Requires credfolio2-ozxj (gqlgen setup) to be completed first (needs schema)

## Scope
- Add URQL dependencies
- Configure URQL client to connect to backend GraphQL endpoint
- Set up code generation for typed query hooks
- Create initial typed hooks for reference letter queries

## Checklist
- [x] Add URQL dependencies (`@urql/core`, `@urql/next`, `urql`)
- [x] Add GraphQL codegen dependencies (`@graphql-codegen/cli`, etc.)
- [x] Configure codegen.ts for URQL plugin
- [x] Copy schema from backend or configure introspection
- [x] Create URQL client provider
- [x] Wrap app with URQL provider
- [x] Create sample query for reference letters
- [x] Run codegen and verify typed hooks work
- [x] Add npm scripts for codegen

## Notes
- Use URQL's Next.js integration for SSR support
- Consider caching strategy (default document cache should be fine initially)