---
# credfolio2-1vtw
title: Improve reference letter testimonial extraction
status: todo
type: task
priority: normal
created_at: 2026-02-03T10:33:18Z
updated_at: 2026-02-05T14:08:13Z
parent: credfolio2-2ex3
---

## Summary
Enhance the LLM prompt for extracting testimonials from reference letters to produce higher quality, performance-focused testimonials.

## Requirements
- Extract 3-5 testimonials per reference letter (currently 2-4)
- Focus testimonials on the person's performance and achievements
- Order testimonials with strongest positive statements first
- Ensure extracted quotes are impactful and meaningful

## Technical Decisions
- **Ranking method**: LLM returns testimonials pre-sorted by sentiment strength (no DB field needed)
- **Scope**: New uploads only - existing testimonials remain unchanged
- **No schema changes**: The LLM simply returns them in order; frontend displays in API order

## Changes Needed
- Update the reference letter extraction prompt (`src/backend/internal/infrastructure/llm/prompts/reference_letter_extraction_system.txt`) to:
  - Request 3-5 testimonials (change from "2-4")
  - Focus on performance-related quotes
  - Return testimonials ordered by sentiment strength (strongest first)

## Checklist
- [ ] Update prompt to request 3-5 testimonials (change "2-4" to "3-5")
- [ ] Add instruction to focus on performance/achievement quotes
- [ ] Add instruction to order by sentiment strength (strongest first)
- [ ] Verify extraction still works correctly with updated prompt
- [ ] Test with a sample reference letter upload

## Definition of Done
- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] All checklist items above are completed
- [ ] Branch pushed and PR created for human review
