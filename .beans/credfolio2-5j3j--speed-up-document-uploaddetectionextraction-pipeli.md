---
# credfolio2-5j3j
title: Speed up document upload/detection/extraction pipeline
status: completed
type: task
priority: high
created_at: 2026-02-05T23:11:46Z
updated_at: 2026-02-06T12:58:04Z
parent: credfolio2-3ram
blocking:
    - credfolio2-1g6t
---

## Summary

The end-to-end pipeline (upload → detection → extraction → import) is very slow. A single document with both career info and testimonial takes 2-3+ minutes due to redundant work and sequential processing.

## Identified Bottlenecks

### 1. Duplicate text extraction (biggest win)
Text is extracted via LLM vision API **twice**: once in `DocumentDetectionWorker` (document_detection.go:96) and again in `DocumentProcessingWorker` (document_processing.go:146). Each call costs 30-60s + LLM tokens. The detection worker discards the extracted text after classification.

**Fix**: Store extracted text in a new `extracted_text` column on the `files` table during detection. Processing worker reads it from DB instead of re-extracting.

### 2. Duplicate file download from storage
Both workers independently download the file from MinIO. If text is stored in DB (fix #1), the processing worker does not need to download the file at all.

### 3. Sequential resume + letter extractions (not viable to parallelize)
Resume extraction produces skills that are fed as context to the letter extraction LLM call. This data dependency prevents parallelization without quality loss. Keeping sequential.

### 4. Detection does not chain into processing (skipped per user decision)
The current frontend flow (detect → review → extract) provides user value. Skipping auto-chaining.

## Checklist

- [x] Migration: add `extracted_text` TEXT column to `files` table
- [x] Detection worker: store extracted text in file record after extraction
- [x] Processing worker: read stored text from file instead of re-extracting
- [x] ~~Processing worker: parallelize resume + letter extractions~~ (not viable — skill dependency)
- [x] ~~Consider auto-chaining detection → processing~~ (skipped per user decision)

### Definition of Done
- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] All other checklist items above are completed
- [x] Branch pushed and PR created for human review
- [x] Automated code review passed (`@review-backend`)
