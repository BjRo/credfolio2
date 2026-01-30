---
# credfolio2-ksna
title: Credibility indicators and hover popovers
status: completed
type: feature
priority: normal
created_at: 2026-01-23T16:28:41Z
updated_at: 2026-01-30T13:25:02Z
parent: credfolio2-1kt0
---

Subtle credibility indicators on skills and experiences, with rich detail on hover.

## Credibility Dot Indicators

### Skills
Display dots next to each skill showing validation count:

- `●` (single dot) = 1 source (resume only)
- `●●` (double dot) = 2 sources (resume + 1 reference)
- `●●●` (triple dot) = 3+ sources (resume + 2+ references)

### Experiences
Display a subtle validation badge/icon when experience has validations:

- No indicator if unvalidated
- Small checkmark or shield icon if validated
- Number badge if multiple validations (e.g., "2")

## Hover Popover

When user hovers over a validated skill or experience, show a popover with:

```
┌─────────────────────────────────┐
│ Leadership                   ●●●│
│ ────────────────────────────────│
│                                 │
│ ✓ Resume                        │
│   Listed in skills section      │
│                                 │
│ ✓ John Smith                    │
│   Engineering Manager, Acme     │
│   "Jane's leadership during the │
│    migration was exceptional... │
│                                 │
│ ✓ Sarah Chen                    │
│   Tech Lead, Acme               │
│   "Led the team through a       │
│    challenging transition..."   │
│                                 │
│ [View full testimonials]        │
└─────────────────────────────────┘
```

## GraphQL Queries

```graphql
# Query validations for a specific skill
query skillValidations($skillId: ID!) {
  skillValidations(skillId: $skillId) {
    referenceLetter {
      id
      extractedData  # for author info
    }
    quoteSnippet
    createdAt
  }
}

# Query validations for a specific experience
query experienceValidations($experienceId: ID!) {
  experienceValidations(experienceId: $experienceId) {
    referenceLetter {
      id
      extractedData
    }
    quoteSnippet
    createdAt
  }
}
```

## Components

### CredibilityDots
- Props: `count: number` (1, 2, 3+)
- Renders appropriate number of dots
- Accessible: includes aria-label

### ValidationPopover
- Props: `validations: Validation[]`, `type: 'skill' | 'experience'`
- Renders list of sources with quotes
- Triggered on hover with delay
- Keyboard accessible (focus)

### Integration Points
- Add CredibilityDots to SkillTag component
- Add validation badge to ExperienceCard component
- Wrap both with ValidationPopover trigger

## Checklist

- [x] Create CredibilityDots component
- [x] Create ValidationPopover component
- [x] Add skillValidations query
- [x] Add experienceValidations query
- [x] Integrate dots into SkillTag
- [x] Integrate badge into ExperienceCard
- [x] Wire up hover/focus to show popover
- [x] Style popover with quotes and attribution
- [x] Add "View full testimonials" link
- [x] Ensure keyboard accessibility
- [x] Add loading state for popover content

## Definition of Done

- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] Visual verification with agent-browser (UI changes)
- [x] All checklist items above are completed
- [x] Branch pushed and PR created for human review
