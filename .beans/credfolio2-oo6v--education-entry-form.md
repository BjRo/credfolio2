---
# credfolio2-oo6v
title: Education entry form
status: draft
type: feature
created_at: 2026-01-26T11:43:07Z
updated_at: 2026-01-26T11:43:07Z
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

- [ ] Add GraphQL mutation for education CRUD
- [ ] Create EducationForm component
- [ ] Add date pickers with validation
- [ ] Implement add mutation
- [ ] Implement edit mutation
- [ ] Implement delete with confirmation
- [ ] Add form trigger to education section
- [ ] Handle validation errors