---
# credfolio2-1vtw
title: Improve reference letter testimonial extraction
status: draft
type: feature
created_at: 2026-02-03T10:33:18Z
updated_at: 2026-02-03T10:33:18Z
parent: 2ex3
---

## Summary
Enhance the LLM prompt for extracting testimonials from reference letters to produce higher quality, performance-focused testimonials.

## Requirements
- Extract 3-5 testimonials per reference letter (instead of current behavior)
- Focus testimonials on the person's performance and achievements
- Order testimonials with strongest positive statements first (by sentiment strength)
- Ensure extracted quotes are impactful and meaningful

## Technical Decisions
- **Ranking method**: Sentiment strength - most positive/enthusiastic language ranks higher
- **Scope**: New uploads only - existing testimonials remain unchanged
- The LLM should score each testimonial's sentiment during extraction
- Store sentiment score to enable sorting on the frontend

## Changes Needed
- Update the reference letter extraction prompt in the backend to:
  - Request 3-5 testimonials (not fewer)
  - Focus on performance-related quotes
  - Include a sentiment/strength score (1-10) for each
- Add `sentimentScore` field to Testimonial model/schema
- Order testimonials by sentiment score (descending) when returning from API

## Checklist
- [ ] Review current extraction prompt (locate in codebase)
- [ ] Update prompt to request 3-5 performance-focused testimonials
- [ ] Add sentiment scoring instruction to prompt
- [ ] Add `sentimentScore` field to Testimonial in GraphQL schema
- [ ] Add database migration for sentiment_score column
- [ ] Update extraction response parsing to capture score
- [ ] Sort testimonials by sentiment score in resolver
- [ ] Update frontend to display in sorted order (or rely on API order)

## Definition of Done
- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] All checklist items above are completed
- [ ] Branch pushed and PR created for human review
