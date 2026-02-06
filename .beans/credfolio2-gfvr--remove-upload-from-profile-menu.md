---
# credfolio2-gfvr
title: Remove Upload from profile menu
status: todo
type: task
priority: normal
created_at: 2026-02-06T12:02:44Z
updated_at: 2026-02-06T12:02:44Z
parent: credfolio2-dwid
---

Remove the Upload-related navigation/actions from the profile page. Currently there are two places where Upload appears:

1. **ProfileActions component** (`src/frontend/src/components/profile/ProfileActions.tsx`): Has an "Upload Another Resume" button (`onUploadAnother` prop)
2. **Site header** (`src/frontend/src/components/site-header.tsx`): Has an "Upload" link in the top navigation

The profile page should no longer offer a direct upload action. Users should access upload functionality through other entry points (e.g. the main page or a dedicated route).

## Checklist
- [ ] Remove the `onUploadAnother` button from `ProfileActions` component
- [ ] Remove the `onUploadAnother` prop and handler from the profile page
- [ ] Clean up any unused imports/code related to the removed upload action
- [ ] Decide whether to also remove the "Upload" link from the site header (clarify with user)

## Definition of Done
- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Visual verification with agent-browser (for UI changes)
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review
- [ ] Automated code review passed (`@review-backend` and/or `@review-frontend`)