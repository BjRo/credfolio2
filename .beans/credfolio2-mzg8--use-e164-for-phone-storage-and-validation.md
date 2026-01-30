---
# credfolio2-mzg8
title: Use E164 for phone storage and validation
status: draft
type: task
created_at: 2026-01-30T15:55:22Z
updated_at: 2026-01-30T15:55:22Z
parent: credfolio2-abtx
---

Store phone numbers in E164 format (+{country code}{number}) for consistent international phone number handling.

## Background

E164 is the international standard for phone number formatting:
- Format: `+{country code}{subscriber number}` (e.g., `+14155551234`)
- Max 15 digits (excluding the `+`)
- Ensures unambiguous international dialing
- Required by many telephony APIs (Twilio, etc.)

## Checklist

- [ ] Add phone number validation using E164 format in backend
- [ ] Update database schema to store phone in E164 format
- [ ] Add frontend phone input with country code selector
- [ ] Validate and normalize phone numbers on input
- [ ] Add migration for existing phone data (if any)
- [ ] Write tests for phone validation logic

## Definition of Done
- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Visual verification with agent-browser (for UI changes)
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review