---
# credfolio2-wena
title: 'Unify profile data model: materialize extracted data into profile tables'
status: in-progress
type: feature
priority: high
created_at: 2026-01-27T11:58:34Z
updated_at: 2026-01-27T12:07:07Z
parent: credfolio2-v5dw
---

## Problem

Currently there are two separate models for education and work experience:

1. **Extracted data** — stored as JSON in the `resumes.extracted_data` JSONB column, exposed via `Education` and `WorkExperience` GraphQL types (read-only)
2. **Profile data** — stored in `profile_education` and `profile_experiences` tables, exposed via `ProfileEducation` and `ProfileExperience` GraphQL types (editable)

This creates unnecessary complexity:
- Client-side merge logic to deduplicate and combine both sources
- Two separate GraphQL types and query paths per entity
- Fragile matching by institution+degree / company+title keys
- "Edit a resume-extracted item creates a new profile item" edge case

## Proposed Design

**Single model per entity.** When a resume is extracted, immediately materialize `profile_education` and `profile_experiences` rows from the extracted data. Store the original extracted JSON in an `original_data JSONB` column on each row for audit/reference.

This means:
- Frontend only queries `profile.educations` and `profile.experiences` — no merge logic
- The `source` column (already exists) tracks whether an entry is `manual` or `resume_extracted`
- The `source_resume_id` column (already exists) links back to the source resume
- The resume's `extracted_data` JSON remains as an immutable audit trail
- The `Education` / `WorkExperience` GraphQL types become internal to the extraction pipeline

## Checklist

### Backend

- [x] Add `original_data JSONB` column to `profile_education` and `profile_experiences` tables (migration)
- [x] Update resume extraction pipeline to create `profile_education` and `profile_experiences` rows after extraction
  - Set `source = 'resume_extracted'` and `source_resume_id`
  - Store original extracted JSON in `original_data`
  - Create rows via existing repository methods
- [x] Update `Profile` query resolver to no longer need resume extracted data for education/experience display
- [x] Consider keeping `Education` / `WorkExperience` GraphQL types for the extraction result but remove them from profile display path

### Frontend

- [x] Simplify `EducationSection` — remove dual-source merge logic, only consume `profileEducations`
- [x] Simplify `WorkExperienceSection` — remove dual-source merge logic, only consume `profileExperiences`
- [x] Update profile page to stop passing `extractedData.education` / `extractedData.experience`
- [x] Remove the `EducationItem` / unified adapter types that exist only for merging
- [x] Update tests to reflect simplified props

### Cleanup

- [x] Remove unused frontend types (`Education`, `WorkExperience` re-exports from profile types if no longer needed)
- [x] Update GraphQL queries — `GetProfile` becomes the single source for education + experience data

## Definition of Done
- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] Visual verification with agent-browser (for UI changes)
- [x] All other checklist items above are completed