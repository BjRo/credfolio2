---
# credfolio2-dhie
title: Fix polling error handling in frontend
status: todo
type: task
priority: normal
created_at: 2026-02-08T11:15:01Z
updated_at: 2026-02-08T11:15:01Z
parent: credfolio2-nihn
---

Add consecutive error detection to polling logic to prevent infinite polling on network failures.

## Problem (from @review-frontend)

Current polling implementation lacks consecutive error detection, causing it to poll forever when the backend is down or network is unavailable. This creates poor user experience and wastes resources.

## Current Behavior

When backend is unavailable:
- Polling continues indefinitely
- User sees perpetual loading state
- No error message or recovery option
- Browser makes requests every few seconds forever

## Impact

- **User Experience**: Users stuck in loading state with no feedback
- **Resource Waste**: Unnecessary network requests and CPU usage
- **Battery Drain**: Mobile devices affected by continuous polling

## Files Affected

- `src/frontend/src/components/upload/extraction-progress.tsx`
- `src/frontend/src/components/upload/detection-progress.tsx`
- Any other components using polling patterns

## Acceptance Criteria

- [ ] Polling stops after N consecutive errors (e.g., 5 failures)
- [ ] User sees clear error message when polling fails
- [ ] Retry button provided for manual retry
- [ ] Error state distinguishes between network failure and backend error
- [ ] Tests verify polling stops after consecutive failures
- [ ] Tests verify retry button works correctly

## Proposed Implementation

```typescript
const [consecutiveErrors, setConsecutiveErrors] = useState(0);
const MAX_CONSECUTIVE_ERRORS = 5;

useEffect(() => {
  if (consecutiveErrors >= MAX_CONSECUTIVE_ERRORS) {
    // Stop polling, show error
    return;
  }
  // Continue polling logic
}, [consecutiveErrors]);
```

## Reference

See: /documentation/reviews/2026-02-08-comprehensive-codebase-review.md#warnings-3

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