---
# credfolio2-ctot
title: Edit profile header section
status: draft
type: feature
created_at: 2026-01-29T13:51:30Z
updated_at: 2026-01-29T13:51:30Z
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

- [ ] Design edit interaction pattern (inline vs modal vs drawer)
- [ ] Implement name editing
- [ ] Implement email editing with validation
- [ ] Implement phone editing with validation  
- [ ] Implement location editing
- [ ] Implement summary editing (expandable text area)
- [ ] Add save/cancel controls
- [ ] Connect to GraphQL mutation for persisting changes
- [ ] Add optimistic UI updates
- [ ] Handle error states

## Definition of Done
- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Visual verification with agent-browser (for UI changes)
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review