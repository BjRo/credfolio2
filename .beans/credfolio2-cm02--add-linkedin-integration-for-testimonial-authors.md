---
# credfolio2-cm02
title: Add LinkedIn integration for testimonial authors
status: todo
type: feature
priority: normal
created_at: 2026-01-30T13:47:16Z
updated_at: 2026-02-03T10:53:14Z
parent: credfolio2-2ex3
---

Allow editing all author information including LinkedIn profile links.

## User Story
As a profile owner, I want to edit testimonial author details (name, title, company, LinkedIn) so I can correct extraction errors and add verification links.

## Current State (Already Implemented)
- ✅ `Author` entity exists with `linkedInUrl` field
- ✅ `updateAuthor` GraphQL mutation exists
- ✅ LinkedIn icon displays next to author name when URL exists
- ✅ Clicking LinkedIn icon opens profile in new tab

## What's Left (Frontend Only)

**Add edit capability to author info:**
- Edit button (pencil icon) on author info section
- Edit modal with fields:
  - Name (required)
  - Title (optional)
  - Company (optional)
  - LinkedIn URL (optional, validated)
  - Relationship dropdown

**Validation:**
- LinkedIn URL must be `https://linkedin.com/in/...` or `https://www.linkedin.com/in/...`
- Name is required, others optional

## Checklist
- [ ] Create AuthorEditModal component
- [ ] Add edit button (pencil icon) to author info display in TestimonialsSection
- [ ] Wire up `updateAuthor` mutation call in modal
- [ ] Add LinkedIn URL validation (regex for linkedin.com/in/...)
- [ ] Update UI optimistically on successful edit

## Definition of Done
- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Visual verification with agent-browser (for UI changes)
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review
