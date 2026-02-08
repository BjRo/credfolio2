---
# credfolio2-vogi
title: Research VCR-like setup for frontend testing
status: todo
type: task
priority: low
created_at: 2026-02-08T11:15:15Z
updated_at: 2026-02-08T11:15:15Z
parent: credfolio2-nihn
---

Research and prototype VCR-like (record/replay) setup for frontend testing to enable on-demand testing with real GraphQL responses.

## Goal

Enable frontend tests to record real GraphQL responses and replay them deterministically, similar to Ruby's VCR gem or Polly.js for HTTP.

## Benefits

- **Realistic Tests**: Use real backend responses instead of mocks
- **Easier Maintenance**: Record once, replay many times
- **Regression Detection**: Detect when backend responses change unexpectedly
- **On-Demand Testing**: Run frontend tests without backend running

## Research Questions

1. **Tool Selection**: Which library/approach?
   - MSW (Mock Service Worker) with recording mode
   - Polly.js
   - Custom urql exchange for recording
   - graphql-codegen with fixture generation

2. **Recording Strategy**: How to capture responses?
   - Manual recording via test suite
   - Record from real user sessions (dev environment)
   - Generate from GraphQL schema + faker

3. **Storage Format**: Where/how to store recordings?
   - JSON files in `__fixtures__/` directory
   - Per-test fixtures vs shared fixture library
   - How to handle dynamic data (timestamps, IDs)

4. **Replay Mechanism**: How to inject recordings?
   - urql custom exchange
   - MSW request handlers
   - Test setup utilities

5. **Maintenance**: How to keep recordings fresh?
   - Automated re-recording on backend changes
   - Version tracking for fixtures
   - Diff tool for response changes

## Acceptance Criteria

- [ ] Research document comparing 3+ approaches with pros/cons
- [ ] Proof-of-concept implementation for one component test
- [ ] Decision on recommended approach documented
- [ ] Migration plan for existing tests (if approach is viable)
- [ ] Cost/benefit analysis vs current mocking strategy

## Files to Explore

- `src/frontend/src/test/mocks/urql-next.tsx` (current mocking)
- `src/frontend/vitest.config.ts` (test setup)
- Existing component tests

## Reference

New research area (not from codebase review, but related to test quality improvements)

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