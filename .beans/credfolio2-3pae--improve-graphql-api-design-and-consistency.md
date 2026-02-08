---
# credfolio2-3pae
title: Improve GraphQL API design and consistency
status: todo
type: task
priority: high
created_at: 2026-02-08T11:14:40Z
updated_at: 2026-02-08T11:14:40Z
parent: credfolio2-nihn
---

Address GraphQL schema design issues including naming consistency, API surface consolidation, and DoS risks.

## Issues (from @review-backend and @review-frontend)

1. **Unbounded Arrays (DoS Risk)** - `testimonials`, `skillValidations`, `experienceValidations` arrays have no pagination/limits
2. **Naming Consistency** - Inconsistent conventions across queries, mutations, and input types
3. **API Surface Consolidation** - Mutations that could be combined or simplified
4. **Internal Implementation Leakage** - Some queries expose internal structure unnecessarily

## Impact

- **Security**: Unbounded arrays allow DoS attacks (request profile with 100k testimonials)
- **Developer Experience**: Inconsistent naming makes API harder to learn and use
- **Maintenance**: Redundant mutations increase API surface and test burden

## Files Affected

- `src/backend/internal/graphql/schema/schema.graphqls`
- `src/backend/internal/graphql/resolver/schema.resolvers.go`
- `src/frontend/src/graphql/queries.graphql`
- `src/frontend/src/graphql/mutations.graphql`

## Acceptance Criteria

- [ ] All array fields have pagination parameters (`first`, `after`) or reasonable limits
- [ ] Naming convention documented and applied consistently (queries, mutations, inputs, enums)
- [ ] API surface audit completed with consolidation opportunities documented
- [ ] Breaking changes planned with migration path for frontend

## Reference

See: /documentation/reviews/2026-02-08-comprehensive-codebase-review.md#warnings-4

## Definition of Done
- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Visual verification via `@qa` subagent (via Task tool, for UI changes)
- [ ] ADR written via `/decision` skill (if new dependencies, patterns, or architectural changes were introduced)
- [ ] All other checklist items above are completed
- [ ] Branch pushed to remote
- [ ] PR created for human review
- [ ] Automated code review passed via `@review-backend`, `@review-frontend`, and/or `@review-ai` (for LLM changes) subagents (via Task tool)