---
# credfolio2-h2i8
title: Skills management UI
status: in-progress
type: feature
priority: high
created_at: 2026-01-23T16:29:31Z
updated_at: 2026-01-27T17:58:24Z
parent: credfolio2-v5dw
---

Add, edit, remove, and categorize skills.

## Features

### Add Skill
- Input field with category dropdown
- Autocomplete from common skills (nice-to-have)
- Enter to add quickly

### Remove Skill
- X button on skill tag
- Confirmation for skills from documents

### Edit Skill
- Click to rename
- Change category via dropdown

## SkillManager Component

- Shows all skills by category
- Add skill input per category or global
- Inline editing on click

## Data Model Decision

Manual skills need a `profile_skills` table to coexist with extraction-sourced skills. Each skill should track:
- Name and normalized name (for deduplication)
- Category (using existing `SkillCategory` enum)
- Source: `manual` vs `extraction` (with optional source_resume_id)
- Display order within category

## Checklist

### Backend

- [ ] Create `profile_skills` database migration with columns: id, profile_id, name, normalized_name, category, source, source_resume_id, display_order, created_at, updated_at
- [ ] Add `ProfileSkill` domain type in `internal/domain/profile.go`
- [ ] Create `profile_skills_repository.go` with CRUD operations
- [ ] Add GraphQL input types: `CreateSkillInput`, `UpdateSkillInput`
- [ ] Add GraphQL mutations: `createSkill`, `updateSkill`, `deleteSkill`
- [ ] Add resolver implementations for skill mutations
- [ ] Update Profile query to merge manual skills with extraction skills

### Frontend

- [ ] Generate GraphQL types after backend is ready
- [ ] Create `AddSkillInput` component (inline input with category dropdown)
- [ ] Add "+" button per category section to trigger add input
- [ ] Add "×" remove button to skill tags
- [ ] Implement skill rename (click to edit inline)
- [ ] Implement category change via dropdown
- [ ] Show warning dialog when deleting extraction-sourced skills
- [ ] Handle optimistic updates for snappy UX

## Out of Scope

- Drag-and-drop reordering (deferred)
- Skill autocomplete/suggestions (nice-to-have, separate task)