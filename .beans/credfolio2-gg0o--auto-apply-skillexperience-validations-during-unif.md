---
# credfolio2-gg0o
title: Auto-apply skill/experience validations during unified import
status: completed
type: bug
priority: normal
created_at: 2026-02-06T09:47:46Z
updated_at: 2026-02-06T10:01:53Z
---

## Problem
When importing both a resume and reference letter together, the import creates profile skills (from resume) and testimonials (from reference letter) but does NOT create skill_validations or experience_validations to cross-reference them.

## Fix
After both resume and reference letter are materialized in importDocumentResults, automatically cross-reference skill mentions and experience mentions from the reference letter with the materialized profile skills and experiences.

**Bug fix:** The LLM extraction puts skills in `discoveredSkills` (not `skillMentions`). Updated `CrossReferenceValidations` to match against both arrays, with deduplication.

## Checklist
- [x] Add CrossReferenceValidations method to MaterializationService
- [x] Call it from ImportDocumentResults when both documents are present
- [x] Match both SkillMentions and DiscoveredSkills from reference letter data
- [x] Add tests for the cross-referencing logic (including discovered skills and deduplication)
- [x] pnpm lint passes
- [x] pnpm test passes