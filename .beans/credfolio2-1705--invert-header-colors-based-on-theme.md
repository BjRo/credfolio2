---
# credfolio2-1705
title: Invert header colors based on theme
status: completed
type: feature
priority: normal
created_at: 2026-01-29T11:50:02Z
updated_at: 2026-01-29T12:34:41Z
---

Change header styling so that:
- Light theme: dark header background
- Dark theme: light/white header background

This creates visual contrast between the header and body content.

Also update dark mode card colors so cards stand out from the background (like shadcn reference):
- Use visible borders on cards instead of shadows
- Slight background difference between page and cards
- More visible border color in dark mode

## Checklist
- [x] Investigate current header implementation
- [x] Update header styles for light theme (dark background)
- [x] Update header styles for dark theme (light background)
- [x] Ensure text/icons have appropriate contrast in both modes
- [x] Update card styling: borders instead of shadows
- [x] Adjust dark mode CSS variables (card, border, input)
- [x] Swap background/card colors in dark mode (cards darker than page)
- [x] Visual verification with agent-browser

## Definition of Done
- [x] Tests written (if applicable) - existing tests pass, no new tests needed for CSS-only change
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] Visual verification with agent-browser (for UI changes)
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review