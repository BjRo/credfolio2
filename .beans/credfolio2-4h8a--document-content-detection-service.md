---
# credfolio2-4h8a
title: Document content detection service
status: in-progress
type: task
priority: high
created_at: 2026-02-05T18:01:39Z
updated_at: 2026-02-05T19:58:19Z
parent: credfolio2-3ram
blocking:
    - credfolio2-u1eh
    - credfolio2-nl46
    - credfolio2-d5jd
---

## Summary

Create a lightweight LLM-based document content detection service that quickly classifies what a document contains (career information, testimonials, or both) without running full extraction. This is the foundation for the unified upload flow.

## Background

The existing extractors (resume extraction, reference letter extraction) are heavyweight operations using structured output with JSON schemas. This task creates a much cheaper/faster classification step that runs first, so the user can confirm what to extract before committing to full extraction.

## Checklist

### LLM Prompt
- [x] Create detection system prompt (`document_detection_system.txt`) in `src/backend/internal/infrastructure/llm/prompts/`
  - Input: raw document text (from existing `ExtractText`)
  - Output: structured classification with fields:
    - `has_career_info` (bool) — resume/CV content detected
    - `has_testimonial` (bool) — reference letter / testimonial content
    - `testimonial_author` (string, optional) — detected author name
    - `confidence` (float, 0-1) — overall detection confidence
    - `summary` (string) — brief description of what was found
    - `document_type_hint` (string) — "resume", "reference_letter", "hybrid", "unknown"
  - Keep prompt concise to minimize token usage and latency

### Backend Service
- [x] Add `DetectDocumentContent(ctx, text string)` method to `DocumentExtractor` interface and implementation
  - Uses the lightweight detection prompt
  - Returns structured detection result
  - Should be significantly faster than full extraction (~2-5s vs ~15-30s)
- [x] Add detection result types to domain model
- [x] Configure LLM model for detection (consider using a smaller/faster model)

### GraphQL API
- [x] Add `DocumentDetectionResult` type to GraphQL schema
- [x] Add `detectDocumentContent(userId: ID!, file: Upload!)` mutation
  - Accepts uploaded file
  - Runs text extraction (reuse existing `ExtractText`)
  - Runs lightweight detection on extracted text
  - Returns detection results synchronously (not async job)
  - Also stores the file and extracted text for later use (avoid re-extracting)
- [x] Handle error cases: unreadable file, empty content, extraction failure

### Testing
- [x] Unit tests for detection prompt/logic with various document types
  - Pure resume
  - Pure reference letter
  - Hybrid document (both career info and testimonial)
  - Unreadable/empty document
- [x] Integration test for GraphQL mutation

## Design Notes

- Detection should run **synchronously** (not via job queue) since it's lightweight and the user needs results immediately to proceed
- The extracted text should be cached/stored so the subsequent full extraction doesn't need to re-extract it
- Consider storing as a new `Document` record or reusing the existing `File` + a detection cache

## Definition of Done
- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] All checklist items above are completed
- [x] Branch pushed and PR created for human review