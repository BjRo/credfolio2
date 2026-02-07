---
# credfolio2-0gh5
title: Deep Link from Skill Validation Popover
status: in-progress
type: feature
priority: normal
created_at: 2026-02-07T09:29:33Z
updated_at: 2026-02-07T13:18:01Z
parent: credfolio2-klgo
---

Add a per-validation "View in source →" link to the skill/experience validation popover that opens the PDF viewer with the validating quote highlighted.

## Context

Currently, the ValidationPopover shows validation details (author, quote snippet) with a single "View full testimonials →" link at the bottom that scrolls to the testimonials section. This isn't very useful because:
1. It doesn't link to the specific source document
2. It groups all validations under one link, rather than linking each individually

The testimonials section (credfolio2-kahb) already deep-links each testimonial to the PDF viewer with quote highlighting. We need to bring the same pattern to the validation popover.

## Design Decision

- **Per-validation text link**: Each validation entry gets a small "View in source →" text link below its quote snippet
- **Remove bottom link**: The "View full testimonials →" link at the bottom is removed entirely (redundant once per-validation links exist)
- **Consistent with testimonials**: Uses the same `buildViewerUrl()` utility and link styling as TestimonialsSection

## Checklist

### Backend / GraphQL Query Changes
- [x] Update `GetSkillValidations` query in `queries.graphql` to include `referenceLetter { id file { id url } }` (file sub-selection is missing today)
- [x] Update `GetExperienceValidations` query similarly to include `referenceLetter.file { id url }`
- [x] Run GraphQL codegen (`pnpm --filter frontend codegen`) to regenerate types

### Frontend: ValidationPopover Changes
- [x] Import `buildViewerUrl` from `@/lib/viewer` in `ValidationPopover.tsx`
- [x] Add per-validation "View in source →" link after each validation's quote snippet:
  - Compute URL: `buildViewerUrl(referenceLetter.id, validation.quoteSnippet)`
  - Only show when `referenceLetter?.id` and `referenceLetter?.file?.url` both exist
  - Style: small text link (`text-xs text-primary hover:underline`), opens in new tab
  - Place it after the blockquote, within the validation entry div
- [x] Remove the bottom "View full testimonials →" anchor link entirely
- [x] Handle edge case: if a validation has no quote snippet, link to the viewer without highlight

### Tests
- [x] Add/update tests in `ValidationPopover.test.tsx`:
  - Test that "View in source →" link appears for validations with reference letter + file
  - Test that link is absent when no file exists
  - Test that the old "View full testimonials →" link is gone
  - Test correct URL construction with `buildViewerUrl`

## Technical Notes

- The `referenceLetter` type in GraphQL already exposes `file { id url }` — we just aren't querying it in `GetSkillValidations`/`GetExperienceValidations`
- `buildViewerUrl(letterId, highlightText)` already handles URL encoding and truncation (max 500 chars)
- The `quoteSnippet` field is the short validation quote — shorter than full testimonial quotes, but should work fine with the PDF text search
- Pattern reference: `TestimonialsSection.tsx` lines 328-334 show how `getSourceUrl` uses `buildViewerUrl`

## Definition of Done
- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [ ] Visual verification with agent-browser (for UI changes)
- [x] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review
- [ ] Automated code review passed (`@review-backend` and/or `@review-frontend`)