---
# credfolio2-q3qj
title: Enforce dev-workflow compliance with mandatory checklists
status: in-progress
type: task
created_at: 2026-01-26T16:46:01Z
updated_at: 2026-01-26T16:46:01Z
---

Add mandatory Definition of Done checklist items and pre-completion reminder to ensure consistent test running, linting, and visual verification before marking work complete.

## Changes
1. Add mandatory checklist template for beans (Definition of Done)
2. Add prominent pre-completion reminder section to CLAUDE.md

## Checklist
- [x] Research how to add bean templates or default checklist items
- [x] Implement mandatory Definition of Done checklist for beans
- [x] Add pre-completion reminder section to CLAUDE.md
- [ ] Push to branch and open PR

## Definition of Done
- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Visual verification with agent-browser (for UI changes)
- [ ] All other checklist items above are completed