---
# credfolio2-1g6t
title: Optimize LLM invocation costs
status: todo
type: task
priority: normal
created_at: 2026-02-06T12:05:19Z
updated_at: 2026-02-06T12:34:09Z
parent: credfolio2-3ram
---

After completing the pipeline speed optimizations (credfolio2-5j3j), review and optimize the cost of LLM invocations across the document processing pipeline.

## Context

The backend makes several LLM calls during document processing (detection, text extraction, resume extraction, reference letter extraction). Each call has token costs that add up, especially for longer documents. This bean focuses on reducing those costs without sacrificing quality.

## Areas to investigate

1. **Model selection** — Are we using the most cost-effective model for each task? Lower-tier models may be sufficient for simpler tasks (e.g. document type detection) while reserving more capable models for complex extraction.

2. **Prompt optimization** — Review prompt lengths and structure. Are prompts unnecessarily verbose? Can system prompts be shortened? Are we sending more context than needed?

3. **Token usage** — Audit input/output token counts per LLM call. Identify which calls are the most expensive and whether their cost is justified by the quality of output.

4. **Caching opportunities** — Are there cases where the same or similar LLM calls are made repeatedly (e.g. re-processing the same document)? Could results be cached?

5. **Batching / consolidation** — Can multiple small LLM calls be consolidated into fewer calls with structured output? (e.g. detect + extract in a single call)

## Dependencies

- Should be done **after** credfolio2-5j3j (pipeline speed optimization), since that work eliminates duplicate LLM calls and changes the extraction flow

## Checklist
- [ ] Audit all LLM invocations in the backend (detection, extraction, etc.)
- [ ] Document current model, token usage, and estimated cost per call
- [ ] Identify cost reduction opportunities (model downgrades, prompt trimming, caching)
- [ ] Implement changes and measure cost impact
- [ ] Verify extraction quality is maintained after optimizations

## Definition of Done
- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review
- [ ] Automated code review passed (`@review-backend` and/or `@review-frontend`)