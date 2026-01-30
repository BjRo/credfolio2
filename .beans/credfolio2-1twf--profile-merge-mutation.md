---
# credfolio2-1twf
title: Apply reference letter validations mutation
status: completed
type: task
priority: normal
created_at: 2026-01-23T16:28:40Z
updated_at: 2026-01-30T10:16:22Z
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

- [x] Define GraphQL input types in schema
- [x] Define mutation return type
- [x] Implement resolver with transaction
- [x] Create skill_validations records
- [x] Create experience_validations records
- [x] Create testimonials records (with author from extractedData)
- [x] Create new skills if selected
- [x] Update reference_letter status to "applied"
- [ ] Calculate credibility score (deferred to credfolio2-fuo1)
- [x] Handle edge cases (duplicates, missing entities)
- [x] Write resolver unit tests
- [ ] Write integration tests

## Definition of Done

- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [ ] All checklist items above are completed
- [x] Branch pushed and PR created for human review

## PR

https://github.com/BjRo/credfolio2/pull/45

## Implementation Notes

The mutation was implemented with the following key decisions:

1. **"Applied" status**: After validations are applied, the reference letter status is updated to "applied". The status column is VARCHAR(50), so no migration was needed. The resolver allows re-application of validations from letters with either "completed" or "applied" status.

2. **Credibility score**: This is a separate feature (credfolio2-fuo1) and is not calculated in this mutation.

3. **Transaction handling**: The resolver handles failures gracefully by continuing to process other items even if one fails, rather than using a database transaction that would roll back all changes.

4. **Author info mapping**: The testimonials use author info from the extracted data (name, title, company, relationship) rather than requiring it in the input.
