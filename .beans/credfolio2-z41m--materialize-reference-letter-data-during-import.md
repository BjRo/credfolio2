---
# credfolio2-z41m
title: Materialize reference letter data during import
status: todo
type: feature
priority: high
created_at: 2026-02-05T23:06:50Z
updated_at: 2026-02-05T23:06:50Z
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
- [ ] Add `MaterializeReferenceLetterData` to `MaterializationService` — create testimonial rows from `ReferenceLetterExtractedData` (author info, testimonial quotes, skill mentions)
- [ ] Add `testimonials` field to `ImportedCount` GraphQL type
- [ ] Add reference letter materialization block to `ImportDocumentResults` resolver (after resume block)
- [ ] Add resolver tests for import with reference letter data

### Frontend
- [ ] Update `ImportedCount` display to show testimonial count when present
- [ ] Verify testimonials section on profile page renders imported data

### Definition of Done
- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Visual verification with agent-browser (for UI changes)
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review