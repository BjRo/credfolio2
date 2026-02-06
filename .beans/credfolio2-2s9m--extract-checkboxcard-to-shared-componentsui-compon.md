---
# credfolio2-2s9m
title: Extract CheckboxCard to shared @/components/ui/ component
status: in-progress
type: task
created_at: 2026-02-06T14:37:19Z
updated_at: 2026-02-06T14:37:19Z
---

## Summary

The CheckboxCard pattern (outer div with `role="checkbox"`, keyboard handling, inner Checkbox, conditional styling) is duplicated across 4 files:

- `ExtractionReview.tsx` (lines 385-426) — already extracted as an internal component
- `CorroborationsSection.tsx` (lines 62-97 and 120-158) — inline duplication
- `TestimonialsSection.tsx` (lines 49-98) — inline duplication
- `DiscoveredSkillsSection.tsx` (lines 49-88) — inline duplication

The `ExtractionReview.tsx` version is already a clean abstraction with `CheckboxCardProps` interface. Move it to `@/components/ui/checkbox-card.tsx` and update all consumers.

## Checklist

- [ ] Create `@/components/ui/checkbox-card.tsx` based on ExtractionReview's CheckboxCard
- [ ] Update ExtractionReview.tsx to import from shared location
- [ ] Update CorroborationsSection.tsx to use shared CheckboxCard
- [ ] Update TestimonialsSection.tsx to use shared CheckboxCard
- [ ] Update DiscoveredSkillsSection.tsx to use shared CheckboxCard

## Definition of Done
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR updated