---
# credfolio2-m607
title: Create Author entity for testimonials
status: draft
type: task
priority: high
created_at: 2026-01-30T13:50:32Z
updated_at: 2026-01-30T13:52:47Z
parent: credfolio2-2ex3
blocking:
    - credfolio2-cm02
    - credfolio2-lvdo
---

Extract author information into a proper entity with identity, enabling LinkedIn linking, editing, and deduplication.

## Context
Currently, author info is stored as denormalized text fields directly on the `testimonials` table:
- `author_name` (TEXT)
- `author_title` (TEXT)  
- `author_company` (TEXT)
- `relationship` (enum)

This makes it impossible to:
- Link the same author across multiple testimonials
- Edit author info in one place
- Add LinkedIn URL to an author (would need to duplicate across testimonials)
- Properly group testimonials by author

## Solution
Create a separate `authors` table and link testimonials to it.

### Database Schema
```sql
CREATE TABLE authors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    profile_id UUID NOT NULL REFERENCES profiles(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    title TEXT,
    company TEXT,
    linkedin_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(profile_id, name, company)  -- prevent obvious duplicates
);

-- Add foreign key to testimonials (relationship STAYS on testimonial)
ALTER TABLE testimonials 
    ADD COLUMN author_id UUID REFERENCES authors(id),
    -- Keep old columns temporarily for migration
    ALTER COLUMN author_name DROP NOT NULL;

-- Migrate existing data
INSERT INTO authors (profile_id, name, title, company)
SELECT DISTINCT t.profile_id, t.author_name, t.author_title, t.author_company
FROM testimonials t
WHERE t.author_id IS NULL;

UPDATE testimonials t
SET author_id = a.id
FROM authors a
WHERE t.profile_id = a.profile_id 
  AND t.author_name = a.name 
  AND COALESCE(t.author_company, '') = COALESCE(a.company, '');

-- Later migration: drop old columns (keep relationship!)
ALTER TABLE testimonials
    DROP COLUMN author_name,
    DROP COLUMN author_title,
    DROP COLUMN author_company;
    -- NOTE: relationship column STAYS on testimonials
```

### Domain Model
```go
type Author struct {
    ID          uuid.UUID
    ProfileID   uuid.UUID
    Name        string
    Title       *string
    Company     *string
    LinkedInURL *string
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

type Testimonial struct {
    // ... existing fields ...
    AuthorID     uuid.UUID
    Author       *Author       // relation
    Relationship TestimonialRelationship  // STAYS HERE - same person can have different relationships
}
```

### Design Decision: Relationship stays on Testimonial
The `relationship` field remains on `Testimonial`, not `Author`, because:
- Same person could be your manager at one company and peer at another
- Each testimonial captures a specific context/relationship

### Extraction Flow Change
When applying reference letter validations:
1. Extract author info from reference letter
2. Check if author with same name+company exists for this profile
3. If exact match: reuse existing author
4. If multiple potential matches: return candidates, let user choose (future enhancement)
5. If no match: create new author

### GraphQL Changes
```graphql
type Author {
    id: ID!
    name: String!
    title: String
    company: String
    linkedInUrl: String
    testimonials: [Testimonial!]!
}

type Testimonial {
    # Remove: authorName, authorTitle, authorCompany
    author: Author!
    relationship: TestimonialRelationship!  # STAYS HERE
}

input UpdateAuthorInput {
    name: String
    title: String
    company: String
    linkedInUrl: String
}

mutation updateAuthor(id: ID!, input: UpdateAuthorInput!): Author!
```

## Tasks
- [ ] Create migration for `authors` table
- [ ] Create migration to add `author_id` to testimonials
- [ ] Create data migration to populate authors from existing testimonials
- [ ] Update domain model (`Author` struct, update `Testimonial`)
- [ ] Create `AuthorRepository` interface and implementation
- [ ] Update `TestimonialRepository` to handle author relation
- [ ] Update GraphQL schema (Author type, mutations)
- [ ] Update resolvers
- [ ] Update extraction flow to find/create authors
- [ ] Update frontend queries and components

## Definition of Done
- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review