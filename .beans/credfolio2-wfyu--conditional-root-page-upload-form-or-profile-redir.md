---
# credfolio2-wfyu
title: 'Conditional root page: upload form or profile redirect'
status: completed
type: feature
priority: normal
created_at: 2026-01-27T13:07:27Z
updated_at: 2026-01-27T16:41:25Z
parent: credfolio2-v5dw
---

The root page (/) should conditionally render based on whether the current user already has an extracted resume/profile:

- **No profile exists**: Show the resume upload form directly on the root page, so new users land on the upload experience immediately.
- **Profile exists**: Redirect to the user's profile page (/profile/{resumeId}).

This eliminates the need for users to navigate to a separate upload page and streamlines the entry point into the app.

## Checklist
- [x] Add a backend query (or use existing) to check if the current user has a resume/profile
- [x] Update the root page to fetch profile existence on load
- [x] Render the upload form inline when no profile exists
- [x] Redirect to /profile/{resumeId} when a profile exists
- [x] Handle loading state while checking for profile

## Definition of Done
- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] Visual verification with agent-browser (for UI changes)
- [x] All other checklist items above are completed