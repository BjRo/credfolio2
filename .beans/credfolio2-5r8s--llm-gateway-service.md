---
# credfolio2-5r8s
title: LLM Gateway Service
status: completed
type: epic
priority: normal
created_at: 2026-01-20T11:24:26Z
updated_at: 2026-01-23T16:39:37Z
parent: credfolio2-tikg
blocking:
    - credfolio2-tmlf
---

Build an abstraction layer for LLM interactions with fault tolerance patterns.

## Current State

**Status:** Ready for PR review - all implementation complete, pushed to PR #22

**Branch:** `feat/credfolio2-5r8s-llm-gateway-service`

## What Was Built

### Core LLM Infrastructure (`src/backend/internal/infrastructure/llm/`)

1. **AnthropicProvider** (`anthropic.go`) - Claude API integration using official `anthropic-sdk-go`
   - Standard completions and structured outputs (via Beta API)
   - Image/PDF vision support for document extraction
   - Proper error conversion to domain types

2. **ResilientProvider** (`resilient.go`) - Fault tolerance using `failsafe-go`
   - Retry with exponential backoff (configurable attempts, delays)
   - Circuit breaker (configurable failure threshold)
   - Request timeout (default 120s)

3. **DocumentExtractor** (`extraction.go`) - Text extraction from images/PDFs
   - Uses Claude's vision capabilities
   - Configurable prompts and models

4. **LoggingProvider** (`logging.go`) - Request/response logging wrapper

### Test Page (`/extract-test`)

- **Backend:** `POST /api/extract` endpoint in `src/backend/internal/handler/extract.go`
- **Frontend:** `src/frontend/src/app/extract-test/page.tsx`
- Drag-and-drop file upload, extraction results display

## Configuration

Requires `ANTHROPIC_API_KEY` in environment. Without it, extraction returns 503.

## Files Changed

- `src/backend/go.mod` / `go.sum` - Added anthropic-sdk-go, failsafe-go
- `src/backend/internal/domain/llm.go` - Added OutputSchema field
- `src/backend/internal/infrastructure/llm/*.go` - Core implementation
- `src/backend/internal/handler/extract.go` - Extraction endpoint
- `src/backend/cmd/server/main.go` - Wired up LLM provider stack
- `src/frontend/src/app/extract-test/page.tsx` - Test UI

## Checklist

- [x] Design LLM gateway interface (Provider, Message, Response)
- [x] Add request/response logging
- [x] Create document-to-text extraction using Claude vision
- [x] Add provider fallback mechanism (stub for now)
- [x] Add dependencies (anthropic-sdk-go, failsafe-go)
- [x] Add OutputSchema field to LLMRequest for structured outputs
- [x] Rewrite Anthropic provider using anthropic-sdk-go
- [x] Rewrite resilient provider using failsafe-go
- [x] Delete old hand-rolled circuitbreaker.go and retry.go
- [x] Update tests for SDK-based implementation
- [x] Update documentation (doc.go)
- [x] Add E2E test page (/extract-test)

## Bug Fixes (2026-01-23)

1. **PDF media type error** - PDFs were being sent as image blocks, but Anthropic API requires document blocks. Fixed in `anthropic.go` to use `NewDocumentBlock` with `Base64PDFSourceParam` for PDFs.

2. **Server WriteTimeout too short** - Default 15s timeout caused connection drops during 60s+ extractions. Increased to 120s in `config.go`.

3. **CORS issues** - Added Next.js API route proxy (`/api/extract`) to avoid browser CORS problems with long requests.

## Follow-up Work

See **credfolio2-he9y** for remaining improvements (progress indication, browser timeout handling).

## Next Steps (after PR merge)

- Mark this bean as completed
- Mark credfolio2-vwxr (test page bean) as completed
