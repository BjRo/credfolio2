---
# credfolio2-ksna
title: Source attribution badges
status: draft
type: task
created_at: 2026-01-23T16:28:41Z
updated_at: 2026-01-23T16:28:41Z
parent: credfolio2-1kt0
---

Visual indicators showing which document contributed each piece of data.

## Badge Types

- "From Resume" - Data from initial resume
- "From [Author Name]'s Letter" - Data from reference letter
- "Validated by [Author]" - Existing data confirmed by letter

## Where to Show

- Skills section: small badge per skill
- Experience: badge on items with added highlights
- Testimonials: inherently show source

## Interaction

- Hover/click badge to see full source info
- Tooltip with document name, date uploaded
- Link to original file (if accessible)

## Design

- Subtle by default (icon or small text)
- Expandable on interaction
- Color coding by source type
- Accessible (not just color)

## Checklist

- [ ] Design badge variants
- [ ] Create SourceBadge component
- [ ] Add to skill tags
- [ ] Add to experience cards where applicable
- [ ] Implement tooltip with details