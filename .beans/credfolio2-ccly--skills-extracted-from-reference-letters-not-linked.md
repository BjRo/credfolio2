---
# credfolio2-ccly
title: Skills extracted from reference letters not linked to source document
status: todo
type: bug
created_at: 2026-02-02T14:18:26Z
updated_at: 2026-02-02T14:18:26Z
parent: credfolio2-2ex3
---

## Problem

When skills are extracted from reference letters, they are not properly:
1. Marked as credible/validated
2. Linked back to the reference letter that supports them

This makes it unclear which skills have third-party validation vs self-reported skills.

## Expected Behavior

When a skill is extracted from a reference letter:
- The skill should be marked as "validated" or "credible"
- The skill should link to the source reference letter that validates it
- The UI should visually distinguish validated skills from self-reported ones

## Current Behavior

Skills extracted from reference letters appear the same as self-reported skills, with no indication of their source or validation status.

## Impact

- Users cannot see which of their skills have third-party validation
- The trust signal from reference letters is lost
- Viewers of the profile cannot distinguish validated vs claimed skills

## Technical Investigation Needed

- [ ] Check how skills are currently extracted from reference letters
- [ ] Verify the data model supports skill → reference letter linking
- [ ] Determine if the link exists in the database but isn't surfaced in the UI

## Definition of Done
- [ ] Investigate root cause (data model vs extraction vs UI issue)
- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Visual verification with agent-browser (for UI changes)
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review