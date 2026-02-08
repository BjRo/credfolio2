---
# credfolio2-v27d
title: Improve type system consistency
status: todo
type: task
priority: normal
created_at: 2026-02-08T11:04:07Z
updated_at: 2026-02-08T11:04:07Z
parent: credfolio2-nihn
---

Address type redundancy and correctness issues across the codebase.

## Issues (from @review-backend)

1. **Type Redundancy** - TestimonialRelationship vs AuthorRelationship serve same purpose
2. **Incorrect Type Usage** - ExperienceSource enum used for skills (should be SkillSource)
3. **Missing Domain Types** - Some GraphQL inputs lack domain equivalents

## Impact

- Confusion about which enum to use
- Type system doesn't accurately represent domain
- Tight coupling between GraphQL and domain layer

## Files Affected

- `src/backend/internal/domain/entities.go`
- `src/backend/internal/domain/extraction.go`
- `src/backend/internal/graphql/schema/schema.graphqls`

## Acceptance Criteria

- [ ] Consolidate TestimonialRelationship and AuthorRelationship into single enum
- [ ] Create SkillSource enum and use it correctly
- [ ] Add domain types for key GraphQL inputs
- [ ] Update all usages throughout codebase

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