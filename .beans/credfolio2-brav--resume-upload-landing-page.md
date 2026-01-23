---
# credfolio2-brav
title: Resume upload landing page
status: draft
type: feature
created_at: 2026-01-23T16:27:28Z
updated_at: 2026-01-23T16:27:28Z
parent: credfolio2-6oza
---

Build the first screen users see: a clean upload interface for resumes.

## Design

- Full-page upload zone (not small dropzone in corner)
- Clear headline: "Upload your resume to get started"
- Drag-and-drop with visual feedback
- Click-to-browse fallback
- Supported formats: PDF, DOCX
- File size limit indicator

## Behavior

1. User drops/selects file
2. Validate file type client-side
3. Start upload immediately
4. Transition to processing view (skeleton)

## Technical

- Route: `/` (homepage)
- Use existing FileUpload component as base
- Style with shadcn/ui + custom styling
- No authentication required initially

## Checklist

- [ ] Create landing page layout
- [ ] Adapt FileUpload component for full-page use
- [ ] Add file validation (type, size)
- [ ] Connect to upload mutation
- [ ] Transition to processing/skeleton view
- [ ] Add error handling for failed uploads