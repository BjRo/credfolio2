---
# credfolio2-35s5
title: Fix performance bottlenecks
status: todo
type: task
priority: high
created_at: 2026-02-08T11:03:47Z
updated_at: 2026-02-08T11:03:47Z
parent: credfolio2-nihn
---

Address critical performance issues identified in backend and frontend code.

## Critical Issues

### Backend (from @review-backend)
1. **N+1 Query Pattern** - Validation count resolvers query per skill/experience

### Frontend (from @review-frontend)
2. **Home Page Waterfall** - Client-side fetch before redirect
3. **File Upload Bypass** - Uses XMLHttpRequest instead of urql
4. **Testimonials Waterfall** - Sequential queries instead of single query

## Impact

- Poor performance at scale (100 skills = 100 queries)
- Slow page loads and navigation
- Inconsistent error handling (file upload)

## Files Affected

- `src/backend/internal/graphql/resolver/schema.resolvers.go`
- `src/frontend/src/app/page.tsx`
- `src/frontend/src/components/upload/document-upload.tsx`
- `src/frontend/src/components/profile/testimonials-section.tsx`

## Acceptance Criteria

- [ ] N+1 queries resolved with dataloader or eager loading
- [ ] Home page uses Server Component with `redirect()`
- [ ] File upload uses proper GraphQL mutation
- [ ] Testimonials fetched in single query with fragments

## Reference

See: /documentation/reviews/2026-02-08-comprehensive-codebase-review.md

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