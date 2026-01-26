---
# credfolio2-wupt
title: Experience entry form
status: in-progress
type: feature
priority: normal
created_at: 2026-01-23T16:29:28Z
updated_at: 2026-01-26T11:46:44Z
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
- [ ] Add Radix UI Dialog dependency
- [ ] Create shadcn/ui-style dialog, input, label, textarea, checkbox
- [ ] Create MonthYearPicker component
- [ ] Create HighlightsEditor component

### Phase 3: Experience Form
- [ ] Create WorkExperienceForm component
- [ ] Create WorkExperienceFormDialog (modal wrapper)
- [ ] Create DeleteExperienceDialog (confirmation)
- [ ] Add GraphQL mutations to frontend

### Phase 4: Integration
- [ ] Update WorkExperienceSection with edit triggers
- [ ] Connect form to mutations
- [ ] End-to-end testing