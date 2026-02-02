---
# credfolio2-ucw3
title: Add source badge to testimonials
status: completed
type: feature
priority: normal
created_at: 2026-01-30T13:47:08Z
updated_at: 2026-02-02T13:39:11Z
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
- [x] Update testimonials query to include `referenceLetter { file { url } }`
- [x] Design source badge component (document icon + tooltip)
- [x] Add badge to TestimonialCard component
- [x] Link badge to PDF URL (opens in new tab)
- [x] Handle case where reference letter has no file attached (hide badge)

## Notes
- This feature is **independent** of the Author entity work
- Can be implemented immediately

## Definition of Done
- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] Visual verification with agent-browser (for UI changes)
  - Profile page renders correctly
  - Testimonials section shows empty state (no fixture reference letter available)
  - Source badge tested via unit tests (renders when file URL exists)
- [x] All other checklist items above are completed
- [x] Branch pushed and PR created for human review: https://github.com/BjRo/credfolio2/pull/54