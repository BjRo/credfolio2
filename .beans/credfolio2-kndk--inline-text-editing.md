---
# credfolio2-kndk
title: Inline text editing
status: draft
type: feature
priority: normal
created_at: 2026-01-23T16:29:27Z
updated_at: 2026-01-28T14:50:56Z
parent: credfolio2-dwid
---

Click on text fields to edit them in place.

## Editable Fields

- Name
- Headline
- Summary
- Location fields
- Contact info (email, phone)
- Experience descriptions
- Skill names

## EditableText Component

Props:
- value: string
- onSave: (newValue) => void
- multiline?: boolean
- placeholder?: string

States:
- View mode: displays text, shows edit icon on hover
- Edit mode: input/textarea, save/cancel buttons
- Saving: disabled with spinner
- Error: shows validation message

## Behavior

- Click text or edit icon to enter edit mode
- Escape to cancel
- Enter to save (Shift+Enter for multiline newline)
- Blur saves (debounced)
- Optimistic update with rollback on error

## Checklist

- [ ] Add GraphQL mutation for updating profile text fields
- [ ] Create EditableText component
- [ ] Create EditableTextArea component (multiline)
- [ ] Add edit mode styling
- [ ] Implement save handlers (connect to mutation)
- [ ] Add to profile header fields
- [ ] Add to experience cards
- [ ] Handle validation errors