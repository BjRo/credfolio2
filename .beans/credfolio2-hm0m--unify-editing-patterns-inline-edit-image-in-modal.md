---
# credfolio2-hm0m
title: 'Unify editing patterns: inline edit + image in modal for both profile header and testimonials'
status: in-progress
type: feature
priority: normal
created_at: 2026-02-03T16:36:16Z
updated_at: 2026-02-03T16:37:01Z
---

## Problem

The profile header and testimonials editing have inconsistent patterns:

**Profile Header:**
- ✅ Has inline editing (avatar can be clicked directly to upload)
- ❌ Image is NOT part of the edit modal

**Testimonials:**
- ✅ Image is part of the edit modal (AuthorEditModal)
- ❌ NO inline editing capability

## Goal

Unify both components to support BOTH patterns:
1. **Inline editing** - Click directly on avatar/image to upload (like profile header)
2. **Image in modal** - Image editing also available inside the edit modal (like testimonials)

## Files to Modify

**Profile Header:**
- `src/frontend/src/components/profile/ProfileHeaderForm.tsx` - Add image upload section
- `src/frontend/src/components/profile/ProfileHeaderFormDialog.tsx` - Pass image props

**Testimonials:**
- `src/frontend/src/components/profile/TestimonialsSection.tsx` - Add inline edit capability to author avatar
- `src/frontend/src/components/profile/AuthorEditModal.tsx` - Keep existing image in modal

## Checklist

- [x] Add image upload section to ProfileHeaderForm (matching AuthorEditModal pattern)
- [x] Pass required props for image handling through ProfileHeaderFormDialog
- [x] Add inline avatar click-to-upload on testimonial author avatars (matching ProfileAvatar pattern)
- [ ] Ensure both inline and modal image edits work correctly together
- [x] Test both editing flows work as expected

## Definition of Done
- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [ ] Visual verification with agent-browser (for UI changes)
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review