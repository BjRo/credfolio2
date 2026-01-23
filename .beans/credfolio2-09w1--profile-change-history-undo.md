---
# credfolio2-09w1
title: Profile Change History & Undo
status: draft
type: epic
created_at: 2026-01-23T16:26:59Z
updated_at: 2026-01-23T16:26:59Z
parent: credfolio2-dwid
---

Track all changes to the profile and allow undoing any modification.

## Features

### Change History
- Every extraction and edit creates a history entry
- View timeline of changes
- See what changed at each point

### Undo Capabilities
- Undo last change (single step)
- Undo entire reference letter enhancement (bulk undo)
- Restore to any previous point in history

### Manual Post-fix
- After accepting reference letter changes, can still edit
- Edits are tracked as their own history entries
- Clear distinction between LLM-extracted and human-edited

## Data Model

```
profile_versions
- id, profile_id, version_number
- snapshot (full profile state)
- change_type (extraction, manual_edit, undo)
- change_source (resume, reference_letter_id, user)
- created_at
```

## UI

- History panel/drawer accessible from profile
- "Undo" button prominently available after changes
- Visual indicator when viewing non-current version