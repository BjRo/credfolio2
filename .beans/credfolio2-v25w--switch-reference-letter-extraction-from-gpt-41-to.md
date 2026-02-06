---
# credfolio2-v25w
title: Switch reference letter extraction to Claude Haiku 4.5
status: in-progress
type: task
priority: normal
created_at: 2026-02-06T13:06:21Z
updated_at: 2026-02-06T13:27:47Z
parent: credfolio2-3ram
---

## Summary

The letter_data_extraction step was using GPT-4.1 and taking ~19.5s at $0.023. After benchmarking four models, Claude Haiku 4.5 was chosen as the winner.

### Benchmark Results

| Model | Duration | Cost |
|-------|----------|------|
| GPT-4.1 | 18.69s | $0.023 |
| Claude Sonnet 4.5 | 19.92s | $0.034 |
| GPT-4.1 mini | 26.62s | $0.003 |
| **Claude Haiku 4.5** | **10.85s** | **$0.012** |

## Approach

Changed the default model for reference letter extraction to `anthropic/claude-haiku-4-5-20251001`. This is a configuration change in `config.go` defaults, overridable via the `REFERENCE_EXTRACTION_MODEL` env var.

**Achieved impact:** ~19.5s → ~10.8s (44% faster). Good cost at $0.012.

## Checklist

- [x] Update the default model for reference letter extraction in config (or env var)
- [x] Run a comparison test: benchmarked GPT-4.1, GPT-4.1-mini, Claude Sonnet 4.5, Claude Haiku 4.5 via Braintrust traces
- [ ] Verify structured output schema compliance with the new model
- [ ] Check that testimonial quotes, skill mentions, and discovered skills are extracted accurately
- [x] Update any hardcoded model references (e.g. the `ModelVersion` in document_processing.go:318)

## Definition of Done
- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [ ] All other checklist items above are completed
- [x] Branch pushed and PR created for human review (PR #94)
- [x] Automated code review passed (`@review-backend` — LGTM, no blocking issues)