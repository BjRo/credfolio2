---
# credfolio2-2s9m
title: Extract CheckboxCard to shared @/components/ui/ component
status: in-progress
type: task
priority: normal
created_at: 2026-02-06T14:37:19Z
updated_at: 2026-02-06T20:32:15Z
---

## Summary

The CheckboxCard pattern (outer div with `role="checkbox"`, keyboard handling, inner Checkbox, conditional styling) is duplicated across 4 files:

- `ExtractionReview.tsx` (lines 385-426) — already extracted as an internal component
- `CorroborationsSection.tsx` (lines 62-97 and 120-158) — inline duplication
- `TestimonialsSection.tsx` (lines 49-98) — inline duplication
- `DiscoveredSkillsSection.tsx` (lines 49-88) — inline duplication

The `ExtractionReview.tsx` version is already a clean abstraction with `CheckboxCardProps` interface. Move it to `@/components/ui/checkbox-card.tsx` and update all consumers.

## Checklist

- [x] Create `@/components/ui/checkbox-card.tsx` based on ExtractionReview's CheckboxCard
- [x] Update ExtractionReview.tsx to import from shared location
- [x] Update CorroborationsSection.tsx to use shared CheckboxCard
- [x] Update TestimonialsSection.tsx to use shared CheckboxCard
- [x] Update DiscoveredSkillsSection.tsx to use shared CheckboxCard

## Definition of Done
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review
- [ ] Automated code review passed (`@review-backend` and/or `@review-frontend`)
