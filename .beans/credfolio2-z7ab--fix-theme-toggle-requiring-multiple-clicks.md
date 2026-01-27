---
# credfolio2-z7ab
title: Fix theme toggle requiring multiple clicks
status: in-progress
type: bug
priority: normal
created_at: 2026-01-27T17:40:04Z
updated_at: 2026-01-27T17:41:05Z
---

The theme toggle skip logic only handles one direction: skipping 'light' when system resolves to light. But it doesn't handle cycling from 'dark' → 'system' when the OS also prefers dark — that produces no visible change, requiring another click.

**Root cause:** The skip check compares the next theme name against `resolvedTheme`, but when next is `"system"`, it won't literally match `"dark"` even though system resolves to dark.

**Fix:** Use `systemTheme` from next-themes to compute what the next theme would actually resolve to visually, then skip if that matches the current `resolvedTheme`.

## Checklist
- [x] Update `cycleTheme` to use `systemTheme` for proper resolved comparison
- [x] Add test for cycling from dark → system when OS prefers dark
- [x] Verify existing tests still pass

## Definition of Done
- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [ ] Visual verification with agent-browser (for UI changes)
- [x] All other checklist items above are completed