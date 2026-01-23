---
# credfolio2-6oza
title: Resume Upload & Initial Extraction
status: draft
type: epic
created_at: 2026-01-23T16:26:55Z
updated_at: 2026-01-23T16:26:55Z
parent: credfolio2-dwid
---

First screen experience: upload a resume, see it processed, view the result.

## User Experience

1. Land on clean upload page (no distractions, clear CTA)
2. Drag-drop or click to upload resume (PDF, DOCX)
3. See animated skeleton/shimmer while processing
4. Smooth transition to profile view when ready

## Technical Components

- Landing page with upload zone
- Resume-specific extraction schema (distinct from reference letters)
- Async processing via River job
- Skeleton UI with realistic content shapes
- Profile creation from extracted data

## Dependencies

- Requires credfolio2-he9y (Async resume parsing via GraphQL) completion
- Builds on existing File Upload Pipeline

## Acceptance Criteria

- [ ] Upload screen is the app's entry point
- [ ] Processing takes no more than 30s for typical resume
- [ ] Skeleton UI makes wait feel shorter
- [ ] Profile is immediately viewable after extraction