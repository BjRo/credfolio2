---
# credfolio2-wupt
title: Experience entry form
status: draft
type: feature
created_at: 2026-01-23T16:29:28Z
updated_at: 2026-01-23T16:29:28Z
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

- [ ] Add GraphQL mutations for experience CRUD (add, update, delete)
- [ ] Create ExperienceForm component
- [ ] Add date pickers with validation
- [ ] Create highlights bullet editor
- [ ] Implement add experience (connect to mutation)
- [ ] Implement edit experience (connect to mutation)
- [ ] Implement delete with confirmation (connect to mutation)
- [ ] Add form trigger to experience section