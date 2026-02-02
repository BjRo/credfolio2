---
# credfolio2-9qfq
title: Replace source icon with kebab menu in testimonials
status: todo
type: feature
created_at: 2026-02-02T14:21:54Z
updated_at: 2026-02-02T14:21:54Z
parent: credfolio2-2ex3
---

## Overview

Replace the current "source" icon on testimonial cards with a kebab menu ("...") that provides a consistent action pattern across the profile page, matching the existing implementation in WorkExperience and Education sections.

## User Story

As a user viewing testimonials, I want a consistent interaction pattern across all profile sections, so the UI feels cohesive and predictable.

## Current State

- Testimonials show a source icon to view the reference letter
- WorkExperience and Education use a kebab menu for actions

## Proposed Change

Replace the source icon with a kebab menu that contains:
- "View source document" - links to the reference letter PDF
- (Future) Additional actions as needed (edit, delete, etc.)

## Implementation

- [ ] Review kebab menu implementation in WorkExperience/Education components
- [ ] Extract shared menu component if not already abstracted
- [ ] Replace source icon with kebab menu in TestimonialsSection
- [ ] Ensure consistent styling and positioning
- [ ] Add "View source document" action that opens the reference letter

## Definition of Done
- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Visual verification with agent-browser (for UI changes)
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review