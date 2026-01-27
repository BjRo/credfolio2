---
# credfolio2-oo6v
title: Education entry form
status: completed
type: feature
priority: high
created_at: 2026-01-26T11:43:07Z
updated_at: 2026-01-27T13:03:08Z
parent: credfolio2-v5dw
---

Form for adding/editing education entries.

## Form Fields

- Institution name (required)
- Degree/certification (required)
- Field of study
- Start date (month/year picker)
- End date (month/year, or "in progress" checkbox)
- Description/achievements
- GPA (optional)

## Use Cases

1. **Add new education** - Empty form, "Add" button
2. **Edit existing** - Pre-filled form, "Save" button
3. **Delete** - Confirmation dialog

## EducationForm Component

- Modal or slide-over panel (consistent with ExperienceForm)
- Validation with inline errors
- Date pickers with reasonable constraints

## Checklist

### Backend

- [x] Create `profile_education` database migration (model after `profile_experiences`)
- [x] Add `ProfileEducation` domain type in `internal/domain/profile.go`
- [x] Create `profile_education_repository.go` with CRUD operations
- [x] Add GraphQL input types: `CreateEducationInput`, `UpdateEducationInput`
- [x] Add GraphQL mutations: `createEducation`, `updateEducation`, `deleteEducation`
- [x] Add resolver implementations for education mutations

### Frontend

- [x] Generate GraphQL types after backend is ready
- [x] Create `EducationForm` component (modal, consistent with `WorkExperienceForm`)
- [x] Create `EducationFormDialog` wrapper component
- [x] Reuse `MonthYearPicker` for date fields
- [x] Implement add education flow
- [x] Implement edit education flow (pre-fill form)
- [x] Implement delete with `DeleteEducationDialog`
- [x] Add "Add Education" button trigger to education section
- [x] Handle validation errors inline

## Definition of Done
- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] Visual verification with agent-browser (for UI changes)
- [x] All other checklist items above are completed

## Reference

- Follow patterns from `WorkExperienceForm` and `WorkExperienceFormDialog`
- Database table structure should mirror `profile_experiences`
- Use existing `Education` GraphQL type fields as guide for form fields