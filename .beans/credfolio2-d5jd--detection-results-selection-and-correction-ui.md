---
# credfolio2-d5jd
title: Detection results, selection and correction UI
status: in-progress
type: task
priority: high
created_at: 2026-02-05T18:02:39Z
updated_at: 2026-02-05T21:18:43Z
parent: credfolio2-3ram
blocking:
    - credfolio2-sxcf
---

## Summary

Build the UI components for displaying detection results, letting users select what to extract, correcting misdetections, and handling low-confidence scenarios. This is the decision-making step between detection and full extraction.

## Background

After the lightweight detection scan, users need to see what was found, confirm or correct it, and choose what to extract. This step is critical for user trust — they should feel in control of what happens with their document.

## Dependencies

- Requires: credfolio2-nl46 (Upload page with drop zone) — this UI plugs into the multi-step flow
- Requires: credfolio2-4h8a (Detection service) — detection result types

## Checklist

### Detection Results Display
- [x] Content type checkboxes showing what was detected
  - "Career Information (resume/CV content)" — checked by default if detected
  - "Testimonial from [Author Name]" — checked by default if detected
  - User can uncheck to skip extraction of that content type
- [x] Document summary text from detection
- [x] Testimonial author name display when detected

### Confidence Handling
- [x] High confidence (≥0.7): Show checkboxes pre-selected, no warning
- [x] Low confidence (<0.7): Show warning banner + force manual selection
  - "We're not sure what this document contains. Please tell us:"
  - Radio options: "Resume / CV", "Reference letter", "Both"
  - Free-text fallback: "Something else: ___"
- [x] Confidence threshold should be configurable (default 0.7)

### Correction UI ("Not what you expected?")
- [x] Expandable/collapsible correction section below detection results
- [x] Quick correction options:
  - "This is just a resume" (unchecks testimonial)
  - "This is just a reference letter" (unchecks career info)
- [x] Free-text option: "Tell us more about this document: ___"
- [x] Corrections should log feedback to backend (for improving detection)

### Error States
- [x] Unreadable/corrupted document: falls through to low confidence flow
  - Show "Upload different document" button that returns to upload step
- [x] Empty detection (nothing found): Similar to low confidence — ask user to classify manually

### Proceed Action
- [x] "Extract selected content" button to proceed to extraction
  - Disabled if nothing is selected
  - Shows what will happen: "We'll extract career information and testimonials from this document"
- [x] "Cancel" or "Upload different document" option

### Testing
- [x] Unit tests for detection results component
- [x] Test high confidence rendering (checkboxes, no warning)
- [x] Test low confidence rendering (warning, manual selection)
- [x] Test correction UI interactions
- [x] Test error states
- [x] Test proceed button states (enabled/disabled)

## Design Notes

- Keep the UI clean and simple for the happy path (high confidence)
- Only surface complexity when needed (low confidence, corrections)
- Use existing Checkbox, Button, Input components from the UI library
- The correction section should be subtle — most users won't need it

## Definition of Done
- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] Visual verification with agent-browser (for UI changes)
- [x] All checklist items above are completed
- [ ] Branch pushed and PR created for human review