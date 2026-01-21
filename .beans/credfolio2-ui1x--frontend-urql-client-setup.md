---
# credfolio2-ui1x
title: 'Frontend: URQL client setup'
status: todo
type: feature
priority: normal
created_at: 2026-01-21T14:28:04Z
updated_at: 2026-01-21T14:28:41Z
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
- [ ] Add URQL dependencies (`@urql/core`, `@urql/next`, `urql`)
- [ ] Add GraphQL codegen dependencies (`@graphql-codegen/cli`, etc.)
- [ ] Configure codegen.ts for URQL plugin
- [ ] Copy schema from backend or configure introspection
- [ ] Create URQL client provider
- [ ] Wrap app with URQL provider
- [ ] Create sample query for reference letters
- [ ] Run codegen and verify typed hooks work
- [ ] Add npm scripts for codegen

## Notes
- Use URQL's Next.js integration for SSR support
- Consider caching strategy (default document cache should be fine initially)