---
# credfolio2-539r
title: Disable GraphQL Playground in production
status: draft
type: task
created_at: 2026-01-21T15:11:12Z
updated_at: 2026-01-21T15:11:12Z
parent: credfolio2-wxn8
---

Ensure the GraphQL Playground UI is not accessible in production environments.

## Context
The GraphQL Playground is useful for development but exposes the API schema and allows arbitrary query execution, which is a security risk in production.

## Acceptance Criteria
- [ ] Playground endpoint returns 404 or is not mounted in production
- [ ] Environment detection uses `CREDFOLIO_ENV` or similar
- [ ] Playground remains available in development

## Implementation Notes
- Check gqlgen handler configuration options
- Consider using build tags or runtime environment checks