---
# credfolio2-wupt
title: Experience entry form
status: completed
type: feature
priority: normal
created_at: 2026-01-23T16:29:28Z
updated_at: 2026-01-26T16:35:30Z
parent: credfolio2-v5dw
---

Form for adding/editing work experience entries.

## Form Fields

- Company name (required)
- Job title (required)
- Location
- Start date (month/year picker)
- End date (month/year, or "current" checkbox)
- Description (rich text or plain)
- Highlights (bullet point editor)

## Use Cases

1. **Add new experience** - Empty form, "Add" button
2. **Edit existing** - Pre-filled form, "Save" button
3. **Delete** - Confirmation dialog

## ExperienceForm Component

- Modal or slide-over panel
- Validation with inline errors
- Date pickers with reasonable constraints
- Highlights editor (add/remove/reorder bullets)

## Checklist

### Phase 1: Backend
- [x] Create database migration for profiles and profile_experiences tables
- [x] Add Profile and ProfileExperience domain entities
- [x] Add ProfileRepository and ProfileExperienceRepository interfaces
- [x] Implement Postgres repositories
- [x] Add GraphQL schema (input types, mutations, Profile query)
- [x] Implement resolvers for experience CRUD

### Phase 2: Frontend UI Components
- [x] Add Radix UI Dialog dependency
- [x] Create shadcn/ui-style dialog, input, label, textarea, checkbox
- [x] Create MonthYearPicker component
- [x] Create HighlightsEditor component

### Phase 3: Experience Form
- [x] Create WorkExperienceForm component
- [x] Create WorkExperienceFormDialog (modal wrapper)
- [x] Create DeleteExperienceDialog (confirmation)
- [x] Add GraphQL mutations to frontend

### Phase 4: Integration
- [x] Update WorkExperienceSection with edit triggers
- [x] Connect form to mutations
- [x] End-to-end testing (build passes, backend tests pass)