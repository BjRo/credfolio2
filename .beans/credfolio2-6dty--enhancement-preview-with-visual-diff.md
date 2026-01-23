---
# credfolio2-6dty
title: Enhancement preview with visual diff
status: draft
type: feature
priority: normal
created_at: 2026-01-23T16:28:37Z
updated_at: 2026-01-23T16:29:49Z
parent: credfolio2-1kt0
blocking:
    - credfolio2-1twf
---

Show what a reference letter would add/change to the profile.

## Preview Screen

After reference letter is processed, show:

### New Items (green highlight/badge)
- New skills mentioned
- New accomplishments
- Testimonial quotes

### Enhanced Items (blue highlight)
- Existing skills now validated by reference
- Positions with added context

### Testimonial Section
- Author name, title, relationship
- Key quotes extracted
- Recommendation strength indicator

## Visual Treatment

- Green "NEW" badges on new items
- Blue "ENHANCED" badges on validated items
- Side-by-side or inline diff view
- Clear before/after comparison

## Actions

- "Accept All" button
- "Reject All" button
- Individual toggle per item (future enhancement)

## Checklist

- [ ] Design preview layout
- [ ] Create EnhancementPreview component
- [ ] Show new skills with badges
- [ ] Show new testimonials
- [ ] Show enhanced/validated items
- [ ] Add Accept/Reject buttons
- [ ] Connect to merge mutation