---
# credfolio2-fuo1
title: Credibility score calculation and display
status: draft
type: feature
priority: normal
created_at: 2026-01-29T13:32:25Z
updated_at: 2026-01-29T13:32:25Z
parent: credfolio2-1kt0
---

Calculate and display profile credibility percentage with breakdown panel.

## Credibility Score Calculation

### Formula
```
credibilityScore = (validatedSkills + validatedExperiences) / (totalSkills + totalExperiences) * 100
```

Where:
- `validatedSkills` = skills with at least 1 validation from a reference letter
- `validatedExperiences` = experiences with at least 1 validation
- `totalSkills` = all skills in profile
- `totalExperiences` = all experiences in profile

### GraphQL Query

```graphql
query profileCredibility($profileId: ID!) {
  profileCredibility(profileId: $profileId) {
    overallScore: Float!           # 0.0 to 100.0
    skills {
      validated: Int!
      total: Int!
      percentage: Float!
    }
    experiences {
      validated: Int!
      total: Int!
      percentage: Float!
    }
    sources {
      referenceLetter {
        id
        extractedData  # for author name
      }
      validationCount: Int!
    }
  }
}
```

## UI Components

### CredibilityScoreBar (Always Visible)
Displays in profile header:

```
┌─────────────────────────────────────────────────┐
│ Profile Credibility              72% backed     │
│ ██████████████████░░░░░░░        by references  │
│                                                 │
│ 3 sources: Resume + 2 reference letters  [▼]   │
└─────────────────────────────────────────────────┘
```

- Progress bar with percentage
- Source count summary
- Expand/collapse toggle

### CredibilityBreakdown (Expandable)

```
┌─────────────────────────────────────────────────┐
│ Profile Credibility Breakdown              [×]  │
├─────────────────────────────────────────────────┤
│                                                 │
│  Skills         ████████████████░░░░  9/12 (75%)│
│  Experiences    ██████████████████░░  4/5  (80%)│
│                                                 │
│  Sources Contributing                           │
│  ┌─────────────────────────────────────────┐   │
│  │ Resume                      Base profile│   │
│  │ John Smith                 +8 validations│   │
│  │ Sarah Chen                 +5 validations│   │
│  └─────────────────────────────────────────┘   │
│                                                 │
│  Tip: 4 skills have no reference backing        │
│     [+ Add Reference Letter]                    │
└─────────────────────────────────────────────────┘
```

- Category breakdown bars
- Source list with contribution counts
- Encouragement message + CTA

## Checklist

- [ ] Implement credibility calculation logic in backend
- [ ] Add profileCredibility query to GraphQL schema
- [ ] Create CredibilityScoreBar component
- [ ] Create CredibilityBreakdown component
- [ ] Add expand/collapse interaction
- [ ] Style progress bars
- [ ] Add source list with icons
- [ ] Add "Add Reference Letter" CTA
- [ ] Integrate into profile page header
- [ ] Handle edge cases (no skills, no validations)

## Definition of Done

- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Visual verification with agent-browser (UI changes)
- [ ] All checklist items above are completed
- [ ] Branch pushed and PR created for human review
