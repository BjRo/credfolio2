---
# credfolio2-mt4c
title: Hover-only kebab menus
status: completed
type: task
priority: normal
created_at: 2026-01-28T14:17:12Z
updated_at: 2026-01-28T14:25:21Z
parent: credfolio2-v5dw
---

Show the three-dot kebab menu only on hover, not permanently visible.

## Scope

- **Education cards**: Show menu when hovering over the entire card
- **Work experience cards**: Show menu when hovering over the entire card
- **Skill pills**: Show menu when hovering over the individual pill only

## Implementation

- [x] Add hover state to education card component
- [x] Add hover state to work experience card component
- [x] Add hover state to skill pill component
- [x] Ensure keyboard accessibility (menu still reachable via focus)

## Definition of Done

- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] Visual verification with agent-browser (for UI changes)
- [x] All other checklist items above are completed