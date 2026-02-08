---
# credfolio2-offq
title: Fix LLM security vulnerabilities
status: todo
type: task
priority: critical
created_at: 2026-02-08T11:03:02Z
updated_at: 2026-02-08T11:03:02Z
parent: credfolio2-nihn
---

Address critical security vulnerabilities in LLM integration identified in codebase review.

## Critical Issues (from @review-ai)

1. **Prompt Injection** - User text embedded without isolation delimiters
2. **Missing Output Validation** - Extracted fields not validated (XSS/SQL injection risk)
3. **Uncontrolled Text Length** - No size limits on documents (DoS risk)

## Impact

- Security risk: Users can manipulate LLM outputs
- Data quality risk: Invalid data persisted to database
- Cost risk: Large documents cause expensive API calls

## Files Affected

- `src/backend/internal/infrastructure/llm/prompts/*.txt`
- `src/backend/internal/infrastructure/llm/extraction.go`
- `src/backend/internal/job/*_processing.go`

## Acceptance Criteria

- [ ] All prompts use XML tags or markdown code blocks to isolate user content
- [ ] All extracted fields validated and sanitized before database persistence
- [ ] Document size limits enforced (50KB resumes, 100KB letters)
- [ ] Tests verify prompt injection attempts are blocked

## Reference

See: /documentation/reviews/2026-02-08-comprehensive-codebase-review.md#critical-issues-3

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