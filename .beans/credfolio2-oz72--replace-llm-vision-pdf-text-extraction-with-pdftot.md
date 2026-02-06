---
# credfolio2-oz72
title: Replace LLM vision PDF text extraction with pdftotext
status: completed
type: task
priority: normal
created_at: 2026-02-06T13:06:11Z
updated_at: 2026-02-06T13:26:37Z
parent: credfolio2-3ram
---

## Summary

The pdf_text_extraction step currently sends the entire PDF as a base64 document block to Claude Sonnet 4.5's vision API, taking ~29s and costing $0.036 per document. Most resumes are **text-based PDFs** (created by Word, Google Docs, LaTeX) — not scanned images — so this is massively overkill.

## Approach

Use a Go PDF text extraction library (e.g. `github.com/ledongthuc/pdf` or `github.com/dslipak/pdf`) to extract text locally in milliseconds. Only fall back to the LLM vision path for scanned/image-based PDFs where the Go library returns empty/garbage text.

This keeps the solution pure Go with no external binary dependencies, simplifying the build and deployment.

**Expected impact:** 29s → <1s for text-based PDFs (majority of uploads). Zero quality risk — better for text PDFs, same for scanned.

## Checklist

- [x] Evaluate and add a Go PDF text extraction library (e.g. `github.com/ledongthuc/pdf`)
- [x] Create a local PDF text extraction function that uses the Go library
- [x] Add a quality heuristic to detect if extracted text is usable (non-empty, contains real words, reasonable length)
- [x] If local extraction succeeds quality check, use it directly and skip the LLM vision call
- [x] If local extraction fails or returns garbage (scanned PDF), fall back to existing LLM vision extraction
- [x] Preserve the existing OCR normalization rules in the LLM fallback path
- [x] Add telemetry/span attributes to track which path was used (local vs LLM)
- [x] Write tests for the local extraction path and the fallback logic

## Definition of Done
- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] All other checklist items above are completed
- [x] Branch pushed and PR created for human review
- [x] Automated code review passed (`@review-backend` and/or `@review-frontend`)