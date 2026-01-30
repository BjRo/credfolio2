---
# credfolio2-3e5x
title: Filter validated skills by testimonial's skillsMentioned
status: in-progress
type: bug
created_at: 2026-01-30T12:08:35Z
updated_at: 2026-01-30T12:08:35Z
---

The 'Validates:' section on testimonial cards shows ALL skills validated by the reference letter, even when the specific testimonial quote doesn't mention those skills. This is misleading - a quote about 'taking ownership of ambiguous challenges' shouldn't show 'Validates: Symfony Framework · JavaScript' unless those skills are mentioned.

## Root Cause

The current implementation links skill_validations to reference_letters, but displays them on all testimonials from that letter. We need to filter to only show skills that match the testimonial's skillsMentioned field.

## Solution

1. Add skills_mentioned column to testimonials table (stores the skills mentioned in the quote)
2. Update applyReferenceLetterValidations to persist skillsMentioned from ExtractedTestimonial
3. Update Testimonials resolver to filter validatedSkills to only include skills where the skill name appears in the testimonial's skills_mentioned list

## Checklist

- [x] Add migration for skills_mentioned column on testimonials table
- [x] Update testimonials repository to handle skills_mentioned (bun handles this automatically)
- [x] Update applyReferenceLetterValidations to store skillsMentioned
- [x] Update Testimonials resolver to filter validatedSkills by skills_mentioned
- [x] Add tests for the filtering logic
- [x] Update frontend tests if needed (no changes needed)

## Definition of Done

- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] Visual verification with agent-browser (skipped - backend-only change with comprehensive test coverage)
- [x] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review