---
# credfolio2-kahb
title: Deep Link from Testimonials Section
status: draft
type: feature
priority: normal
created_at: 2026-02-07T09:29:40Z
updated_at: 2026-02-07T09:29:40Z
parent: credfolio2-klgo
---

Modify the existing "View source document" link in the testimonials section to use the new PDF viewer with quote highlighting instead of opening the raw PDF.

## Checklist

- [ ] In TestimonialsSection.tsx, locate the "View source document" dropdown menu item
- [ ] Change the `href` from the raw file URL to the viewer URL:
  - Old: `{testimonial.referenceLetter.file.url}` (raw presigned URL)
  - New: `/viewer?letterId={testimonial.referenceLetter.id}&highlight={encodeURIComponent(testimonial.quote)}`
- [ ] Ensure `referenceLetter.id` is available in the GraphQL query (it should be — verify)
- [ ] Handle long quotes: if `testimonial.quote` exceeds ~500 chars, truncate for the URL param (the viewer's substring search will still match)
- [ ] Keep `target="_blank" rel="noopener noreferrer"` behavior
- [ ] Update existing tests in TestimonialsSection.test.tsx to reflect the new URL format
- [ ] Verify the link works end-to-end with the viewer page

## Technical Notes

- The current implementation at TestimonialsSection.tsx:238-242 uses the raw file URL
- The `GetTestimonials` query already includes `referenceLetter { id, file { id, url } }`
- The testimonial quote is available as `testimonial.quote` in the component
- Consider adding a utility function `buildViewerUrl(letterId: string, highlightText: string): string` shared between SkillsSection and TestimonialsSection

## Definition of Done
- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Visual verification with agent-browser (for UI changes)
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review
- [ ] Automated code review passed (`@review-backend` and/or `@review-frontend`)