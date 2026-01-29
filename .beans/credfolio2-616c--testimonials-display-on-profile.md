---
# credfolio2-616c
title: Testimonials display on profile
status: draft
type: feature
priority: normal
created_at: 2026-01-29T13:33:00Z
updated_at: 2026-01-29T13:33:00Z
parent: credfolio2-1kt0
---

Display testimonials section on profile page with full quotes and attribution.

## Section Design

New section on profile page: "What Others Say"

```
┌─────────────────────────────────────────────────────────┐
│  What Others Say                                        │
│                                                         │
│  ┌─────────────────────────────────────────────────┐   │
│  │  "Jane's technical leadership during our cloud  │   │
│  │   migration was exceptional. She not only led   │   │
│  │   the architecture decisions but mentored the   │   │
│  │   entire team through the transition."          │   │
│  │                                                  │   │
│  │   — John Smith                                   │   │
│  │     Engineering Manager at Acme Corp            │   │
│  │     Relationship: Direct manager (2020-2024)    │   │
│  │                                                  │   │
│  │   Validates: Leadership · Mentoring · Cloud     │   │
│  └─────────────────────────────────────────────────┘   │
│                                                         │
│  ┌─────────────────────────────────────────────────┐   │
│  │  "Working with Jane on the payments team was    │   │
│  │   a pleasure. Her deep knowledge of Go and      │   │
│  │   distributed systems elevated our whole team." │   │
│  │                                                  │   │
│  │   — Sarah Chen                                   │   │
│  │     Tech Lead at Acme Corp                      │   │
│  │     Relationship: Peer                          │   │
│  │                                                  │   │
│  │   Validates: Go · Distributed Systems           │   │
│  └─────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────┘
```

## GraphQL Query

```graphql
query profileTestimonials($profileId: ID!) {
  testimonials(profileId: $profileId) {
    id
    quote
    authorName
    authorTitle
    authorCompany
    relationship
    referenceLetter {
      id
    }
    createdAt
  }
}
```

## Components

### TestimonialsSection
- Fetches testimonials for profile
- Renders list of TestimonialCard components
- Shows empty state if no testimonials
- Header: "What Others Say"

### TestimonialCard
- Quote text (styled as blockquote)
- Author attribution line
- Relationship badge
- "Validates" skill links (clickable, scroll to skill)

### Skill Links
- Show which skills this testimonial validates
- Clicking scrolls to skill in Skills section
- Visual connection between testimonial and skills

## Empty State

When no testimonials:
```
┌─────────────────────────────────────────────────────────┐
│  What Others Say                                        │
│                                                         │
│  No testimonials yet.                                   │
│  Add a reference letter to include testimonials         │
│  from people who've worked with you.                    │
│                                                         │
│  [+ Add Reference Letter]                               │
└─────────────────────────────────────────────────────────┘
```

## Checklist

- [ ] Add testimonials query to GraphQL schema
- [ ] Create TestimonialsSection component
- [ ] Create TestimonialCard component
- [ ] Style quote blockquote
- [ ] Add author attribution styling
- [ ] Add relationship badge
- [ ] Create "Validates" skill links
- [ ] Implement scroll-to-skill on click
- [ ] Design and implement empty state
- [ ] Add to profile page layout
- [ ] Handle loading state

## Definition of Done

- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Visual verification with agent-browser (UI changes)
- [ ] All checklist items above are completed
- [ ] Branch pushed and PR created for human review
