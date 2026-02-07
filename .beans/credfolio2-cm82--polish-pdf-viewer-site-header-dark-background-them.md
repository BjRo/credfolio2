---
# credfolio2-cm82
title: 'Polish PDF viewer: site header, dark background, theme handling'
status: todo
type: task
created_at: 2026-02-07T14:32:39Z
updated_at: 2026-02-07T14:32:39Z
parent: credfolio2-klgo
---

Improve the PDF viewer's visual design: show the Credfolio site header, use a dark gray background for the content area, and ensure consistent theme behavior.

## Context

The PDF viewer currently:
- Uses `h-screen` which pushes the SiteHeader (rendered in layout.tsx) off-screen
- Has a very faint `bg-muted/30` background that doesn't provide enough contrast with white PDF pages
- PDF pages are intentionally `bg-white` regardless of theme (correct behavior)

## Design Decisions

1. **Site header**: Keep the SiteHeader visible on the viewer page by changing the viewer layout from `h-screen` to account for the header height. The Credfolio branding and theme toggle remain accessible.
2. **Background**: Use a fixed dark gray background (e.g., `bg-neutral-700`) for the PDF content area regardless of light/dark theme — similar to how Google Drive's PDF viewer works. This provides consistent high contrast against white PDF pages.
3. **Theme**: The viewer toolbars (`bg-background`) continue to follow the theme. Only the PDF content area gets the fixed dark gray.

## Checklist

- [ ] Change the viewer page layout from `h-screen` to accommodate the SiteHeader height (e.g., `h-[calc(100vh-var(--header-height))]` or a flex layout approach)
- [ ] Verify the SiteHeader is visible and functional on the viewer page (logo, theme toggle)
- [ ] Change PDF content area background from `bg-muted/30` to a fixed dark gray (e.g., `bg-neutral-700` or similar)
- [ ] Ensure the viewer toolbars (page nav, zoom, document title) still follow the active theme
- [ ] Test in both light and dark mode to verify:
  - SiteHeader visible in both modes
  - PDF content area is consistently dark gray in both modes
  - PDF pages remain white and readable in both modes
  - Toolbar text/icons are legible in both modes

## Technical Notes

- `SiteHeader` is rendered in `layout.tsx` for all routes — it's already there, just hidden by the viewer's `h-screen`
- The viewer has two toolbars: page-level (back button + title in `viewer/page.tsx`) and component-level (page nav + zoom in `PDFViewer.tsx`)
- PDF page wrappers use `bg-white` with a comment explaining why — don't change this
- Consider using CSS custom properties for the header height to keep the calc maintainable

## Definition of Done
- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Visual verification with agent-browser (for UI changes)
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review
- [ ] Automated code review passed (`@review-backend` and/or `@review-frontend`)