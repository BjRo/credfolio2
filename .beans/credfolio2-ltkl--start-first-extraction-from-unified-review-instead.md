---
# credfolio2-ltkl
title: Start first extraction from unified review instead of resume review
status: todo
type: task
priority: normal
created_at: 2026-02-06T12:02:52Z
updated_at: 2026-02-06T12:02:52Z
parent: credfolio2-dwid
---

Change the reference letter extraction flow so that the first extraction starts from the unified review UI (ExtractionReview) rather than the resume-specific review flow.

### Current behavior
When a reference letter is uploaded from the profile page via `ReferenceLetterUploadModal`:
1. Modal uploads file → backend extracts data → returns reference letter ID
2. Redirects to `/profile/{id}/reference-letters/{letterID}/preview` (ValidationPreviewPage)
3. ValidationPreviewPage shows corroborations, testimonials, discovered skills
4. User applies validations via `ApplyReferenceLetterValidations` mutation

### Desired behavior
The first extraction for a reference letter should go through the unified review flow (`ExtractionReview` component in `src/frontend/src/components/upload/ExtractionReview.tsx`) instead of the separate resume-specific `ValidationPreviewPage`. This consolidates the review experience so all document types use the same extraction review UI.

### Key files
- `src/frontend/src/components/upload/ExtractionReview.tsx` — unified review UI
- `src/frontend/src/components/profile/ReferenceLetterUploadModal.tsx` — current modal flow
- `src/frontend/src/app/profile/[id]/reference-letters/[referenceLetterID]/preview/page.tsx` — current reference letter preview
- `src/frontend/src/app/profile/[id]/page.tsx` — profile page handler for upload success

## Checklist
- [ ] Route reference letter uploads through the unified review flow
- [ ] Ensure ExtractionReview handles reference-letter-only data (no career info section)
- [ ] Update redirect after reference letter upload to point to unified review
- [ ] Verify import flow works correctly for reference letters through unified review
- [ ] Consider whether the standalone ValidationPreviewPage is still needed or can be removed

## Definition of Done
- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Visual verification with agent-browser (for UI changes)
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review
- [ ] Automated code review passed (`@review-backend` and/or `@review-frontend`)