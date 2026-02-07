---
# credfolio2-7rfu
title: Remove current page ring highlight from PDF viewer
status: in-progress
type: task
priority: normal
created_at: 2026-02-07T14:28:19Z
updated_at: 2026-02-07T14:37:03Z
parent: credfolio2-klgo
---

Remove the red ring highlight around the current page in the PDF viewer.

## Context

The PDF viewer shows a `ring-2 ring-primary/30` border around the "current" page to indicate which page is active. However, since the viewer is scrollable and the page number in the top-left is not updated on scroll, this highlight is misleading — it stays on page 1 even when the user has scrolled to page 2 or 3. The page navigation controls (prev/next, direct input) also feel unnecessary given the scroll-based UX.

## Checklist

- [x] Remove the `ring-2 ring-primary/30` conditional class from the page wrapper div in `PDFViewer.tsx` (line ~222)
- [x] Update or remove any tests that assert the ring styling
- [x] Verify visually that pages render cleanly without the ring

## Technical Notes

- File: `src/frontend/src/components/viewer/PDFViewer.tsx` line 222
- Current code: `pageNumber === currentPage && "ring-2 ring-primary/30"`
- Simply remove the conditional — the `shadow-md bg-white` classes provide enough visual separation between pages

## Definition of Done
- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] Visual verification with agent-browser (for UI changes)
- [x] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review
- [ ] Automated code review passed (`@review-backend` and/or `@review-frontend`)