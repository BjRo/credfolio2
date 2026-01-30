---
# credfolio2-ucw3
title: Add source badge to testimonials
status: draft
type: feature
priority: normal
created_at: 2026-01-30T13:47:08Z
updated_at: 2026-01-30T13:50:46Z
parent: credfolio2-2ex3
---

Show which reference letter each testimonial came from to build trust and enable source verification.

## User Story
As a profile viewer, I want to see where each testimonial came from so I can trust that these are real recommendations from actual documents.

## Implementation

### UI Design
- Small document icon/badge on each testimonial card
- Display the reference letter title or author name
- Position: top-right corner of card or near the quote
- **Click behavior: Open the raw PDF file in a new tab**

### Data Available
The `Testimonial` type already has a `referenceLetter` field in GraphQL:
```graphql
type Testimonial {
  referenceLetter: ReferenceLetter
}

type ReferenceLetter {
  id: ID!
  extractedAuthor: ExtractedAuthor
  file: File  # Contains the uploaded PDF
}

type File {
  id: ID!
  url: String!  # Presigned URL to download
}
```

### Tasks
- [ ] Update testimonials query to include `referenceLetter { file { url } }`
- [ ] Design source badge component (document icon + tooltip)
- [ ] Add badge to TestimonialCard component
- [ ] Link badge to PDF URL (opens in new tab)
- [ ] Handle case where reference letter has no file attached (hide badge)

## Notes
- This feature is **independent** of the Author entity work
- Can be implemented immediately

## Definition of Done
- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Visual verification with agent-browser (for UI changes)
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review