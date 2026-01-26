---
# credfolio2-4tht
title: Skills display with categories
status: completed
type: task
priority: normal
created_at: 2026-01-23T16:28:01Z
updated_at: 2026-01-26T10:50:20Z
parent: credfolio2-umxd
---

Show skills grouped by category with tag/badge styling.

## Design

- Section header: "Skills"
- Subsections: Technical, Soft Skills, Languages
- Skills as tags/badges
- Possibly with hover for source info

## SkillTag Component

- Skill name
- Optional: proficiency indicator (for languages)
- Optional: source badge (tiny indicator of where it came from)

## SkillsSection Component

- Groups skills by category
- Handles empty categories gracefully
- Responsive: wrap nicely on mobile

## Checklist

- [x] Create SkillTag component
- [x] Create SkillsSection component
- [ ] ~~Group by category~~ (deferred: resume data model only has string[] for skills, no categories)
- [x] Style tags appropriately (with hover effects, proper spacing)
- [x] Add accessible list structure
- [ ] ~~Add language proficiency display~~ (deferred: data model doesn't support proficiency levels)