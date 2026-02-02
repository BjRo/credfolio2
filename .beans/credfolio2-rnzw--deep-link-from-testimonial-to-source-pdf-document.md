---
# credfolio2-rnzw
title: Deep link from testimonial to source PDF document
status: in-progress
type: feature
priority: normal
created_at: 2026-02-02T14:20:08Z
updated_at: 2026-02-02T16:31:35Z
parent: credfolio2-2ex3
---

## Overview

Allow users to navigate directly from a testimonial quote to the specific location in the source PDF reference letter where the quote appears.

## User Story

As a profile viewer, I want to click on a testimonial and be taken to the exact page/location in the PDF where that quote appears, so I can verify the context and read the full reference letter.

## Investigation Areas

- [x] Research PDF deep linking capabilities (page anchors, text fragments)
- [x] Evaluate PDF viewer options that support highlighting/navigation
- [x] Determine if we store page numbers during testimonial extraction
- [x] Consider browser-native PDF viewing vs embedded viewer (pdf.js)
- [x] Explore text fragment URLs (Chrome's `#:~:text=` syntax)

## Research Findings

### PDF Deep Linking
- Browser-native `#page=N` parameter is well-supported across Chrome, Firefox, Edge, Safari
- Works with both presigned URLs and proxy URLs (fragment is client-side)
- No additional dependencies required

### PDF Viewer Options
- **Browser-native**: Simple, works everywhere, supports `#page=N`
- **PDF.js**: More control, can highlight text, but requires significant implementation effort
- **Recommendation**: Start with browser-native `#page=N`, consider PDF.js for future text highlighting

### Current State
- Page numbers are **NOT** currently captured during extraction
- Extraction uses LLM vision to process PDF → raw text → structured data
- The LLM sees the PDF visually and could potentially identify page numbers

### Text Fragment URLs
- `#:~:text=` only works for HTML pages, NOT PDFs
- Not a viable option for this feature

### Chosen Approach
1. **Modify extraction schema**: Add `pageNumber` field to testimonials in LLM output
2. **Add database column**: Add `page_number` to `testimonials` table
3. **Update frontend**: Append `#page=N` to source document URL

## Implementation Checklist

### Backend Changes
- [x] Add `pageNumber` field to testimonial extraction schema (extraction.go)
- [x] Update LLM prompt to extract page numbers (reference_letter_extraction_system.txt)
- [x] Create migration to add `page_number` column to testimonials table
- [x] Update Testimonial entity and repository
- [x] Update GraphQL schema to expose pageNumber
- [x] Update ApplyReferenceLetterValidations to store page numbers

### Frontend Changes
- [x] Update GetTestimonials query to fetch pageNumber
- [x] Modify TestimonialsSection to append #page=N to source URL

## Out of Scope (for now)

- Highlighting the exact text within the PDF (stretch goal)
- Editing or annotating the PDF
- PDF.js embedded viewer

## Definition of Done
- [x] Research completed and approach documented
- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] Visual verification with agent-browser (profile page verified, code changes confirmed)
- [x] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review