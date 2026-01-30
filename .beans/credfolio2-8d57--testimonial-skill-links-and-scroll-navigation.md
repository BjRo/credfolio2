---
# credfolio2-8d57
title: Testimonial skill links and scroll navigation
status: in-progress
type: feature
priority: normal
created_at: 2026-01-30T11:26:41Z
updated_at: 2026-01-30T11:26:41Z
parent: credfolio2-1kt0
---

Add "Validates" skill links to testimonial cards that show which skills each testimonial validates, with click-to-scroll navigation.

## Context

The testimonials section displays quotes from reference letters. Each reference letter may validate specific skills (via skill_validations table). We want to show these skill connections directly on testimonial cards.

## Design

From the original spec (credfolio2-616c):

```
│  │   Validates: Leadership · Mentoring · Cloud     │   │
│  └─────────────────────────────────────────────────┘   │
```

## Implementation Approach

Since skill_validations are linked to reference_letters (not directly to testimonials), we can:
1. For each testimonial, look up its reference_letter_id
2. Query skill_validations for that reference_letter_id
3. Display the validated skill names as clickable links
4. On click, scroll to the skill in the Skills section

## GraphQL Changes

Extend Testimonial type to include validated skills:

```graphql
type Testimonial {
  # ... existing fields ...
  validatedSkills: [ProfileSkill!]!  # Skills validated by this testimonial's reference letter
}
```

## Frontend Changes

### TestimonialCard Component
- Add "Validates:" section below author attribution
- Display skill names as clickable pills/tags
- Each skill links to scroll target in Skills section

### Skills Section
- Add `id` attributes to skill elements for scroll targeting
- Smooth scroll behavior on testimonial skill click

## Checklist

- [x] Add validatedSkills field to Testimonial GraphQL type
- [x] Implement resolver to fetch skill validations for testimonial's reference letter
- [x] Update TestimonialCard to display validated skill links
- [x] Add scroll target IDs to skills in SkillsSection
- [x] Implement smooth scroll on skill click
- [x] Add tests for new functionality

## Definition of Done

- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] Visual verification with agent-browser (UI changes)
- [x] All checklist items above are completed
- [ ] Branch pushed and PR created for human review