---
# credfolio2-h2i8
title: Skills management UI
status: draft
type: feature
created_at: 2026-01-23T16:29:31Z
updated_at: 2026-01-23T16:29:31Z
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

## Checklist

- [ ] Add GraphQL mutations for skill CRUD (add, update, delete)
- [ ] Add "add skill" UI per category
- [ ] Add remove button to skill tags
- [ ] Implement skill rename (connect to mutation)
- [ ] Implement category change (connect to mutation)
- [ ] Handle skills with sources carefully (warn on delete)

## Out of Scope

- Drag-and-drop reordering (deferred)