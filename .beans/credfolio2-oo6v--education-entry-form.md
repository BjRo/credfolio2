---
# credfolio2-oo6v
title: Education entry form
status: todo
type: feature
priority: high
created_at: 2026-01-26T11:43:07Z
updated_at: 2026-01-26T16:41:06Z
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

- [ ] Create `profile_education` database migration (model after `profile_experiences`)
- [ ] Add `ProfileEducation` domain type in `internal/domain/profile.go`
- [ ] Create `profile_education_repository.go` with CRUD operations
- [ ] Add GraphQL input types: `CreateEducationInput`, `UpdateEducationInput`
- [ ] Add GraphQL mutations: `createEducation`, `updateEducation`, `deleteEducation`
- [ ] Add resolver implementations for education mutations

### Frontend

- [ ] Generate GraphQL types after backend is ready
- [ ] Create `EducationForm` component (modal, consistent with `WorkExperienceForm`)
- [ ] Create `EducationFormDialog` wrapper component
- [ ] Reuse `MonthYearPicker` for date fields
- [ ] Implement add education flow
- [ ] Implement edit education flow (pre-fill form)
- [ ] Implement delete with `DeleteEducationDialog`
- [ ] Add "Add Education" button trigger to education section
- [ ] Handle validation errors inline

## Reference

- Follow patterns from `WorkExperienceForm` and `WorkExperienceFormDialog`
- Database table structure should mirror `profile_experiences`
- Use existing `Education` GraphQL type fields as guide for form fields