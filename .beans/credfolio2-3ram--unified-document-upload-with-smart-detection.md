---
# credfolio2-3ram
title: Unified document upload with smart detection
status: draft
type: feature
created_at: 2026-02-03T12:34:56Z
updated_at: 2026-02-03T12:34:56Z
---

## Summary

Create a unified document upload flow that uses lightweight AI detection to identify content types (career info, testimonials) and lets users choose what to extract. This replaces the need for separate resume and reference letter upload paths.

## Background

Reference letters often contain career information (positions, skills) in addition to testimonials. Currently, users would need to upload the same document twice to extract both types of content. This feature creates a smarter, unified upload experience.

### Design Decisions Made

1. **Lightweight detection first** - Quick LLM scan to classify content before full extraction (cheaper, faster)
2. **Unified upload page** - Single entry point at `/upload` for all document types
3. **User correction** - Allow users to override detection if wrong + report extraction issues
4. **Incremental rollout** - Keep existing separate flows initially, deprecate later if unified works well
5. **Merge behavior** - Duplicate skills/positions from different documents are merged, not duplicated
6. **Confidence display** - Only surface confidence to users when it's low (clean UI for happy path, warning when uncertain)
7. **Error handling**:
   - Low confidence → Force manual selection ("What does this document contain?")
   - Unreadable/corrupted → Ask user to try a different document

## User Flow

```
Upload document
       ↓
Lightweight scan (classify content types)
       ↓
Show detection results with checkboxes:
  ☑ Career Information
  ☑ Testimonial (from Jane Doe)
  + "Not what you expected?" correction options
       ↓
User confirms selection → Full extraction runs
       ↓
Review extracted data
  + "Report extraction issue" feedback option
       ↓
Import to profile (merge with existing data)
```

## Checklist

### Backend

- [ ] Create lightweight document detection prompt
  - Input: document text/content
  - Output: has_career_info, has_testimonial, testimonial_author, confidence, summary
- [ ] Create `/api/documents/detect` endpoint for lightweight detection
- [ ] Create `/api/documents/upload` unified upload endpoint
  - Accepts document + user's extraction preferences
  - Runs appropriate extractors based on selection
  - Returns extracted data for review
- [ ] Create `/api/documents/import` endpoint to save reviewed data to profile
- [ ] Add feedback logging for detection/extraction issues
- [ ] Update GraphQL schema if needed for new mutations

### Frontend

- [ ] Create unified upload page at `/upload`
- [ ] Build document drop zone component (reuse existing if possible)
- [ ] Build detection results UI
  - Content type checkboxes with descriptions
  - Testimonial author display when detected
  - Low confidence warning (only shown when confidence is low)
- [ ] Build "Not what you expected?" correction UI
  - Radio options: "just a resume", "just a reference letter"
  - Free-text option for other issues
- [ ] Build manual selection fallback UI (for low confidence detection)
  - "What does this document contain?" with radio options
- [ ] Build error state for unreadable/corrupted documents
  - Message asking user to try a different/clearer document
- [ ] Build review/preview UI showing extracted data
  - Career info: positions, skills, education
  - Testimonial: quote, author, verified skills
- [ ] Build "Report extraction issue" feedback UI
- [ ] Handle loading states for detection and extraction
- [ ] Redirect to profile after successful import

### Integration

- [ ] Ensure extracted career info merges correctly with existing profile data
- [ ] Ensure testimonials link to correct profile
- [ ] Handle edge cases: empty detection, low confidence, extraction failures

### Testing

- [ ] Unit tests for detection prompt/logic
- [ ] Integration tests for upload flow
- [ ] Test with various document types:
  - Pure resume
  - Pure reference letter
  - Hybrid document (both)
  - Edge cases (letters of recommendation without career details)

## Definition of Done

- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Visual verification with agent-browser (for UI changes)
- [ ] All checklist items above are completed
- [ ] Branch pushed and PR created for human review

## Resolved Questions

- **URL**: `/upload`
- **Confidence display**: Only show when low (as a warning)
- **Low confidence handling**: Force manual selection
- **Unreadable documents**: Ask user to try a different document

## Future Considerations (out of scope)

- Deprecating separate resume/reference letter upload flows
- Re-processing previously uploaded documents
- Batch upload of multiple documents