---
# credfolio2-sxcf
title: Extraction review, import and feedback UI
status: draft
type: task
priority: high
created_at: 2026-02-05T18:02:57Z
updated_at: 2026-02-05T18:02:57Z
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
- [ ] Show progress during full extraction (~15-30s)
  - "Extracting career information..." / "Extracting testimonials..."
  - Progress indicator or animated status
  - Handle timeout gracefully
- [ ] Poll for results if extraction is async (reuse existing polling pattern)

### Career Info Review
- [ ] Display extracted positions/experience in a reviewable list
  - Company, title, dates, description, highlights
  - Visual indicator for each item
- [ ] Display extracted education
  - Institution, degree, field, dates
- [ ] Display extracted skills grouped by category
  - Technical, Soft, Domain skills
- [ ] Show profile summary if extracted

### Testimonial Review
- [ ] Display testimonial author info (name, title, company, relationship)
- [ ] Display extracted quotes with skill mentions highlighted
- [ ] Display discovered skills (not already in profile)
- [ ] Display experience corroborations

### Import Action
- [ ] "Import to profile" button
  - Clear description of what will happen
  - Merges with existing profile data (doesn't replace)
- [ ] Loading state during import
- [ ] Success state with redirect to profile page
- [ ] Handle merge conflicts (e.g., duplicate skills) gracefully
  - Show what was merged vs what was new

### Feedback UI ("Report extraction issue")
- [ ] "Something doesn't look right?" link/button
- [ ] Expandable feedback form
  - Predefined options: "Missing information", "Incorrect data", "Wrong person", "Other"
  - Free-text description field
  - Submit sends to `reportDocumentFeedback` mutation
- [ ] Thank you confirmation after submission
- [ ] Feedback doesn't block import — user can still proceed

### Navigation
- [ ] "Back" button to go to detection/selection step
- [ ] Redirect to `/profile/{resumeId}` after successful import
- [ ] Handle case where user navigates away mid-extraction

### Testing
- [ ] Unit tests for career info review component
- [ ] Unit tests for testimonial review component
- [ ] Unit tests for import flow
- [ ] Unit tests for feedback form
- [ ] Test loading and error states

## Design Notes

- Reuse visual patterns from existing reference letter preview page
- The review should feel like a "preview" — read-only display of what will be imported
- Keep feedback subtle but accessible — don't make it feel like the extraction is unreliable
- Consider showing a "diff" of what's new vs what already exists in the profile

## Definition of Done
- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Visual verification with agent-browser (for UI changes)
- [ ] All checklist items above are completed
- [ ] Branch pushed and PR created for human review