---
# credfolio2-ue4q
title: Reference letter data model and extraction
status: completed
type: task
priority: normal
created_at: 2026-01-23T16:28:38Z
updated_at: 2026-01-30T08:58:36Z
parent: credfolio2-1kt0
blocking:
    - credfolio2-6dty
    - credfolio2-1twf
---

Foundation for the reference letter credibility system: database schema, domain types, and LLM extraction.

## Database Schema

### reference_letters table
- `id` UUID PRIMARY KEY
- `user_id` UUID REFERENCES users
- `file_id` UUID REFERENCES files
- `status` TEXT (pending, processing, completed, failed)
- `extracted_data` JSONB (raw extraction results)
- `error_message` TEXT (if failed)
- `created_at`, `updated_at` TIMESTAMPTZ

### testimonials table
- `id` UUID PRIMARY KEY
- `profile_id` UUID REFERENCES profiles
- `reference_letter_id` UUID REFERENCES reference_letters
- `quote` TEXT NOT NULL
- `author_name` TEXT NOT NULL
- `author_title` TEXT
- `author_company` TEXT
- `relationship` TEXT (manager, peer, direct_report, client, other)
- `created_at`, `updated_at` TIMESTAMPTZ

### skill_validations table
- `id` UUID PRIMARY KEY
- `profile_skill_id` UUID REFERENCES profile_skills
- `reference_letter_id` UUID REFERENCES reference_letters
- `quote_snippet` TEXT (short quote for hover popover)
- `created_at` TIMESTAMPTZ
- UNIQUE (profile_skill_id, reference_letter_id)

### experience_validations table
- `id` UUID PRIMARY KEY
- `profile_experience_id` UUID REFERENCES profile_experiences
- `reference_letter_id` UUID REFERENCES reference_letters
- `quote_snippet` TEXT
- `created_at` TIMESTAMPTZ
- UNIQUE (profile_experience_id, reference_letter_id)

## Extraction Schema (LLM Output)

```json
{
  "author": {
    "name": "John Smith",
    "title": "Engineering Manager",
    "company": "Acme Corp",
    "relationship": "manager"
  },
  "testimonials": [
    {
      "quote": "Jane's leadership during our cloud migration...",
      "skillsMentioned": ["leadership", "cloud architecture"]
    }
  ],
  "skillMentions": [
    {
      "skill": "Go",
      "quote": "Her expertise in Go helped us...",
      "context": "technical skills"
    }
  ],
  "experienceMentions": [
    {
      "company": "Acme Corp",
      "role": "Senior Engineer",
      "quote": "During her time as Senior Engineer..."
    }
  ],
  "discoveredSkills": ["mentoring", "system design"]
}
```

## GraphQL Types

- `ReferenceLetter` type with status, extractedData
- `Testimonial` type
- `SkillValidation` type
- `ExperienceValidation` type
- `uploadReferenceLetter` mutation
- `referenceLetter(id)` query
- `referenceLetters(userId)` query

## Checklist

- [x] Create `reference_letters` migration (added error_message column)
- [x] Create `testimonials` migration
- [x] Create `skill_validations` migration
- [x] Create `experience_validations` migration
- [x] Add Go domain types for all entities
- [x] Add repository methods for CRUD operations
- [x] Add GraphQL types to schema
- [x] Create extraction prompt for reference letters
- [x] Create ReferenceLetterProcessingJob (River job)
- [x] Implement text extraction (reuse from resume)
- [x] Implement structured extraction with JSON schema
- [x] Store extracted_data in reference_letters table
- [x] Write unit tests for extraction
- [x] Write integration tests for job processing

## Definition of Done

- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] All checklist items above are completed
- [x] Branch pushed and PR created for human review (PR #43)
