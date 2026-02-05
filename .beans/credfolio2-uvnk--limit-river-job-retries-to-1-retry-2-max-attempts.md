---
# credfolio2-uvnk
title: Limit River job retries to 1 retry (2 max attempts)
status: in-progress
type: task
priority: normal
created_at: 2026-02-05T16:52:56Z
updated_at: 2026-02-05T16:52:56Z
---

River jobs currently use the default 25 max attempts. For now, limit to 2 max attempts (1 initial + 1 retry) to avoid repeatedly hammering LLM providers on persistent failures.

Fix: override MaxAttempts on each job worker to return 2.

## Checklist
- [x] Set MaxAttempts to 2 on ResumeProcessingWorker
- [x] Set MaxAttempts to 2 on ReferenceLetterProcessingWorker
- [x] Set MaxAttempts to 2 on DocumentProcessingWorker
- [x] Tests pass
- [x] Lint passes

## Definition of Done
- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review
