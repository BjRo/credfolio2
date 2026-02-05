---
# credfolio2-nl46
title: Upload page with drop zone and detection flow
status: in-progress
type: task
priority: high
created_at: 2026-02-05T18:02:20Z
updated_at: 2026-02-05T21:00:27Z
parent: credfolio2-3ram
blocking:
    - credfolio2-d5jd
    - credfolio2-sxcf
---

## Summary

Create the unified `/upload` page with document drop zone and detection flow. This is the first frontend piece — handling file upload through to showing detection results.

## Background

The existing `ResumeUpload` component provides a solid pattern: drag-and-drop, file validation, progress tracking, and polling. This task creates a new page that uploads a document and runs lightweight detection, presenting results to the user.

## Dependencies

- Requires: credfolio2-4h8a (Document content detection service) — the backend detection endpoint

## Checklist

### Page & Routing
- [ ] Create `/upload` route at `src/frontend/src/app/upload/page.tsx`
- [ ] Add navigation link to `/upload` in site header (or appropriate location)

### Document Drop Zone Component
- [ ] Create `DocumentUpload` component (can adapt patterns from existing `ResumeUpload`)
  - Drag-and-drop zone with visual feedback
  - File validation: PDF, DOCX, TXT (max 10MB) — same as existing
  - Upload progress indicator
  - Accept single file at a time
- [ ] Wire up to `detectDocumentContent` GraphQL mutation
  - Upload file
  - Show "Analyzing document..." loading state during detection
  - Display detection results when ready

### Multi-Step Flow Container
- [ ] Create step indicator/progress component for the multi-step flow
  - Steps: Upload → Review Detection → Extract → Review Results → Import
  - Visual indicator of current step
  - Keep it simple — no need for complex wizard framework
- [ ] Manage flow state with React useState (current step, detection results, etc.)

### Loading States
- [ ] Upload progress bar (reuse existing pattern)
- [ ] Detection analysis spinner/skeleton
- [ ] Clear messaging at each stage ("Uploading...", "Analyzing document content...")

### Error Handling
- [ ] File validation errors (wrong type, too large)
- [ ] Upload failure
- [ ] Detection failure / unreadable document
  - Show message: "We couldn't read this document. Please try uploading a clearer version."
- [ ] Duplicate file detection (reuse existing pattern)

### Testing
- [ ] Unit tests for DocumentUpload component
- [ ] Test file validation
- [ ] Test loading states and error states
- [ ] Test flow navigation

## Design Notes

- Follow existing upload component patterns for consistency
- Use existing shadcn-style UI components (Button, Dialog, etc.)
- The drop zone should feel similar to the existing resume upload for UX consistency
- Mobile-friendly: the drop zone should work well on small screens

## Definition of Done
- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Visual verification with agent-browser (for UI changes)
- [ ] All checklist items above are completed
- [ ] Branch pushed and PR created for human review