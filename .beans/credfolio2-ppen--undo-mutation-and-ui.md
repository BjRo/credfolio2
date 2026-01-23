---
# credfolio2-ppen
title: Undo mutation and UI
status: draft
type: feature
created_at: 2026-01-23T16:29:25Z
updated_at: 2026-01-23T16:29:25Z
parent: credfolio2-09w1
---

Allow users to undo recent changes to their profile.

## GraphQL

```graphql
mutation UndoLastChange(profileId: ID!): Profile!
mutation RestoreVersion(profileId: ID!, versionId: ID!): Profile!
query ProfileHistory(profileId: ID!): [ProfileVersion!]!
```

## UI Components

### UndoButton
- Shows after any change
- "Undo" with brief description of what will be undone
- Confirmation for destructive undos

### HistoryPanel
- Drawer/modal showing version timeline
- Click any version to preview
- "Restore" button per version
- Shows change type and description

## Behavior

- Undo creates a new version (for redo capability)
- Clear feedback on what was undone
- Toast notification on successful undo

## Checklist

- [ ] Implement UndoLastChange mutation
- [ ] Implement RestoreVersion mutation
- [ ] Implement ProfileHistory query
- [ ] Create UndoButton component
- [ ] Create HistoryPanel component
- [ ] Add history access to profile page
- [ ] Add confirmation dialogs