---
# credfolio2-gg0o
title: Auto-apply skill/experience validations during unified import
status: in-progress
type: bug
created_at: 2026-02-06T09:47:46Z
updated_at: 2026-02-06T09:47:46Z
---

## Problem
When importing both a resume and reference letter together, the import creates profile skills (from resume) and testimonials (from reference letter) but does NOT create skill_validations or experience_validations to cross-reference them.

## Fix
After both resume and reference letter are materialized in importDocumentResults, automatically cross-reference skill mentions and experience mentions from the reference letter with the materialized profile skills and experiences.

## Checklist
- [ ] Add CrossReferenceValidations method to MaterializationService
- [ ] Call it from ImportDocumentResults when both documents are present
- [ ] Add tests for the cross-referencing logic
- [ ] pnpm lint passes
- [ ] pnpm test passes