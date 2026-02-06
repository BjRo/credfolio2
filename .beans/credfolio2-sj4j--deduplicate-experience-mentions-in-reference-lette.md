---
# credfolio2-sj4j
title: Deduplicate experience mentions in reference letter extraction
status: scrapped
type: task
priority: normal
created_at: 2026-02-06T09:00:04Z
updated_at: 2026-02-06T13:09:16Z
parent: credfolio2-3ram
---

## Summary

When extracting reference letter data, the LLM produces duplicate experience mentions for the same role/company. If the letter discusses "Senior Engineering Manager, Infrastructure at Wellfound" in three different paragraphs, the extraction returns three separate experience mention entries — one per quote/paragraph.

## Current Behavior

The `experienceMentions` array in `ExtractedLetterData` contains one entry per quote that references a role, not one entry per unique role. This makes the "Experience Mentions" section in the ExtractionReview UI noisy and repetitive.

## Expected Behavior

Each unique role/company combination should appear once, with quotes consolidated (e.g., as an array of quotes, or the most representative quote selected).

## Root Cause

The reference letter extraction LLM prompt (`reference_letter_extraction_system.txt`) doesn't instruct the model to deduplicate experience mentions by role+company. It simply says "Extract references to specific roles or companies the candidate held" without guidance on consolidation.

## Impact

Display-only — experience mentions are not materialized into profile tables. But the noisy display degrades the extraction review UX.

## Possible Approaches

1. **Prompt fix** — Update the extraction prompt to instruct the LLM to output one entry per unique role, picking the most representative quote
2. **Post-processing** — Deduplicate in the backend after extraction (group by role+company, keep first or best quote)
3. **Both** — Prompt for deduplication + backend safety net

## Relevant Files

- `src/backend/internal/infrastructure/llm/prompts/reference_letter_extraction_system.txt` — extraction prompt
- `src/backend/internal/infrastructure/llm/extraction.go` — extraction schema definition
- `src/backend/internal/domain/extraction.go` — `ExtractedExperienceMention` type
- `src/frontend/src/components/upload/ExtractionReview.tsx` — display in TestimonialSection