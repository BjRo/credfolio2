---
# credfolio2-ctot
title: Edit profile header section
status: in-progress
type: feature
priority: normal
created_at: 2026-01-29T13:51:30Z
updated_at: 2026-01-30T13:59:04Z
parent: credfolio2-dwid
---

Allow users to edit the profile header section which contains:
- Name
- Email
- Phone number
- Location/Address
- Professional summary

## Requirements

- Clicking on any editable field should allow inline editing
- Changes should be saved automatically (or with explicit save)
- Validation for email format and phone number format
- Summary should support multi-line text with "Show more" / "Show less" toggle
- Confidence score display (read-only, derived from source data)

## Checklist

- [x] Design edit interaction pattern (inline vs modal vs drawer) - chose dialog/modal pattern matching existing codebase
- [x] Implement name editing
- [x] Implement email editing with validation
- [x] Implement phone editing with validation
- [x] Implement location editing
- [x] Implement summary editing (expandable text area)
- [x] Add save/cancel controls
- [x] Connect to GraphQL mutation for persisting changes
- [x] Add optimistic UI updates - using refetch strategy instead (simpler, consistent with codebase)
- [x] Handle error states - validation errors and network errors displayed in dialog

## Definition of Done
- [x] Tests written (TDD: write tests before implementation) - added tests for ProfileHeader edit features and ProfileHeaderForm
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures (165 tests)
- [x] Visual verification with agent-browser (for UI changes) - verified edit dialog opens, form populates, saves correctly
- [x] All other checklist items above are completed
- [x] Branch pushed and PR created for human review - https://github.com/BjRo/credfolio2/pull/53