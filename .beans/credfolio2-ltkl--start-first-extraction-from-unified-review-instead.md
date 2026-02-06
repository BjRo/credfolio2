---
# credfolio2-ltkl
title: Start first extraction from unified review instead of resume review
status: in-progress
type: task
priority: normal
created_at: 2026-02-06T12:02:52Z
updated_at: 2026-02-06T12:17:50Z
parent: credfolio2-dwid
---

Change the first-time document upload to use the unified upload flow (`UploadFlow` with `ExtractionReview`) instead of the resume-only flow (`ResumeUpload`) that skips review entirely.

### Current behavior
When a user first uploads a document via `/upload-resume`:
1. `ResumeUpload` component calls `uploadResume` mutation
2. Polls backend for extraction completion
3. Auto-redirects to profile page — **no review step**
4. User never sees or approves the extracted data before it's applied

### Desired behavior
The first-time extraction should go through the unified upload flow at `/upload` which uses `UploadFlow`:
1. Upload document → `uploadForDetection` mutation
2. Backend detects content type (career info, testimonial, or both)
3. User reviews detection results and selects what to extract
4. Extraction runs via `processDocument` mutation
5. User reviews extracted data in `ExtractionReview` before importing
6. User explicitly imports results via `importDocumentResults`

This gives users the opportunity to review and confirm extraction results before they're applied to their profile.

### Key files
- `src/frontend/src/app/upload-resume/page.tsx` — current first-time upload entry point (to be changed)
- `src/frontend/src/components/ResumeUpload.tsx` — current resume-only flow (no review)
- `src/frontend/src/app/upload/page.tsx` — unified upload page (target flow)
- `src/frontend/src/components/upload/UploadFlow.tsx` — unified flow orchestrator
- `src/frontend/src/components/upload/ExtractionReview.tsx` — unified review UI

## Checklist
- [x] Homepage redirects to `/upload` instead of showing `ResumeUpload` inline
- [x] Verify homepage redirects to `/upload` with unified flow (visual)
- [ ] Verify full upload→review→import→profile flow (requires LLM API keys)

## Definition of Done
- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] Visual verification with agent-browser (for UI changes)
- [ ] All other checklist items above are completed
- [x] Branch pushed and PR created for human review
- [x] Automated code review passed (`@review-frontend` — no critical findings)