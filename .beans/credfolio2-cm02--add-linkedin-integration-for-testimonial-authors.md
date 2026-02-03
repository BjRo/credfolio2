---
# credfolio2-cm02
title: Edit testimonial author details (name, role, image, LinkedIn)
status: todo
type: feature
priority: normal
created_at: 2026-01-30T13:47:16Z
updated_at: 2026-02-03T12:53:37Z
parent: credfolio2-2ex3
---

Allow editing all author information including name, role, profile image, and LinkedIn URL. Handle "unknown" authors with special UI treatment.

## User Story

As a profile owner, I want to edit testimonial author details (name, title, company, image, LinkedIn) so I can:
1. Correct extraction errors when the LLM couldn't determine the author
2. Add a profile photo for visual recognition
3. Add LinkedIn verification links for trust

## Context

Reference letters sometimes don't explicitly name the author, causing the LLM extraction to return inconsistent values. This feature:
1. Standardizes unknown author handling (LLM returns `"unknown"`)
2. Provides visual differentiation for unknown authors
3. Allows post-upload editing of all author fields including a profile image

## Current State (Already Implemented)

- ✅ `Author` entity exists with `name`, `title`, `company`, `linkedInUrl` fields
- ✅ `updateAuthor` GraphQL mutation exists
- ✅ LinkedIn icon displays next to author name when URL exists
- ✅ Clicking LinkedIn icon opens profile in new tab
- ✅ Kebab menu exists on **quotes** (delete testimonial, view source)
- ✅ MinIO file storage infrastructure exists (used for PDFs)

## What's Needed

### Backend Changes

**1. LLM Prompt Update**
- Modify `reference_letter_extraction_system.txt` to instruct the LLM to return `"unknown"` as the author name when it cannot be determined
- Location: `src/backend/internal/infrastructure/llm/prompts/reference_letter_extraction_system.txt`

**2. Author Image Support**
- Add `image_id` column to `authors` table (FK to `files`)
- Database migration with the new column
- Update GraphQL schema: add `imageUrl` field to `Author` type
- Update `updateAuthor` mutation input to accept `imageId`
- Resolver logic to return the file URL when image exists

### Frontend Changes

**1. Kebab Menu on Author Header**
- Add kebab menu (MoreVertical icon) to the author header row in `TestimonialGroupCard`
- Menu items: "Edit author" (opens modal)
- Coexists with existing quote-level kebab menus

**2. Unknown Author Detection & UI**
- Detection: `author.name === "unknown"` or `author.name` is empty/null
- Visual treatment for unknown authors:
  - Dashed border on the card
  - Muted/subtle background color
  - Banner inside the card: "Author not detected — click to add details"
  - "Add details" CTA prominent in the header

**3. AuthorEditModal Component**
- Fields:
  - **Image** — File upload (stored in MinIO), displays current image or placeholder
  - **Name** — Text input (required)
  - **Title** — Text input (optional)
  - **Company** — Text input (optional)
  - **LinkedIn URL** — Text input with validation (optional)
  - **Relationship** — Dropdown: MANAGER, PEER, DIRECT_REPORT, CLIENT, OTHER
- Validation:
  - Name is required (cannot be empty)
  - LinkedIn URL must match `https://(www.)?linkedin.com/in/.*`
- On submit: call `updateAuthor` mutation

**4. Author Image Display**
- Replace the initials avatar circle with the uploaded image when present
- Keep initials avatar as fallback when no image

## Technical Notes

- File reference: [TestimonialsSection.tsx](src/frontend/src/components/profile/TestimonialsSection.tsx)
- Pattern reference: [EducationSection.tsx](src/frontend/src/components/profile/EducationSection.tsx) for kebab menu + edit modal
- LLM prompt: [reference_letter_extraction_system.txt](src/backend/internal/infrastructure/llm/prompts/reference_letter_extraction_system.txt)
- Author schema: [extraction.go:416-439](src/backend/internal/infrastructure/llm/extraction.go#L416-L439)

## Checklist

### Backend
- [ ] Update LLM prompt to return `"unknown"` for unidentifiable authors
- [ ] Create migration: add `image_id` column to `authors` table
- [ ] Update `Author` entity in Go with `ImageID` field
- [ ] Update GraphQL schema: add `imageUrl` to `Author` type
- [ ] Update `UpdateAuthorInput` to accept `imageId`
- [ ] Update resolver to handle image URL resolution

### Frontend
- [ ] Add kebab menu to author header in `TestimonialGroupCard`
- [ ] Create `AuthorEditModal` component with all fields
- [ ] Add image upload component (reuse existing file upload pattern)
- [ ] Wire up `updateAuthor` mutation in modal
- [ ] Add LinkedIn URL validation
- [ ] Add unknown author detection logic
- [ ] Style unknown author cards (dashed border, muted bg, banner)
- [ ] Replace initials avatar with uploaded image when present
- [ ] Update UI optimistically on successful edit

## Definition of Done

- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Visual verification with agent-browser (for UI changes)
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review
