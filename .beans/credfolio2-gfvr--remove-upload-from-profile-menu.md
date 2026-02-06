---
# credfolio2-gfvr
title: Remove Upload link from site header
status: completed
type: task
priority: normal
created_at: 2026-02-06T12:02:44Z
updated_at: 2026-02-06T12:33:02Z
parent: credfolio2-dwid
---

Remove the "Upload" link from the site header navigation (`src/frontend/src/components/site-header.tsx`). The link currently appears in the top-right of the header on every page and points to `/upload`. It shouldn't be shown — first-time users are already redirected to `/upload` from the homepage, and the profile page has its own upload actions.

**Out of scope:** The `ProfileActions` component and its "Upload Another Resume" button stay as-is.

## Checklist
- [x] Remove the "Upload" link from the site header
- [x] Clean up any unused imports/code in site-header.tsx

## Definition of Done
- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] Visual verification with agent-browser (for UI changes)
- [x] All other checklist items above are completed
- [x] Branch pushed and PR created for human review
- [x] Automated code review passed (`@review-backend` and/or `@review-frontend`)