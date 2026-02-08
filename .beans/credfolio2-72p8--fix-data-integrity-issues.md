---
# credfolio2-72p8
title: Fix data integrity issues
status: todo
type: task
priority: high
created_at: 2026-02-08T11:04:00Z
updated_at: 2026-02-08T11:04:00Z
parent: credfolio2-nihn
---

Address data integrity and concurrency issues in the backend.

## Critical Issues (from @review-backend)

1. **Race Condition in Author Creation** - TOCTOU vulnerability in findOrCreateAuthor
2. **Missing Transaction Boundaries** - Delete+create cycles lack atomicity

## Impact

- Duplicate authors created when concurrent letter processing
- Partial updates leave database in inconsistent state (orphaned records)
- Data loss risk if later steps fail

## Files Affected

- `src/backend/internal/service/materialization.go`
- `src/backend/internal/repository/postgres/author_repository.go`

## Acceptance Criteria

- [ ] Author creation uses database-level unique constraint + upsert pattern
- [ ] Delete+create operations wrapped in database transaction
- [ ] Tests verify concurrent author creation doesn't create duplicates
- [ ] Tests verify transaction rollback on failure

## Reference

See: /documentation/reviews/2026-02-08-comprehensive-codebase-review.md#critical-issues

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