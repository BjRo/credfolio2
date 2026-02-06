---
# credfolio2-600v
title: Add item-level select/deselect in extraction review step
status: todo
type: task
priority: normal
created_at: 2026-02-05T23:09:37Z
updated_at: 2026-02-06T07:28:39Z
parent: credfolio2-3ram
---

## Summary

The ExtractionReview step currently shows extracted data (experiences, education, skills, testimonials) in a read-only list with a single "Import to profile" button that imports everything. Users cannot select or deselect individual items before importing.

Users should be able to review extracted data and choose which items to import — e.g., uncheck an incorrectly extracted work experience or deselect skills they don't want on their profile.

## Current Behavior

- `ExtractionReview.tsx` renders `CareerInfoSection` and `TestimonialSection` as read-only displays
- "Import to profile" sends `resumeId` and `referenceLetterID` to `importDocumentResults` mutation
- Backend materializes ALL extracted data — no filtering

## Desired Behavior

- Each extracted item (experience, education entry, skill, testimonial) has a checkbox
- All items start selected by default
- Users can deselect items they don't want imported
- Only selected items get materialized into profile tables

## Implementation Considerations

- Frontend: Add selection state per item category, render checkboxes, pass selection to import mutation
- Backend: `importDocumentResults` mutation or `MaterializeResumeData` needs to accept item indices or IDs to filter which items to materialize
- Alternative: Filter on frontend and send only selected items as a new input shape

## Checklist

- [ ] Design the selection UX (checkboxes per item vs. per category)
- [ ] Add selection state management to ExtractionReview
- [ ] Update import mutation input to support item-level filtering
- [ ] Update MaterializationService to respect selection
- [ ] Update tests

### Definition of Done
- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Visual verification with agent-browser (for UI changes)
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review