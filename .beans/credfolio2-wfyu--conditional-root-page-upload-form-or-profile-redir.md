---
# credfolio2-wfyu
title: 'Conditional root page: upload form or profile redirect'
status: todo
type: feature
priority: normal
created_at: 2026-01-27T13:07:27Z
updated_at: 2026-01-27T13:07:59Z
parent: credfolio2-v5dw
---

The root page (/) should conditionally render based on whether the current user already has an extracted resume/profile:

- **No profile exists**: Show the resume upload form directly on the root page, so new users land on the upload experience immediately.
- **Profile exists**: Redirect to the user's profile page (/profile/{resumeId}).

This eliminates the need for users to navigate to a separate upload page and streamlines the entry point into the app.

## Checklist
- [ ] Add a backend query (or use existing) to check if the current user has a resume/profile
- [ ] Update the root page to fetch profile existence on load
- [ ] Render the upload form inline when no profile exists
- [ ] Redirect to /profile/{resumeId} when a profile exists
- [ ] Handle loading state while checking for profile

## Definition of Done
- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Visual verification with agent-browser (for UI changes)
- [ ] All other checklist items above are completed