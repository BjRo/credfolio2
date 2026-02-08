---
# credfolio2-5od0
title: Improve LLM extraction quality and cost efficiency
status: todo
type: task
priority: normal
created_at: 2026-02-08T11:10:03Z
updated_at: 2026-02-08T11:10:03Z
parent: credfolio2-nihn
---

Address LLM accuracy, data quality, and cost optimization opportunities identified in codebase review.

## Important Issues (from @review-ai)

1. **Resume Summary Synthesis** - LLM generates summary instead of extracting (hallucination risk)
2. **Unknown Author Acceptance** - System accepts "Unknown" as valid author name (data quality issue)
3. **JSON Cleanup Masking Quality** - Aggressive cleanup hides LLM output quality problems
4. **Duplicate Text Extraction** - Same text extracted multiple times (cost/performance waste)

## Optimization Opportunities

5. **Use Haiku for Detection** - Document classification doesn't need Sonnet (10x cost reduction)
6. **Prompt Versioning** - No system for A/B testing prompt improvements
7. **Enhanced Extraction Metadata** - Missing token counts, duration tracking, model versions
8. **Per-Task Timeouts** - All tasks use same timeout regardless of complexity

## Impact

- **Accuracy**: Hallucination risk in summaries, poor data quality from unknown authors
- **Cost**: Unnecessarily expensive model for simple classification, duplicate processing
- **Observability**: Can't track prompt effectiveness or LLM performance over time

## Files Affected

- `src/backend/internal/infrastructure/llm/extraction.go`
- `src/backend/internal/infrastructure/llm/prompts/resume_extraction.txt`
- `src/backend/internal/infrastructure/llm/prompts/letter_extraction.txt`
- `src/backend/internal/job/document_detection.go`
- `src/backend/internal/job/resume_processing.go`
- `src/backend/internal/job/reference_letter_processing.go`
- `src/backend/internal/domain/extraction.go`

## Acceptance Criteria

### Data Quality
- [ ] Resume summaries extracted from text, not synthesized by LLM
- [ ] "Unknown" authors rejected, require actual name extraction
- [ ] JSON cleanup logs warnings when aggressive fixes needed (indicates prompt issues)
- [ ] Text extraction deduplicated to avoid redundant LLM calls

### Cost Optimization
- [ ] Document detection uses Haiku instead of Sonnet (10x cost reduction)
- [ ] Prompt versions tracked in code and logs
- [ ] Extraction metadata includes: tokens used, duration, model version

### Operational Improvements
- [ ] Per-task timeout configuration (detection: 30s, extraction: 2min)
- [ ] Dashboard/logs show prompt effectiveness metrics

## Reference

See: /documentation/reviews/2026-02-08-comprehensive-codebase-review.md#important-issues-4

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