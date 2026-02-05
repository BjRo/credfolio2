---
# credfolio2-fhis
title: 'Fix timeout handling: increase timeout and fix UI polling on failure'
status: completed
type: bug
priority: high
created_at: 2026-02-05T17:12:45Z
updated_at: 2026-02-05T17:18:49Z
---

Two related issues when LLM extraction times out:

1. **Backend**: The 120s ResilientProvider timeout covers ALL retries (not per-attempt), which is too short for structured output with slower models like gpt-5-nano
2. **Frontend**: On timeout/failure, the UI doesn't notice and keeps polling indefinitely — either the error status isn't being set correctly in the DB, or the frontend doesn't handle the "failed" status properly

## Investigation Findings

- **Root cause found**: River's `WorkerDefaults.Timeout()` returns 0, which inherits the client-level `JobTimeoutDefault = 1 * time.Minute`. This 60s River timeout cancels the entire job context — killing BOTH the LLM call AND the subsequent `updateStatusFailed` DB call. That's why the status never gets set to "failed" and the UI polls forever.
- **Frontend**: Both `ResumeUpload` and `ReferenceLetterUploadModal` properly check for `"FAILED"` status and stop polling. No frontend fix needed.
- The `updateStatusFailed` errors were silently ignored — now logged.

## Fixes Applied
1. Override `Timeout()` on resume and reference letter workers to return -1 (no River-level timeout)
2. Increased ResilientProvider timeout from 120s to 300s
3. Changed `updateStatusFailed` to log DB errors instead of silently ignoring

## Checklist
- [x] Investigate: Does updateStatusFailed actually execute on timeout? (River's 60s kills the context first)
- [x] Investigate: How does the frontend poll and does it handle "failed" status? (Yes, correctly)
- [x] Increase resilient timeout to 300s
- [x] Override worker Timeout() to disable River's 60s default
- [x] Frontend polling already handles failure status correctly — no fix needed
- [x] Log updateStatusFailed errors instead of silently ignoring them
- [x] Tests pass
- [x] Lint passes

## Definition of Done
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] All other checklist items above are completed
- [x] Branch pushed and PR created for human review
