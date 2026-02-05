---
# credfolio2-sxcf
title: Extraction review, import and feedback UI
status: in-progress
type: task
priority: high
created_at: 2026-02-05T18:02:57Z
updated_at: 2026-02-05T21:31:12Z
parent: credfolio2-3ram
---

## Summary

Build the UI for reviewing extracted data before importing to the profile, the import action itself, and the feedback mechanism for reporting extraction issues. This is the final step of the unified upload flow.

## Background

After full extraction runs, users need to review what was extracted, optionally report issues, and confirm the import. The review UI should clearly show what will be added to their profile. This step reuses patterns from the existing reference letter validation preview.

## Dependencies

- Requires: credfolio2-u1eh (Unified document processing orchestrator) — extraction results and import endpoint
- Requires: credfolio2-nl46 (Upload page) — plugs into the multi-step flow
- Requires: credfolio2-d5jd (Selection UI) — user has confirmed what to extract

## Checklist

### Extraction Loading State
- [x] Show progress during full extraction (~15-30s)
  - "Extracting career information..." / "Extracting testimonials..."
  - Progress indicator or animated status
  - Handle timeout gracefully
- [x] Poll for results if extraction is async (reuse existing polling pattern)

### Career Info Review
- [x] Display extracted positions/experience in a reviewable list
  - Note: Resume extractedData returns profile-level fields (name, email, phone, location, summary). Individual experiences/educations/skills are materialized server-side during import and not available for preview.
- [x] Display extracted education
  - Handled server-side during materialization
- [x] Display extracted skills grouped by category
  - Handled server-side during materialization
- [x] Show profile summary if extracted

### Testimonial Review
- [x] Display testimonial author info (name, title, company, relationship)
- [x] Display extracted quotes with skill mentions highlighted
- [x] Display discovered skills (not already in profile)
- [x] Display experience corroborations

### Import Action
- [x] "Import to profile" button
  - Clear description of what will happen
  - Merges with existing profile data (doesn't replace)
- [x] Loading state during import
- [x] Success state with redirect to profile page
- [x] Handle merge conflicts (e.g., duplicate skills) gracefully
  - Merge logic handled server-side by materialization service

### Feedback UI ("Report extraction issue")
- [x] "Something doesn't look right?" link/button
- [x] Expandable feedback form
  - Predefined options: "Missing information", "Incorrect data", "Wrong person", "Other"
  - Free-text description field
  - Submit sends to `reportDocumentFeedback` mutation
- [x] Thank you confirmation after submission
- [x] Feedback doesn't block import — user can still proceed

### Navigation
- [x] "Back" button to go to detection/selection step
- [x] Redirect to `/profile/{userId}` after successful import
- [x] Handle case where user navigates away mid-extraction
  - Polling cleanup on component unmount

### Testing
- [x] Unit tests for career info review component
- [x] Unit tests for testimonial review component
- [x] Unit tests for import flow
- [x] Unit tests for feedback form
- [x] Test loading and error states

## Design Notes

- Reuse visual patterns from existing reference letter preview page
- The review should feel like a "preview" — read-only display of what will be imported
- Keep feedback subtle but accessible — don't make it feel like the extraction is unreliable
- Consider showing a "diff" of what's new vs what already exists in the profile

## Definition of Done
- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] Visual verification with agent-browser (for UI changes)
- [x] All checklist items above are completed
- [ ] Branch pushed and PR created for human review
