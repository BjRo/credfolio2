---
# credfolio2-zhnh
title: Reference letter upload from profile
status: in-progress
type: feature
priority: normal
created_at: 2026-01-23T16:28:36Z
updated_at: 2026-01-30T09:05:03Z
parent: credfolio2-1kt0
blocking:
    - credfolio2-6dty
---

Button and flow to upload a reference letter from the profile view.

## User Flow

1. User clicks "Add Reference Letter" button on profile page
2. Modal opens with upload zone (drag-drop or file picker)
3. User drops/selects file (PDF, PNG, JPG, DOCX)
4. Modal shows processing state with spinner
5. On completion, transitions to validation preview (credfolio2-6dty)
6. On error, shows error message with retry option

## UI Components

### "Add Reference Letter" Button
- Location: Profile page header or sources section
- Style: Secondary button with icon
- Text: "Add Reference Letter" or "+ Reference"

### Upload Modal
- Reuse/adapt existing FileUpload component
- Drag-drop zone with file type hints
- File size limit indicator
- Processing spinner with status text

## Backend Integration

### GraphQL Mutation
```graphql
mutation uploadReferenceLetter($userId: ID!, $file: Upload!) {
  uploadReferenceLetter(userId: $userId, file: $file) {
    id
    status
  }
}
```

### Polling for Completion
- Poll `referenceLetter(id)` every 2 seconds
- Check status: pending -> processing -> completed/failed
- On completed: navigate to preview
- On failed: show error message

## Checklist

- [x] Add "Add Reference Letter" button to profile page
- [x] Create ReferenceLetterUploadModal component
- [x] Adapt FileUpload component for reference letters
- [x] Wire up uploadReferenceLetter mutation (using `uploadFile` mutation)
- [x] Implement polling for processing status
- [x] Show processing state in modal
- [x] Handle error states with retry
- [ ] Navigate to validation preview on completion (blocked by credfolio2-6dty - currently closes modal and refetches profile)

## Definition of Done

- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] Visual verification with agent-browser (UI changes)
- [ ] All checklist items above are completed (validation preview navigation blocked)
- [x] Branch pushed and PR created for human review: https://github.com/BjRo/credfolio2/pull/44
