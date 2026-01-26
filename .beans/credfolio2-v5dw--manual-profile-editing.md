---
# credfolio2-v5dw
title: Manual Profile Editing
status: todo
type: epic
priority: normal
created_at: 2026-01-23T16:27:00Z
updated_at: 2026-01-26T11:43:52Z
parent: credfolio2-dwid
---

Allow users to directly edit any part of their profile.

## Editable Fields

- Personal info (name, headline, summary)
- Work experience (add, edit, remove)
- Skills (add, edit, remove, recategorize)
- Education (add, edit, remove)
- Any text field from extractions

## Edit Modes

### Inline Editing
- Click on any text to edit in place
- Auto-save with debounce
- Clear save/cancel affordance

### Form-based Editing
- Full edit form for complex items (positions)
- Add new items via form
- Delete with confirmation

## Validation

- Required fields highlighted
- Date validation (end after start)
- Character limits where appropriate

## Out of Scope

- Drag-and-drop reordering (deferred)
- Change history and undo (see credfolio2-09w1)