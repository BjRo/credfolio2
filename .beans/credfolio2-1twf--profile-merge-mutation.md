---
# credfolio2-1twf
title: Apply reference letter validations mutation
status: draft
type: task
priority: normal
created_at: 2026-01-23T16:28:40Z
updated_at: 2026-01-29T00:00:00Z
parent: credfolio2-1kt0
blocking:
    - credfolio2-ksna
---

GraphQL mutation to apply selected validations from a reference letter to the profile.

## GraphQL Schema

### Input Types

```graphql
input ApplyValidationsInput {
  referenceLetterID: ID!
  skillValidations: [SkillValidationInput!]!
  experienceValidations: [ExperienceValidationInput!]!
  testimonials: [TestimonialInput!]!
  newSkills: [NewSkillInput!]!
}

input SkillValidationInput {
  profileSkillID: ID!
  quoteSnippet: String!
}

input ExperienceValidationInput {
  profileExperienceID: ID!
  quoteSnippet: String!
}

input TestimonialInput {
  quote: String!
  skillsMentioned: [String!]!
}

input NewSkillInput {
  name: String!
  category: SkillCategory!
  quoteContext: String
}
```

### Mutation

```graphql
mutation applyReferenceLetterValidations($input: ApplyValidationsInput!) {
  applyReferenceLetterValidations(input: $input) {
    referenceLetter {
      id
      status
    }
    profile {
      id
      credibilityScore
    }
    appliedCount {
      skillValidations: Int!
      experienceValidations: Int!
      testimonials: Int!
      newSkills: Int!
    }
  }
}
```

## Resolver Behavior

1. Validate reference letter exists and belongs to user
2. Validate all referenced skills/experiences exist
3. Create skill_validations records (upsert to handle duplicates)
4. Create experience_validations records
5. Create testimonials records with author info from extraction
6. Create new skills if selected (with source=reference_letter)
7. Update reference_letter status to "applied"
8. Calculate and return updated credibility score

## Edge Cases

- Duplicate validation (same skill + same letter): upsert/ignore
- Skill deleted between preview and apply: skip gracefully
- Reference letter already applied: return current state

## Checklist

- [ ] Define GraphQL input types in schema
- [ ] Define mutation return type
- [ ] Implement resolver with transaction
- [ ] Create skill_validations records
- [ ] Create experience_validations records
- [ ] Create testimonials records (with author from extractedData)
- [ ] Create new skills if selected
- [ ] Update reference_letter status
- [ ] Calculate credibility score
- [ ] Handle edge cases (duplicates, missing entities)
- [ ] Write resolver unit tests
- [ ] Write integration tests

## Definition of Done

- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] All checklist items above are completed
- [ ] Branch pushed and PR created for human review
