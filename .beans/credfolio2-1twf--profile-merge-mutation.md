---
# credfolio2-1twf
title: Profile merge mutation
status: draft
type: task
priority: normal
created_at: 2026-01-23T16:28:40Z
updated_at: 2026-01-23T16:29:49Z
parent: credfolio2-1kt0
blocking:
    - credfolio2-ksna
---

GraphQL mutation to merge reference letter extractions into profile.

## Mutation

```graphql
mutation MergeReferenceLetterIntoProfile(
  profileId: ID!
  referenceLetterExtractionId: ID!
  options: MergeOptions
) {
  profile: Profile!
  changes: [ProfileChange!]!
}

input MergeOptions {
  includeSkills: Boolean = true
  includeTestimonials: Boolean = true
  includeAccomplishments: Boolean = true
}

type ProfileChange {
  type: ChangeType!
  field: String!
  oldValue: String
  newValue: String!
  source: String!
}
```

## Behavior

1. Create profile version snapshot (for undo)
2. Apply changes based on options
3. Record change details for history
4. Return updated profile and change list

## Checklist

- [ ] Define GraphQL types
- [ ] Implement merge resolver
- [ ] Create version snapshot before merge
- [ ] Apply skill additions
- [ ] Apply testimonial additions
- [ ] Record changes for history
- [ ] Write resolver tests