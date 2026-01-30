---
# credfolio2-cm02
title: Add LinkedIn integration for testimonial authors
status: draft
type: feature
priority: normal
created_at: 2026-01-30T13:47:16Z
updated_at: 2026-01-30T13:50:54Z
parent: credfolio2-2ex3
---

Allow editing all author information including LinkedIn profile links.

## User Story
As a profile owner, I want to edit testimonial author details (name, title, company, LinkedIn) so I can correct extraction errors and add verification links.

## Prerequisites
**Depends on: credfolio2-m607 (Create Author entity)**
- Author must be a proper entity before we can edit it
- LinkedIn URL will be stored on Author, not Testimonial

## Implementation

### Backend (covered by Author entity task)
- `Author` entity includes `linkedin_url` field
- `updateAuthor` mutation for editing all fields

### Frontend Changes

**Author display with edit capability:**
- LinkedIn icon next to author name (when URL exists)
- Clicking LinkedIn icon opens profile in new tab
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

### Tasks
- [ ] Create AuthorEditModal component
- [ ] Add edit button to author info display
- [ ] Implement updateAuthor mutation call
- [ ] Add LinkedIn icon with link
- [ ] Add URL validation for LinkedIn
- [ ] Update optimistically on edit

## Definition of Done
- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Visual verification with agent-browser (for UI changes)
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review