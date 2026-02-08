---
# credfolio2-5ljk
title: Fix UTF-8 truncation and logging bugs in LLM validation
status: in-progress
type: bug
priority: normal
created_at: 2026-02-08T11:45:49Z
updated_at: 2026-02-08T11:46:24Z
---

Fix critical issues found in code review of PR #137 (credfolio2-offq).

## Critical Issues to Fix

1. **UTF-8 Truncation Bug** (validation.go:167, extraction.go:450)
   - Problem: Byte-level string slicing `text[:maxLen]` corrupts multi-byte UTF-8 characters
   - Impact: International names, emoji, and Unicode content can be corrupted (e.g., "José" → "Jos�")
   - Fix: Use proper UTF-8 rune-aware truncation that preserves character boundaries

2. **Truncation Logging Bug** (extraction.go:443, 448)
   - Problem: Logs `len(text)` AFTER truncation, not before
   - Impact: Observability broken - can't tell when documents are actually being truncated
   - Fix: Capture original size before truncation and log that value

## Implementation Plan

### Step 1: Fix UTF-8 truncation helper (TDD)
- Write test in `validation_test.go` for UTF-8 truncation:
  - Test with multi-byte characters (emoji, Chinese, accented Latin)
  - Verify no corruption at boundary
  - Verify length is respected
- Implement UTF-8-safe truncation function in `validation.go`:
  - Use `[]rune` conversion or `utf8.DecodeRuneInString`
  - Truncate at character boundary, not byte boundary
  - Return string <= maxLen bytes that doesn't split characters

### Step 2: Update validation.go truncation (line 167)
- Replace `s[:maxLen]` with UTF-8-safe truncation
- Run tests to verify fix

### Step 3: Fix logging in extraction.go (lines 443, 448)
- Capture `originalSize := len(text)` BEFORE truncation
- Log `originalSize` instead of `len(text)` after truncation
- Update both resume and letter extraction paths

### Step 4: Add integration tests
- Add test in `extraction_test.go` for document with multi-byte characters at truncation boundary
- Verify no corruption and correct logging

## Acceptance Criteria
- [ ] UTF-8-safe truncation function implemented and tested
- [ ] All string truncation uses UTF-8-safe function
- [ ] Logging captures original size before truncation
- [ ] Tests verify multi-byte character handling
- [ ] All existing tests still pass

## Definition of Done
- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Visual verification via `@qa` subagent (via Task tool, for UI changes)
- [ ] ADR written via `/decision` skill (if new dependencies, patterns, or architectural changes were introduced)
- [ ] All other checklist items above are completed
- [ ] Branch pushed to remote
- [ ] PR created for human review
- [ ] Automated code review passed via `@review-backend`, `@review-frontend`, and/or `@review-ai` (for LLM changes) subagents (via Task tool)