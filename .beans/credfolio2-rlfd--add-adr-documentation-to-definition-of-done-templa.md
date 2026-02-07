---
# credfolio2-rlfd
title: Add ADR documentation to Definition of Done template
status: in-progress
type: task
priority: normal
created_at: 2026-02-07T16:33:12Z
updated_at: 2026-02-07T19:23:32Z
parent: credfolio2-ynmd
---

Add a conditional ADR item to the Definition of Done template so the decision skill gets used when appropriate.

## Why

The `/decision` skill exists but rarely gets invoked because nothing prompts Claude to consider whether an ADR is warranted. Adding it to the Definition of Done makes it part of the standard completion check — and with the TaskCompleted hook (credfolio2-tdjg) enforcing checklist completion, it becomes a real gate.

## What

Add a conditional item to the Definition of Done template (`.claude/templates/definition-of-done.md`):

```markdown
## Definition of Done
- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Visual verification via QA subagent (for UI changes)
- [ ] ADR written via `/decision` skill (if new dependencies, patterns, or architectural changes were introduced)
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review
- [ ] Automated code review via review-backend and/or review-frontend subagents
```

### When an ADR is warranted

Make it clear in the template (or in the decision skill description) what qualifies:
- Adding or removing a dependency (Go module, npm package)
- Introducing a new architectural pattern or abstraction
- Changing build/deployment configuration
- Deprecating an existing approach
- Choosing between multiple viable technical options

### When it's NOT needed
- Bug fixes using existing patterns
- Adding features following established conventions
- Routine refactoring without design changes
- Test additions or improvements
- Documentation updates (unless documenting a new documentation strategy)
- Configuration changes that follow existing patterns

## Notes

- The item is conditional ("if...") so Claude can check it off with "N/A — no architectural changes" when it doesn't apply
- Also review the `/decision` skill description to make these trigger criteria explicit
- CLAUDE.md includes the template via `@` reference on line 204, so it picks up changes automatically

## Implementation Plan

### Files to Modify

1. **`.claude/templates/definition-of-done.md`** — Insert the conditional ADR checklist item between "Visual verification" and "All other checklist items"
2. **`.claude/skills/decision/SKILL.md`** — Add a "When This Skill Does NOT Apply" section after the existing "When This Skill Applies" section with explicit exclusion criteria
3. **`.claude/hooks/tests/test-validate-bean-dod.sh`** — Update hardcoded `DOD_BODY` variable to include the new ADR item so the hook test keeps passing

### Files NOT Changed (and why)

- **`validate-bean-dod.sh` hook** — Reads the template dynamically with `grep -qi`; handles the new item automatically
- **`CLAUDE.md`** — Includes the template via `@` reference; updates itself
- **`decisions/README.md`** — Already has matching "When to Document" criteria

## Checklist
- [x] Update DoD template with conditional ADR item
- [x] Update `/decision` skill description with clear trigger and exclusion criteria
- [x] Update hook test fixture (`DOD_BODY`) to include new ADR item
- [x] Run hook tests (`bash .claude/hooks/tests/test-validate-bean-dod.sh`) — all pass
- [x] Run `pnpm lint` and `pnpm test` — all pass

## Definition of Done
- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] ADR written via `/decision` skill (if new dependencies, patterns, or architectural changes were introduced) — N/A, no architectural changes
- [x] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review
- [ ] Automated code review passed via `@review-backend` and/or `@review-frontend` subagents (via Task tool)
