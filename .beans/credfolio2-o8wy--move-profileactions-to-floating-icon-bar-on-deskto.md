---
# credfolio2-o8wy
title: Move ProfileActions to floating icon bar on desktop
status: in-progress
type: feature
priority: normal
created_at: 2026-02-06T22:44:25Z
updated_at: 2026-02-06T22:54:56Z
---

## Context

The profile page currently has an action menu ("Add Reference Letter", "Export PDF", "Upload Another Resume") rendered as a card at the very bottom of the page. This menu should be repositioned: on desktop, it becomes a vertical icon-only bar floating to the right of the ProfileHeader card. On smaller screens, it falls back to the current bottom-card layout.

**Design decisions:**
- Separate floating bar (not attached to header card), with a small gap
- Static positioning (not sticky) — scrolls away with the header
- Tooltips on hover (icons are label-less on desktop)
- Icons scale up slightly on hover
- Breakpoint: `lg` (1024px)

## Approach

Render the actions **in two places** with responsive visibility:
- **Desktop (lg+):** A new `ProfileActionsBar` component renders inside a `relative` wrapper around `ProfileHeader`, positioned with `absolute left-full` to float in the right gutter
- **Mobile (<lg):** The existing `ProfileActions` component stays at the page bottom, hidden on desktop via `lg:hidden`

To ensure the floating bar has room at the `lg` breakpoint (1024px), increase the outer wrapper's right padding to `lg:pr-20` (80px). This shifts the content column ~24px left of center, which is subtle and provides ~28px clearance for the bar.

CSS-only tooltips (no new dependencies) using the `group`/`group-hover` pattern already used throughout the codebase.

## Checklist

### Implementation
- [x] Refactor `ProfileActions.tsx`: extract shared action config, add `ProfileActionsBar` export with vertical icon buttons, CSS tooltips, hover scale effect
- [x] Update `index.ts`: add `ProfileActionsBar` barrel export
- [x] Update `page.tsx`: import `ProfileActionsBar`, wrap header in relative container, position bar with `absolute left-full ml-3`, add `lg:pr-20` to outer wrapper
- [x] Update `ProfileActions.test.tsx`: add `ProfileActionsBar` test suite (aria-labels, tooltips, click handlers, null when no handlers)

### Definition of Done
- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] Visual verification with agent-browser (for UI changes)
- [x] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review
- [ ] Automated code review passed (`@review-backend` and/or `@review-frontend`)