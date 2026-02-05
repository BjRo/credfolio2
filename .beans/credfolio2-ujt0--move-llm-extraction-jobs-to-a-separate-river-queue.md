---
# credfolio2-ujt0
title: Move LLM extraction jobs to a separate River queue with independent worker pool
status: draft
type: task
priority: normal
created_at: 2026-02-05T17:28:21Z
updated_at: 2026-02-05T17:28:21Z
---

Currently all River jobs share a single default queue with MaxWorkers: 10. LLM extraction jobs (resume processing, reference letter processing) can take several minutes with structured output on slower models. If multiple extraction jobs run concurrently, they can exhaust the worker pool and block other (fast) jobs — classic head-of-line blocking.

## Solution
Create a dedicated `llm_extraction` queue with its own worker pool, isolating slow LLM work from the default queue:

```go
Queues: map[string]river.QueueConfig{
    river.QueueDefault: {MaxWorkers: 5},
    "llm_extraction":   {MaxWorkers: 5},
}
```

Then configure the extraction job args to target the `llm_extraction` queue via `InsertOpts.Queue`.

## Checklist
- [ ] Add `llm_extraction` queue to River client config with dedicated MaxWorkers
- [ ] Update `ResumeProcessingArgs.InsertOpts()` to set `Queue: "llm_extraction"`
- [ ] Update `ReferenceLetterProcessingArgs.InsertOpts()` to set `Queue: "llm_extraction"`
- [ ] Consider making queue worker counts configurable via env vars
- [ ] Tests pass
- [ ] Lint passes

## Context
- Related: credfolio2-fhis (timeout handling fix)
- Related: credfolio2-f2na (resilience wrapping)
- Current ResilientProvider timeout: 300s
- Current River worker timeout: 10min (safety net)
- Current MaxAttempts: 2 (1 retry)