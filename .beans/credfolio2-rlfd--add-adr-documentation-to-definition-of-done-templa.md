---
# credfolio2-rlfd
title: Add ADR documentation to Definition of Done template
status: todo
type: task
created_at: 2026-02-07T16:33:12Z
updated_at: 2026-02-07T16:33:12Z
parent: credfolio2-ynmd
---

Add a conditional ADR item to the Definition of Done template so the decision skill gets used when appropriate.

## Why

The `/decision` skill exists but rarely gets invoked because nothing prompts Claude to consider whether an ADR is warranted. Adding it to the Definition of Done makes it part of the standard completion check — and with the TaskCompleted hook (credfolio2-tdjg) enforcing checklist completion, it becomes a real gate.

## What

Add a conditional item to the Definition of Done template (`.claude/templates/definition-of-done.md` once credfolio2-o624 creates it, or CLAUDE.md in the meantime):

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

## Notes

- The item is conditional ("if...") so Claude can check it off with "N/A — no architectural changes" when it doesn't apply
- Also review the `/decision` skill description to make these trigger criteria explicit
- Depends on credfolio2-o624 (DoD template file) if we want a single source of truth, but can also be added directly to CLAUDE.md now

## Definition of Done
- [ ] DoD template updated with conditional ADR item
- [ ] Decision skill description updated with clear trigger criteria
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Branch pushed and PR created for human review