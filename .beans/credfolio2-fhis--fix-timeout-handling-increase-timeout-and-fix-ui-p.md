---
# credfolio2-fhis
title: 'Fix timeout handling: increase timeout and fix UI polling on failure'
status: in-progress
type: bug
priority: high
created_at: 2026-02-05T17:12:45Z
updated_at: 2026-02-05T17:12:45Z
---

Two related issues when LLM extraction times out:

1. **Backend**: The 120s ResilientProvider timeout covers ALL retries (not per-attempt), which is too short for structured output with slower models like gpt-5-nano
2. **Frontend**: On timeout/failure, the UI doesn't notice and keeps polling indefinitely — either the error status isn't being set correctly in the DB, or the frontend doesn't handle the "failed" status properly

## Investigation Findings

- **Backend**: The failsafe timeout only cancels its derived context, NOT the River job context. So `updateStatusFailed(ctx, ...)` should succeed. However, the error from `updateStatusFailed` was silently ignored (`_ = ...`), so DB update failures were invisible.
- **Frontend**: Both `ResumeUpload` and `ReferenceLetterUploadModal` properly check for `"FAILED"` status and stop polling. The GraphQL converter correctly maps "failed" → "FAILED". No frontend fix needed.
- **Root cause**: Most likely the 120s timeout was too short for gpt-5-nano structured output, causing repeated timeouts. With MaxAttempts: 2 and 120s per attempt, the status bounced "failed" → "processing" → "failed" over ~4 minutes — appearing stuck to the user.

## Checklist
- [x] Investigate: Does updateStatusFailed actually execute on timeout? (check if context is still valid)
- [x] Investigate: How does the frontend poll and does it handle "failed" status?
- [x] Increase resilient timeout to 300s
- [x] Frontend polling already handles failure status correctly — no fix needed
- [x] Log updateStatusFailed errors instead of silently ignoring them
- [x] Tests pass
- [x] Lint passes

## Definition of Done
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review
