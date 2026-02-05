---
# credfolio2-u1eh
title: Unified document processing orchestrator
status: draft
type: task
priority: high
created_at: 2026-02-05T18:02:02Z
updated_at: 2026-02-05T18:03:08Z
parent: credfolio2-3ram
blocking:
    - credfolio2-sxcf
---

## Summary

Create a unified document processing pipeline that accepts detection results + user extraction preferences and orchestrates running the appropriate extractor(s) on the document. This replaces the need for separate `uploadResume` and `uploadFile` mutations when using the unified flow.

## Background

Currently, resume and reference letter processing are separate async jobs (`ResumeProcessingWorker`, `ReferenceLetterProcessingWorker`). This task creates an orchestrator that:
1. Takes a previously uploaded & detected document
2. Runs the selected extractors based on user preferences
3. Returns results for user review before profile import

## Dependencies

- Requires: credfolio2-4h8a (Document content detection service) — needs the stored file + extracted text

## Checklist

### Backend Service
- [ ] Create unified processing orchestrator service
  - Accepts: file/document ID, extraction preferences (extract_career_info, extract_testimonial)
  - Reuses extracted text from detection step (avoid re-extracting)
  - Runs `ExtractResumeData` if career info selected
  - Runs `ExtractLetterData` if testimonial selected
  - Can run both sequentially on same document
  - Returns combined extraction results
- [ ] Decide sync vs async processing strategy
  - Option A: Synchronous (simpler, but blocks ~15-30s)
  - Option B: Async job with polling (consistent with existing pattern)
  - Recommendation: Async with polling, matching existing UX patterns

### Data Model
- [ ] Consider a unified `Document` entity that can hold both resume and reference letter extraction data
  - Or: create both a Resume and ReferenceLetter record as needed, linking to same File
- [ ] Store extraction preferences alongside the document for audit trail

### GraphQL API
- [ ] Add `processDocument(fileId: ID!, preferences: DocumentProcessingInput!)` mutation
  - `DocumentProcessingInput`: `extractCareerInfo: Boolean!, extractTestimonial: Boolean!`
  - Returns processing status (job ID or inline results)
- [ ] Add `DocumentProcessingResult` type combining possible extraction outputs
  - Career info: positions, skills, education (same as existing `ResumeExtractedData`)
  - Testimonial: author, quotes, skill mentions (same as existing `ExtractedLetterData`)
- [ ] Add query to poll processing status if async

### Profile Import
- [ ] Add `importDocumentResults(userId: ID!, input: DocumentImportInput!)` mutation
  - Merges career info into profile (reuse existing materialization logic)
  - Applies testimonial data (reuse existing `applyReferenceLetterValidations` logic)
  - Handles deduplication: skills matched by normalized name, experiences by company+role+dates
- [ ] Ensure both data types can be imported from a single document

### Feedback Logging
- [ ] Add simple feedback storage (table or structured log) for:
  - Detection corrections ("user said this was just a resume")
  - Extraction quality issues ("report extraction issue" with free text)
- [ ] Add `reportDocumentFeedback(documentId: ID!, feedback: DocumentFeedbackInput!)` mutation

### Testing
- [ ] Unit tests for orchestrator logic
- [ ] Integration tests for processing + import flow
- [ ] Test edge cases: both extractors, single extractor, failed extraction

## Design Notes

- Reuse as much existing extraction logic as possible — don't duplicate `ExtractResumeData` or `ExtractLetterData`
- The orchestrator is primarily a coordination layer, not new extraction logic
- Consider whether a new River job type is needed or if existing workers can be reused

## Definition of Done
- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] All checklist items above are completed
- [ ] Branch pushed and PR created for human review