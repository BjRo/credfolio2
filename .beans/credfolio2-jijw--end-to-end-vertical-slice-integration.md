---
# credfolio2-jijw
title: 'Minimal E2E integration: Resume upload to profile view'
status: draft
type: task
priority: normal
created_at: 2026-01-20T15:31:10Z
updated_at: 2026-01-23T16:37:31Z
parent: credfolio2-tikg
blocking:
    - credfolio2-brav
    - credfolio2-63w0
---

Wire together all components for the first working end-to-end flow. Focus on **proving the integration works**, not polish.

## Goal

Upload a resume → extract via LLM → see basic profile data. Minimal UI, no fancy loading states.

## Integration Points

1. Frontend: Simple upload form (can reuse existing FileUpload component)
2. Backend: Upload mutation stores file in MinIO + creates DB record
3. Backend: River job triggers LLM extraction with resume schema
4. Backend: Extracted data saved to profile tables
5. Frontend: Basic page to view extracted profile data via GraphQL

## Acceptance Criteria

- [ ] User can upload a resume PDF through the UI
- [ ] File is stored in MinIO with metadata in PostgreSQL
- [ ] River job triggers and completes LLM extraction
- [ ] Extracted profile data (name, experience, skills) is saved to DB
- [ ] User can view extracted data on a basic profile page
- [ ] Errors are logged (fancy error UI not required)

## What This Is NOT

This task is intentionally minimal. The following are **out of scope** (handled by dwid tasks):

- Polished landing page design
- Skeleton/shimmer loading UI
- LinkedIn-style profile layout
- Reference letter enhancement
- Undo/history
- PDF export

## Dependencies

Requires these to be ready first:
- credfolio2-aevf (Resume extraction schema)
- credfolio2-j3ii (Profile database model)
- credfolio2-he9y (Async resume parsing via GraphQL)

## After This

Once this works, the dwid milestone tasks iterate on the foundation:
- Better upload UX (credfolio2-brav)
- Skeleton loading (credfolio2-63w0)
- Polished profile display (credfolio2-umxd tasks)
- Reference letter enhancement (credfolio2-1kt0 tasks)
- And more...