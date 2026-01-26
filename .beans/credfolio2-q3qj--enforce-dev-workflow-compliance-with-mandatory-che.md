---
# credfolio2-q3qj
title: Enforce dev-workflow compliance with mandatory checklists
status: completed
type: task
priority: normal
created_at: 2026-01-26T16:46:01Z
updated_at: 2026-01-26T16:50:38Z
---

Add mandatory Definition of Done checklist items and pre-completion reminder to ensure consistent test running, linting, and visual verification before marking work complete.

## Changes
1. Add mandatory checklist template for beans (Definition of Done)
2. Add prominent pre-completion reminder section to CLAUDE.md

## Checklist
- [x] Research how to add bean templates or default checklist items
- [x] Implement mandatory Definition of Done checklist for beans
- [x] Add pre-completion reminder section to CLAUDE.md
- [x] Push to branch and open PR

## Definition of Done
- [x] Tests written (TDD: write tests before implementation) — N/A: docs-only change
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures — N/A: docs-only change
- [x] Visual verification with agent-browser (for UI changes) — N/A: no UI changes
- [x] All other checklist items above are completed