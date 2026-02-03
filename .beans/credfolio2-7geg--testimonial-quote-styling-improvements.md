---
# credfolio2-7geg
title: Testimonial quote styling improvements
status: in-progress
type: task
priority: normal
created_at: 2026-02-03T10:14:50Z
updated_at: 2026-02-03T10:15:38Z
parent: 1kt0
---

Visual refinements to the testimonial quotes in the "What Others Say" section on the profile page.

## Context

The current implementation has three visual issues:
1. The left border is interrupted between quotes (each `QuoteItem` has its own `border-l-2`)
2. Quotes are plain text without bullet points
3. There's excessive whitespace between the opening quote mark (`"`) and the first word

## Changes

### 1. Continuous left border
Move the `border-l-2` from individual `QuoteItem` divs to a wrapper around all quotes, so the line spans continuously from the first quote to the last within an author card.

### 2. Triangle bullet points
Add a triangle icon (`▸` or `ChevronRight`) as a bullet point before each quote to create a proper bulleted list appearance. Position it inline with the quote text, between the left border and the quote mark.

### 3. Tighten opening quote spacing
Reduce the gap between the decorative opening quote (`"`) and the first word. Currently there's `pl-4` padding creating unnecessary whitespace. Adjust positioning so the quote mark sits closer to the text.

## Files to modify

- `src/frontend/src/components/profile/TestimonialsSection.tsx` - main component (profile view only, not preview page)

## Checklist

- [x] Move border-l-2 to wrapper div for continuous line
- [x] Add triangle bullet icon before each quote
- [x] Tighten spacing between opening quote mark and text
- [x] Verify visual appearance matches expectations

## Definition of Done
- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] Visual verification with agent-browser (for UI changes)
- [x] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review