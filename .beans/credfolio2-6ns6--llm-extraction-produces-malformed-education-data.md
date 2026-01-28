---
# credfolio2-6ns6
title: LLM extraction produces malformed education data
status: in-progress
type: bug
priority: normal
created_at: 2026-01-26T11:16:38Z
updated_at: 2026-01-28T15:30:00Z
parent: credfolio2-dwid
---

The resume extraction LLM produces unreliable data. The extraction is inconsistent - sometimes returning correct data, sometimes returning corrupted or empty data.

## Original Problem

Extracted data had spacing artifacts from PDF text extraction:
- "Co lumb ia University" instead of "Columbia University"
- "201 4-01-0 1" instead of "2014-01-01"

## Current Status: PARTIALLY FIXED, STILL UNRELIABLE

The normalization code fixes spacing artifacts when they occur, but the underlying LLM extraction is fundamentally unreliable:

### Observed Issues (from database analysis):

1. **Inconsistent results** - Same resume uploaded multiple times produces different results
2. **Empty extractions** - Some uploads return no education/experience at all
3. **Split entries** - Single education entry split across multiple JSON objects
4. **Field contamination** - Certifications text appearing in GPA/achievements fields
5. **Missing dates** - endDate often missing even when clearly present in source

### Example extraction results from same PDF:

| Upload | Name | Education Count | Experience Count |
|--------|------|-----------------|------------------|
| 50d27ec3 | George Evans | 1 | 1 |
| 381fa367 | George Evans | 2 | 2 |
| cd69d163 | (empty) | 0 | 0 |
| 9dbb0324 | (empty) | 0 | 0 |

## What Was Implemented

### 1. Prompt externalization (DONE)
- Prompts moved to `src/backend/internal/infrastructure/llm/prompts/`
- Uses Go's `//go:embed` for compile-time embedding
- Easier to iterate on prompts without code changes

### 2. Post-extraction normalization (DONE)
- `normalize.go` with text cleanup functions
- Fixes spacing artifacts in institution names, dates, etc.
- Unit tests for normalization logic

### 3. Improved extraction prompt (TRIED, DID NOT FIX)
- Added explicit rules for field placement
- "ONE ENTRY PER ITEM" rule
- "CERTIFICATIONS ARE NOT EDUCATION" rule
- Still produces unreliable results

## Open Issues

- [ ] LLM extraction is fundamentally unreliable - needs architectural change
- [ ] Consider retry logic with validation
- [ ] Consider different extraction approach (e.g., chain-of-thought)
- [ ] May need to use a more capable model
- [ ] Document text extraction (first stage) may be the root cause

## Potential Future Approaches

1. **Retry with validation** - If extraction looks malformed, retry with different prompt
2. **Two-stage extraction** - First understand document structure, then extract fields
3. **Model upgrade** - Use a more capable model for structured extraction
4. **Improve text extraction** - The first stage (PDF→text) may be producing poor input
5. **Add confidence thresholds** - Reject extractions below quality threshold

## Pull Request

https://github.com/BjRo/credfolio2/pull/38

## Files Changed

- `src/backend/internal/infrastructure/llm/extraction.go` - Added normalization, embed prompts
- `src/backend/internal/infrastructure/llm/normalize.go` - New file with normalization functions
- `src/backend/internal/infrastructure/llm/normalize_test.go` - Unit tests
- `src/backend/internal/infrastructure/llm/prompts/document_extraction.txt` - Document extraction prompt
- `src/backend/internal/infrastructure/llm/prompts/resume_extraction.txt` - Resume extraction prompt

## Definition of Done
- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [ ] LLM extraction works reliably - **BLOCKED: Requires architectural changes**
- [x] Branch pushed and PR created for human review
