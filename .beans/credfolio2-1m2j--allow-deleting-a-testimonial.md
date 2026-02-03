---
# credfolio2-1m2j
title: Allow deleting a testimonial
status: completed
type: feature
priority: normal
created_at: 2026-02-03T10:33:12Z
updated_at: 2026-02-03T12:28:11Z
parent: credfolio2-2ex3
---

## Summary
Add the ability to delete testimonials from the profile page.

## Context
Users need to be able to remove testimonials that are inaccurate, outdated, or that they don't want displayed on their profile.

## Requirements
- Add delete option to existing kebab menu (kebab menu exists in TestimonialsSection.tsx:184-206)
- Implement GraphQL mutation for deleting testimonials
- Add confirmation dialog before deletion
- Handle UI state update after deletion (remove from list)

## Technical Decisions
- **Hard delete**: Testimonials will be permanently removed from the database (no soft delete/recovery)
- Testimonials are linked to reference letters - deletion removes only the testimonial, not the source document

## Checklist
- [x] Add `deleteTestimonial` GraphQL mutation to schema
- [x] Implement resolver for deletion mutation
- [x] Add "Delete" option to testimonial kebab menu in UI
- [x] Create confirmation dialog component
- [x] Wire up mutation call on confirm
- [x] Update local state/cache after successful deletion
- [x] Handle error states (show toast on failure)

## Definition of Done
- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] Visual verification with agent-browser (for UI changes)
- [x] All checklist items above are completed
- [x] Branch pushed and PR created for human review: https://github.com/BjRo/credfolio2/pull/65
