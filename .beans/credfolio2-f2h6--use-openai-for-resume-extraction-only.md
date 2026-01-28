---
# credfolio2-f2h6
title: Use OpenAI for resume extraction only
status: in-progress
type: task
created_at: 2026-01-28T16:32:05Z
updated_at: 2026-01-28T16:32:05Z
---

Configure the provider chains so that:
- Document text extraction uses Anthropic (works well for vision/PDF)
- Resume data extraction uses OpenAI (better structured output support)

## Changes
- [x] Add RESUME_EXTRACTION_PROVIDER config option
- [x] Update main.go to use separate provider for resume extraction chain
- [ ] Test with OpenAI API key configured

## Definition of Done
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [ ] Branch pushed and PR created for human review