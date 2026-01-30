---
# credfolio2-lvdo
title: Collapse testimonials by author
status: draft
type: feature
priority: normal
created_at: 2026-01-30T13:47:26Z
updated_at: 2026-01-30T13:51:05Z
parent: credfolio2-2ex3
---

Group testimonials from the same author to reduce visual repetition and improve scannability.

## User Story
As a profile viewer, I want testimonials from the same person grouped together so I can quickly scan different perspectives without seeing repeated author information.

## Prerequisites
**Depends on: credfolio2-m607 (Create Author entity)**
- With proper Author entity, grouping is by `author_id` (no string matching needed)
- Author deduplication handled at extraction time

## Implementation

### UI Design
**Grouped layout:**
- Show author info once (avatar, name, title, company, relationship, LinkedIn)
- List all quotes from that author as bullet points or sub-cards
- Collapsible: show first quote expanded, rest collapsed
- "Show X more quotes" toggle

### Example
Before (current):
```
┌─────────────────────────────┐
│ "Quote 1..."                │
│ Amit Matani, CEO at Company │
└─────────────────────────────┘
┌─────────────────────────────┐
│ "Quote 2..."                │
│ Amit Matani, CEO at Company │
└─────────────────────────────┘
```

After (grouped):
```
┌─────────────────────────────┐
│ Amit Matani          [LinkedIn] [Edit] │
│ CEO at Company · Manager    │
│                             │
│ • "Quote 1..."        [📄]  │
│ • "Quote 2..."        [📄]  │
│   ▼ Show 2 more quotes      │
└─────────────────────────────┘
```

### Grouping Logic
With Author entity:
```typescript
// Group testimonials by author_id
const groupedByAuthor = testimonials.reduce((acc, t) => {
  const authorId = t.author.id;
  if (!acc[authorId]) {
    acc[authorId] = { author: t.author, testimonials: [] };
  }
  acc[authorId].testimonials.push(t);
  return acc;
}, {});
```

### Tasks
- [ ] Create TestimonialGroup component (shows author + multiple quotes)
- [ ] Implement grouping logic in TestimonialsSection
- [ ] Add expand/collapse functionality
- [ ] Default state: expanded for 1-2 quotes, collapsed for 3+
- [ ] Ensure source badges display on each quote
- [ ] Ensure validated skills still link correctly

## Definition of Done
- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Visual verification with agent-browser (for UI changes)
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review