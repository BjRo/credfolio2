---
# credfolio2-6ns6
title: LLM extraction produces malformed education data
status: in-progress
type: bug
priority: normal
created_at: 2026-01-26T11:16:38Z
updated_at: 2026-01-28T14:54:24Z
parent: credfolio2-dwid
---

The resume extraction LLM sometimes produces corrupted data with unwanted spaces in text fields. Examples from the uploaded PDF "CV_TEMPLATE_0004.pdf":

**Source data (PDF):**
- Institution: "Columbia University"
- Degree: "Bachelor of Science: Computer Information Systems"
- Year: 2018

**Extracted data (corrupted):**
- Institution: "Co lumb ia University" (spaces breaking up words)
- Date: "201 4-01-0 1 - Present" (spaces in date)
- Field: "Comput er Information Systems" (space in word)

## Root Cause Analysis

The text extraction stage (ExtractText using LLM vision) appears to produce garbled text from the PDF due to how the PDF renders characters with individual positioning.

## Solution Implemented

1. **Extracted prompts to dedicated files** using Go's `//go:embed` directive:
   - `src/backend/internal/infrastructure/llm/prompts/document_extraction.txt`
   - `src/backend/internal/infrastructure/llm/prompts/resume_extraction.txt`

2. **Added text normalization** in `normalize.go`:
   - `NormalizeSpacedText()` - removes spurious spaces from OCR artifacts
   - `NormalizeDate()` - validates and cleans ISO date format
   - `HasExcessiveSpacing()` - detects text with spacing artifacts

3. **Post-extraction normalization** in `extraction.go`:
   - All text fields are normalized after LLM extraction
   - Malformed dates are converted to nil
   - Institution/company names are cleaned up

## Checklist

- [x] Extract prompts into dedicated files in `src/backend/internal/infrastructure/llm/prompts/`
  - [x] Create `prompts/` directory
  - [x] Create `prompts/document_extraction.txt` for the default document extraction prompt
  - [x] Create `prompts/resume_extraction.txt` for the resume structured extraction prompt
  - [x] Update extraction.go to load prompts from files (using go:embed)
- [x] Add text normalization to clean up extracted text before structured extraction
  - [x] Add function to normalize whitespace in extracted text
  - [x] Remove spurious spaces within words (OCR artifacts)
  - [x] Normalize date formats
- [x] Add post-extraction validation to reject malformed data
  - [x] Validate date format (YYYY-MM-DD)
  - [x] Validate institution names don't have excessive spaces
  - [x] ~~Add confidence threshold for data quality~~ (Not needed - normalization handles this)
- [x] Write tests for the new validation and normalization logic
- [ ] Manual test with fixture resume

## Definition of Done
- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [ ] Visual verification with agent-browser (for UI changes)
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review
