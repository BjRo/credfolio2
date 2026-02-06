---
# credfolio2-v25w
title: Switch reference letter extraction from GPT-4.1 to GPT-4.1-mini
status: todo
type: task
created_at: 2026-02-06T13:06:21Z
updated_at: 2026-02-06T13:06:21Z
parent: credfolio2-3ram
---

## Summary

The letter_data_extraction step currently uses GPT-4.1 and takes ~19.5s (TTFT: 19.55s) with 3,928 tokens at $0.017. The long time-to-first-token suggests the model is spending excessive time reasoning before generating output.

GPT-4.1-mini is significantly faster for structured output tasks. Since reference letter extraction is fundamentally an **extraction** task (not a reasoning task), a smaller model should maintain quality while dramatically reducing latency.

## Approach

Change the default model for reference letter extraction from `openai/gpt-4o` (which resolves to GPT-4.1 via Braintrust/routing) to `openai/gpt-4.1-mini`. This is a configuration change in `config.go` defaults, overridable via the `REFERENCE_EXTRACTION_MODEL` env var.

**Expected impact:** ~19.5s → ~5-8s. Low quality risk — the task is extraction with structured output, not complex reasoning.

## Checklist

- [ ] Update the default model for reference letter extraction in config (or env var)
- [ ] Run a comparison test: upload the fixture resume/letter with GPT-4.1 vs GPT-4.1-mini and compare extraction quality
- [ ] Verify structured output schema compliance with the new model
- [ ] Check that testimonial quotes, skill mentions, and discovered skills are extracted accurately
- [ ] Update any hardcoded model references (e.g. the `ModelVersion` in document_processing.go:318)

## Definition of Done
- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review
- [ ] Automated code review passed (`@review-backend` and/or `@review-frontend`)