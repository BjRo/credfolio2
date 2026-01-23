---
# credfolio2-v5dw
title: Manual Profile Editing
status: draft
type: epic
created_at: 2026-01-23T16:27:00Z
updated_at: 2026-01-23T16:27:00Z
parent: credfolio2-dwid
---

Allow users to directly edit any part of their profile.

## Editable Fields

- Personal info (name, headline, summary)
- Work experience (add, edit, remove, reorder positions)
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

### Drag-and-Drop
- Reorder positions
- Reorder skills within categories

## Validation

- Required fields highlighted
- Date validation (end after start)
- Character limits where appropriate

## Change Tracking

- All edits create history entries
- Can undo individual edits