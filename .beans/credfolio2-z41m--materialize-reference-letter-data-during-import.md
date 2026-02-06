---
# credfolio2-z41m
title: Materialize reference letter data during import
status: in-progress
type: task
priority: high
created_at: 2026-02-05T23:06:50Z
updated_at: 2026-02-06T08:10:25Z
parent: credfolio2-3ram
---

## Summary

When a document contains both career info and testimonial content, both extractions run successfully (resume + reference letter). However, at import time, only resume data gets materialized into profile tables. The reference letter extracted data sits unused in `reference_letters.extracted_data` JSONB.

## Root Cause

The `ImportDocumentResults` resolver (`schema.resolvers.go:862-918`) only has a resume materialization block. After it, the code jumps to "get profile and return" — there's no corresponding block for reference letters.

Additionally:
- `MaterializationService` only has `MaterializeResumeData` — no `MaterializeReferenceLetterData`
- `ImportedCount` GraphQL type only has resume-oriented fields (`experiences`, `educations`, `skills`) — no `testimonials`
- The profile `testimonials` table exists but is never populated from extraction

## What needs to happen

### Backend
- [x] Add `MaterializeReferenceLetterData` to `MaterializationService` — create testimonial rows from `ReferenceLetterExtractedData` (author info, testimonial quotes, skill mentions)
- [x] Add `testimonials` field to `ImportedCount` GraphQL type
- [x] Add reference letter materialization block to `ImportDocumentResults` resolver (after resume block)
- [x] Add resolver tests for import with reference letter data

### Frontend
- [x] Update `ImportedCount` display to show testimonial count when present
- [x] Verify testimonials section on profile page renders imported data

### Definition of Done
- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] Visual verification with agent-browser (for UI changes)
- [x] All other checklist items above are completed
- [x] Branch pushed and PR created for human review