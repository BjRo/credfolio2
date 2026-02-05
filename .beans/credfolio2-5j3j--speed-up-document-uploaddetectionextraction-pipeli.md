---
# credfolio2-5j3j
title: Speed up document upload/detection/extraction pipeline
status: todo
type: feature
priority: high
created_at: 2026-02-05T23:11:46Z
updated_at: 2026-02-05T23:11:46Z
---

## Summary

The end-to-end pipeline (upload → detection → extraction → import) is very slow. A single document with both career info and testimonial takes 2-3+ minutes due to redundant work and sequential processing.

## Identified Bottlenecks

### 1. Duplicate text extraction (biggest win)
Text is extracted via LLM vision API **twice**: once in `DocumentDetectionWorker` (document_detection.go:96) and again in `DocumentProcessingWorker` (document_processing.go:146). Each call costs 30-60s + LLM tokens. The detection worker discards the extracted text after classification.

**Fix**: Store extracted text in a new `extracted_text` column on the `files` table during detection. Processing worker reads it from DB instead of re-extracting.

### 2. Duplicate file download from storage
Both workers independently download the file from MinIO. If text is stored in DB (fix #1), the processing worker doesn't need to download the file at all.

### 3. Sequential resume + letter extractions
When both `extractCareerInfo` and `extractTestimonial` are true, resume and reference letter LLM extractions run **sequentially** in the same job (document_processing.go:156-166). They have no dependencies on each other — both use the same extracted text.

**Fix**: Run both extractions in parallel using goroutines (with proper error collection).

### 4. Detection doesn't chain into processing
Detection and processing are independent workflows requiring separate user actions. Detection already knows the document type — it could automatically trigger the processing job.

**Fix**: Consider an optional "auto-process" flag or have the frontend automatically start processing after detection completes, removing user wait time between steps.

### 5. Detection steps are sequential
Within detection, `ExtractText()` (30-60s) must complete before `DetectDocumentContent()` (5-10s) can start. These are inherently sequential (detection needs the text), but the text extraction is the dominant cost.

## Priority Order

1. **Store extracted text from detection** — eliminates the most expensive redundant operation (~30-60s saved)
2. **Parallelize resume + letter extractions** — cuts extraction time nearly in half for hybrid documents
3. **Auto-chain detection → processing** — removes user wait time between steps
4. **Skip file re-download in processing** — minor I/O savings if text is in DB

## Checklist

- [ ] Migration: add `extracted_text` TEXT column to `files` table
- [ ] Detection worker: store extracted text in file record after extraction
- [ ] Processing worker: read stored text from file instead of re-extracting
- [ ] Processing worker: parallelize resume + letter extractions with goroutines
- [ ] Consider auto-chaining detection → processing (may need frontend changes)

### Definition of Done
- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review