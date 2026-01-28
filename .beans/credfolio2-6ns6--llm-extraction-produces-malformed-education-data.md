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

## Current Approach: Add OpenAI Provider

Adding OpenAI as an alternative LLM provider to test if it produces more reliable extraction results.

### Checklist
- [x] Add `openai-go` dependency to go.mod
- [x] Create `openai.go` provider implementing `domain.LLMProvider`
- [x] Add `OpenAIConfig` to config.go with `OPEN_AI_API_KEY` env var
- [x] Update main.go with provider registry for all available providers
- [x] Restructure prompts: split into system/user with template support
- [x] Add per-use-case provider chains (ProviderChain type)
- [x] Add LLM-based normalization rules to prompts
- [ ] Test extraction with OpenAI model (blocked by network restrictions in devcontainer)
- [x] Run lint and tests

## Previous Open Issues (now blocked on OpenAI testing)

- [ ] LLM extraction is fundamentally unreliable - trying different provider
- [x] May need to use a more capable model - **trying OpenAI**

## Potential Future Approaches (if OpenAI doesn't help)

1. **Retry with validation** - If extraction looks malformed, retry with different prompt
2. **Two-stage extraction** - First understand document structure, then extract fields
3. **Improve text extraction** - The first stage (PDF→text) may be producing poor input
4. **Add confidence thresholds** - Reject extractions below quality threshold

## Pull Request

https://github.com/BjRo/credfolio2/pull/38

## Files Changed

- `src/backend/internal/infrastructure/llm/extraction.go` - Per-use-case provider chains, system/user prompts
- `src/backend/internal/infrastructure/llm/provider_chain.go` - **NEW** ProviderChain and ProviderRegistry types
- `src/backend/internal/infrastructure/llm/openai.go` - **NEW** OpenAI provider implementation
- `src/backend/internal/infrastructure/llm/normalize.go` - Text cleanup functions (kept as fallback)
- `src/backend/internal/infrastructure/llm/normalize_test.go` - Unit tests
- `src/backend/internal/infrastructure/llm/prompts/document_extraction_system.txt` - System prompt with normalization rules
- `src/backend/internal/infrastructure/llm/prompts/document_extraction_user.txt` - User prompt
- `src/backend/internal/infrastructure/llm/prompts/resume_extraction_system.txt` - System prompt with normalization rules
- `src/backend/internal/infrastructure/llm/prompts/resume_extraction_user.txt` - User template
- `src/backend/internal/config/config.go` - Added OpenAI and LLM config
- `src/backend/cmd/server/main.go` - Provider registry and chain setup

## Definition of Done
- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [ ] LLM extraction works reliably - **BLOCKED: Requires architectural changes**
- [x] Branch pushed and PR created for human review
