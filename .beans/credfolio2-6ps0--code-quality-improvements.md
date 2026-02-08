---
# credfolio2-6ps0
title: Code quality improvements
status: todo
type: task
priority: low
created_at: 2026-02-08T11:04:19Z
updated_at: 2026-02-08T11:04:19Z
parent: credfolio2-nihn
---

Address code quality issues like dead code, duplication, and missing patterns.

## Issues

### Backend (from @review-backend)
1. **Dead Code** - Unused getProfileSkillsContext function
2. **Duplicate Functions** - mapAuthorRelationship duplicated
3. **Silent Error Suppression** - Duplicate suppression hides bugs

### Frontend (from @review-frontend)
4. **Query Duplication** - No GraphQL fragments, 160+ lines duplicated
5. **Missing Optimizations** - groupExperiencesByCompany not memoized

### LLM (from @review-ai)
6. **Duplicate Text Extraction** - Same text extracted multiple times

## Impact

- Code bloat and maintenance burden
- Hidden bugs from silent error handling
- Larger bundle sizes and slower queries

## Files Affected

- `src/backend/internal/service/materialization.go`
- `src/frontend/src/graphql/queries.graphql`
- `src/frontend/src/components/profile/*.tsx`

## Acceptance Criteria

- [ ] Dead code removed
- [ ] Duplicate functions consolidated
- [ ] GraphQL fragments created for common field selections
- [ ] Expensive computations memoized
- [ ] Text extraction optimized to avoid duplicates

## Reference

See: /documentation/reviews/2026-02-08-comprehensive-codebase-review.md#suggestions

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