---
# credfolio2-zhnh
title: Add reference letter upload from profile
status: draft
type: feature
priority: normal
created_at: 2026-01-23T16:28:36Z
updated_at: 2026-01-23T16:29:49Z
parent: credfolio2-1kt0
blocking:
    - credfolio2-6dty
---

Button and flow to upload a reference letter from the profile view.

## Trigger

- "Add Reference Letter" button on profile page
- Opens modal or drawer with upload interface

## Upload Interface

- Similar to initial upload (drag-drop)
- File types: PDF, images (PNG, JPG), DOCX
- Processing indicator inline

## Flow

1. User clicks "Add Reference Letter"
2. Modal opens with upload zone
3. User drops/selects file
4. Modal shows processing state
5. On completion, modal shows enhancement preview
6. User reviews and accepts/rejects (next ticket)

## Checklist

- [ ] Add "Add Reference Letter" button to profile
- [ ] Create upload modal/drawer
- [ ] Reuse/adapt FileUpload component
- [ ] Trigger reference letter extraction job
- [ ] Show processing state in modal
- [ ] Transition to preview when done