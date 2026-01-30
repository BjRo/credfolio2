---
# credfolio2-616c
title: Testimonials display on profile
status: completed
type: feature
priority: normal
created_at: 2026-01-29T13:33:00Z
updated_at: 2026-01-30T11:20:31Z
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

- [x] Add testimonials query to GraphQL schema
- [x] Create TestimonialsSection component
- [x] Create TestimonialCard component
- [x] Style quote blockquote
- [x] Add author attribution styling
- [x] Add relationship badge
- [ ] Create "Validates" skill links (deferred - needs skill-testimonial linking in data model)
- [ ] Implement scroll-to-skill on click (deferred - depends on skill links)
- [x] Design and implement empty state
- [x] Add to profile page layout
- [x] Handle loading state

## Definition of Done

- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] Visual verification with agent-browser (UI changes)
- [x] All checklist items above are completed (except deferred items)
- [x] Branch pushed and PR created for human review (https://github.com/BjRo/credfolio2/pull/47)
