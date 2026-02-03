---
# credfolio2-c9rj
title: Duplicate Document Detection
status: completed
type: feature
priority: normal
created_at: 2026-01-29T13:43:34Z
updated_at: 2026-02-03T11:49:42Z
parent: credfolio2-2ex3
---

Detect when a document (resume or reference letter) has already been uploaded by the same user to prevent accidental re-uploads and allow intentional re-imports.

## Context

Currently, users can upload the same document multiple times without any warning. This leads to:
- Duplicate processing costs (LLM extraction is expensive)
- Confusion about which version of the profile data is current
- Wasted storage space

## Solution

Calculate and store a content hash (SHA-256) for each uploaded file, then check for duplicates before processing.

## Checklist

### Backend: Schema & Storage
- [x] Add `content_hash` column (VARCHAR(64)) to `files` table via migration
- [x] Add unique index on `(user_id, content_hash)` in `files` table
- [x] Update file upload handler to calculate SHA-256 hash during upload

### Backend: Duplicate Detection Logic
- [x] Create a `CheckDuplicateFile` query/resolver that checks if hash exists for user
- [x] Return existing file/resume/reference_letter info if duplicate found
- [x] Add GraphQL mutation option `forceReimport: Boolean` to bypass duplicate check

### Frontend: User Confirmation Flow
- [x] Before upload completes, check for duplicate via GraphQL query
- [x] If duplicate found, show confirmation dialog with options:
  - "This document was uploaded on [date]. Re-import anyway?"
  - Show previous extraction status (completed/failed/pending)
- [x] If user confirms, call mutation with `forceReimport: true`
- [x] If user cancels, redirect to existing profile/document view

### Testing
- [x] Unit tests for hash calculation
- [x] Integration test: CheckDuplicateFile query tests (duplicate found, not found, edge cases)
- [x] Integration test: file repository GetByUserIDAndContentHash tests
- [x] E2E test: verify confirmation dialog appears and works correctly

## Definition of Done
- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] Visual verification with agent-browser (for UI changes)
- [x] All other checklist items above are completed
- [x] Branch pushed and PR created for human review

## Technical Notes

- Use SHA-256 for content hashing (crypto-secure, widely supported)
- Hash should be calculated server-side to ensure consistency
- The hash is on file content only, not metadata (so renamed files are still detected)
- Consider: should we also detect near-duplicates? (out of scope for v1)