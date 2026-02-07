---
# credfolio2-cm82
title: 'Polish PDF viewer: site header, dark background, theme handling'
status: completed
type: task
priority: normal
created_at: 2026-02-07T14:32:39Z
updated_at: 2026-02-07T15:32:12Z
parent: credfolio2-klgo
---

## Polish PDF viewer: site header, dark background, theme handling

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

## Implementation Plan

### Architecture Overview

The fix involves three files with purely CSS/layout changes, plus test updates. No new components, hooks, or backend changes are needed. The SiteHeader is already rendered in layout.tsx for all routes -- it just needs to not be pushed off-screen by the viewer.

### Step-by-step Implementation

#### Step 1: Add `--header-height` CSS custom property

**File:** `src/frontend/src/app/globals.css`

Add `--header-height: 3.5rem;` to the `:root` block. This matches the `h-14` on the SiteHeader inner div. Using a CSS custom property keeps the value DRY and makes `calc()` expressions maintainable. No need for a dark-mode variant -- the height is theme-independent.

This property will be consumed in the viewer page via `h-[calc(100dvh-var(--header-height))]`.

#### Step 2: Fix the viewer page layout height

**File:** `src/frontend/src/app/viewer/page.tsx`

**Main container (line 172):** Change `h-screen` to `h-[calc(100dvh-var(--header-height))]`. Using `100dvh` (dynamic viewport height) instead of `100vh` handles mobile browser chrome better. This makes the viewer fill exactly the remaining viewport below the SiteHeader.

**Error/loading states (lines 48, 60):** The `ErrorPage` and `LoadingSkeleton` components use `min-h-screen` for centering. These are rendered as the entire page content, so the SiteHeader is already visible above them. However, for consistency they could use `min-h-[calc(100dvh-var(--header-height))]` so the centering accounts for the header. This is optional polish -- since `min-h-screen` still works (they scroll naturally), it is acceptable to leave them as-is and address only the main viewer container.

#### Step 3: Change PDF content area background to fixed dark gray

**File:** `src/frontend/src/components/viewer/PDFViewer.tsx`

**PDF content area (line 202):** Change `bg-muted/30` to `bg-neutral-700`. This gives a consistent dark gray background (#404040) in both light and dark mode, providing high contrast against white PDF pages. The `bg-neutral-700` class is built into Tailwind CSS 4 and does NOT follow the theme -- it is a fixed color.

The toolbars already use `bg-background` (page toolbar in `viewer/page.tsx` line 174, component toolbar in `PDFViewer.tsx` line 134). These are theme-aware and need no changes.

The PDF page wrappers (line 219) use `bg-white` with an intentional comment. Do NOT change these.

#### Step 4: Update tests

**File:** `src/frontend/src/app/viewer/page.test.tsx`

Existing tests verify the viewer renders correctly (toolbar, PDFViewer component, back navigation, info banner). The layout change from `h-screen` to `h-[calc(...)]` should not break any of these tests since they test behavior, not CSS class names.

No test changes are strictly required unless we want to add an explicit assertion that the viewer container does NOT use `h-screen`. Consider adding one test:

- **"viewer container accounts for header height"**: Render the viewer page in success state, find the main flex container, and assert it does NOT have the class `h-screen`.

**File:** `src/frontend/src/components/viewer/PDFViewer.test.tsx`

Similarly, the PDFViewer component tests mock react-pdf and test behavior. The background color change from `bg-muted/30` to `bg-neutral-700` will not affect any existing tests. No changes needed.

#### Step 5: Visual verification

Use `agent-browser` to verify in both light and dark mode:
1. SiteHeader is visible at the top of the viewer page
2. Theme toggle in the header works
3. PDF content area has consistent dark gray background in both modes
4. PDF pages are white and readable in both modes
5. Toolbar text/icons are legible in both modes
6. The viewer fills exactly the viewport below the header (no scrollbar on the body)
7. The info banner (highlight not found) still renders correctly between the page toolbar and PDF content

### Files to Modify

| File | Change |
|------|--------|
| `src/frontend/src/app/globals.css` | Add `--header-height: 3.5rem` to `:root` |
| `src/frontend/src/app/viewer/page.tsx` | Change `h-screen` → `h-[calc(100dvh-var(--header-height))]` |
| `src/frontend/src/components/viewer/PDFViewer.tsx` | Change `bg-muted/30` → `bg-neutral-700` |
| `src/frontend/src/app/viewer/page.test.tsx` | Optional: add test for no `h-screen` class |

### Risk Assessment

- **Low risk**: All changes are CSS-only. No logic, state, or data flow changes.
- **Header height coupling**: The `--header-height` variable must stay in sync with the `h-14` on SiteHeader. This is acceptable because the header height changes rarely, and the CSS variable makes it explicit.
- **`100dvh` browser support**: `dvh` is supported in all modern browsers (Chrome 108+, Firefox 108+, Safari 15.4+). This project targets modern browsers.
- **Tailwind `bg-neutral-700`**: This is a core Tailwind color. In Tailwind CSS 4, all color utilities are available by default without explicit configuration.

## Checklist

- [x] Add `--header-height: 3.5rem` CSS custom property to `:root` in globals.css
- [x] Change viewer page container from `h-screen` to `h-[calc(100dvh-var(--header-height))]`
- [x] Verify the SiteHeader is visible and functional on the viewer page (logo, theme toggle)
- [x] Change PDF content area background from `bg-muted/30` to `bg-neutral-700`
- [x] Ensure the viewer toolbars (page nav, zoom, document title) still follow the active theme
- [x] Test in both light and dark mode to verify:
  - SiteHeader visible in both modes
  - PDF content area is consistently dark gray in both modes
  - PDF pages remain white and readable in both modes
  - Toolbar text/icons are legible in both modes

## Definition of Done
- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] Visual verification with agent-browser (for UI changes)
- [x] All other checklist items above are completed
- [x] Branch pushed and PR created for human review
- [x] Automated code review passed (`@review-backend` and/or `@review-frontend`)