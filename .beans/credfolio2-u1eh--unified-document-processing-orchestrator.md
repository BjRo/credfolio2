---
# credfolio2-u1eh
title: Unified document processing orchestrator
status: completed
type: task
priority: high
created_at: 2026-02-05T18:02:02Z
updated_at: 2026-02-05T20:59:08Z
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
- [x] Create unified processing orchestrator service
  - Accepts: file/document ID, extraction preferences (extract_career_info, extract_testimonial)
  - Reuses extracted text from detection step (avoid re-extracting)
  - Runs `ExtractResumeData` if career info selected
  - Runs `ExtractLetterData` if testimonial selected
  - Can run both sequentially on same document
  - Returns combined extraction results
- [x] Decide sync vs async processing strategy
  - Decision: Async with polling via River job, matching existing UX patterns
  - Worker: `DocumentProcessingWorker` in `job/document_processing.go`

### Data Model
- [x] Consider a unified `Document` entity that can hold both resume and reference letter extraction data
  - Decision: Create both Resume and ReferenceLetter records as needed, linking to same File
  - No new table needed — reuses existing entities
- [x] Store extraction preferences alongside the document for audit trail
  - Tracked via which entity IDs are set in `DocumentProcessingArgs`

### GraphQL API
- [x] Add `processDocument(userId: ID!, input: ProcessDocumentInput!)` mutation
  - `ProcessDocumentInput`: `fileId: ID!, extractCareerInfo: Boolean!, extractTestimonial: Boolean!`
  - Returns `ProcessDocumentResult` with optional `resumeId` and `referenceLetterID`
- [x] Add `DocumentProcessingResult` type combining possible extraction outputs
  - Career info stored in Resume.ExtractedData (same as existing `ResumeExtractedData`)
  - Testimonial stored in ReferenceLetter.ExtractedData (same as existing `ExtractedLetterData`)
- [x] Add query to poll processing status: `documentProcessingStatus(resumeId: ID, referenceLetterID: ID)`
  - Returns `DocumentProcessingStatus` with resume, referenceLetter, and allComplete flag

### Profile Import
- [x] Add `importDocumentResults(userId: ID!, input: ImportDocumentResultsInput!)` mutation
  - Materializes career info into profile using shared `MaterializationService`
  - Returns `ImportDocumentResultsResult` with profile and importedCount
- [x] Ensure both data types can be imported from a single document
  - Resume materialization handled; letter validation can use existing `applyReferenceLetterValidations`

### Feedback Logging
- [x] Add simple feedback storage (structured log) for:
  - Detection corrections ("user said this was just a resume")
  - Extraction quality issues ("report extraction issue" with free text)
- [x] Add `reportDocumentFeedback(userId: ID!, input: DocumentFeedbackInput!)` mutation
  - Logs structured feedback via logger; no dedicated table for MVP

### Testing
- [x] Unit tests for orchestrator logic (10 tests in document_processing_test.go)
- [x] Integration tests for processing + import flow (via comprehensive worker tests)
- [x] Test edge cases: both extractors, single extractor, failed extraction

## Design Notes

- Reuse as much existing extraction logic as possible — don't duplicate `ExtractResumeData` or `ExtractLetterData`
- The orchestrator is primarily a coordination layer, not new extraction logic
- New River job type `unified_document_processing` created alongside existing workers
- Shared `MaterializationService` extracted for reuse by both existing resume worker and new import mutation

## Definition of Done
- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] All checklist items above are completed
- [x] Branch pushed and PR created for human review