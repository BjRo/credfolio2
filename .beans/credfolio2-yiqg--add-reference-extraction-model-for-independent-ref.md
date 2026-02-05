---
# credfolio2-yiqg
title: Add REFERENCE_EXTRACTION_MODEL for independent reference letter model config
status: completed
type: task
priority: normal
created_at: 2026-02-05T16:27:49Z
updated_at: 2026-02-05T16:31:28Z
---

Add a separate REFERENCE_EXTRACTION_MODEL env var so reference letter extraction can use a different provider/model than resume extraction. Currently both share RESUME_EXTRACTION_MODEL via ResumeExtractionChain.

## Checklist
- [x] Add ReferenceExtractionModel to LLMConfig with ParseReferenceExtractionModel()
- [x] Add REFERENCE_EXTRACTION_MODEL to config loading
- [x] Add ReferenceExtractionChain to DocumentExtractorConfig
- [x] Update ExtractLetterData to use the new chain
- [x] Wire up the new chain in main.go createLLMExtractor
- [x] Update .env.example with REFERENCE_EXTRACTION_MODEL
- [x] Update tests

## Definition of Done
- [x] Tests written (TDD: write tests before implementation)
- [x] pnpm lint passes with no errors
- [x] pnpm test passes with no failures
- [x] All other checklist items above are completed
- [x] Branch pushed and PR created for human review