---
# credfolio2-h2op
title: Work experience timeline/cards
status: todo
type: task
priority: normal
created_at: 2026-01-23T16:28:00Z
updated_at: 2026-01-26T10:19:17Z
parent: credfolio2-umxd
---

Display work history as a visual timeline or card list.

## Design Options

1. **Timeline** - Vertical line with positions branching off
2. **Cards** - Stacked cards, most recent first

Recommendation: Cards for MVP, timeline as enhancement

## ExperienceCard Component

- Company logo placeholder
- Company name, job title
- Date range (formatted nicely)
- Location
- Description (collapsible if long)
- Highlights as bullet list
- Source indicator (which document contributed this)

## Behavior

- Most recent first
- Current job highlighted
- Expand/collapse for details

## Checklist

- [ ] Create ExperienceCard component
- [ ] Create ExperienceList container
- [ ] Add expand/collapse behavior
- [ ] Style with proper spacing
- [ ] Add source indicator badge
- [ ] Handle "current" job styling