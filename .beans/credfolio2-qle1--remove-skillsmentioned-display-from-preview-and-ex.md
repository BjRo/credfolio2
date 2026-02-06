---
# credfolio2-qle1
title: Remove skillsMentioned display from preview and extraction review UIs
status: in-progress
type: task
created_at: 2026-02-06T22:12:05Z
updated_at: 2026-02-06T22:12:05Z
---

The `skillsMentioned` field displayed in testimonial preview and extraction review screens is LLM extraction metadata used internally by the backend to filter `validatedSkills`. Showing it to users is misleading — especially labeled "Skills validated:" in the preview when nothing has been validated yet.

## Changes
- Remove skills badges from reference letter preview (`preview/TestimonialsSection.tsx`)
- Remove skills badges from extraction review (`ExtractionReview.tsx`)
- Update related tests

The profile page's `validatedSkills` display (which uses real ProfileSkill objects) remains unchanged — that's the correct, functional display.

## Definition of Done
- [ ] Remove skillsMentioned display from preview/TestimonialsSection.tsx
- [ ] Remove skillsMentioned display from ExtractionReview.tsx
- [ ] Update tests for both components
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Branch pushed and PR created for human review
- [ ] Automated code review passed (`@review-frontend`)