---
# credfolio2-5r8s
title: LLM Gateway Service
status: in-progress
type: epic
priority: normal
created_at: 2026-01-20T11:24:26Z
updated_at: 2026-01-23T14:56:28Z
parent: credfolio2-tikg
blocking:
    - credfolio2-tmlf
---

Build an abstraction layer for LLM interactions with fault tolerance patterns.

## Goals
- Abstract away direct LLM SDK usage
- Implement circuit breaker for reliability
- Support multiple LLM providers (Anthropic first)
- Handle document vision/text extraction

## Checklist
- [x] Design LLM gateway interface (Provider, Message, Response)
- [x] Add request/response logging
- [x] Create document-to-text extraction using Claude vision
- [x] Add provider fallback mechanism (stub for now)

### Refactoring (SDK + failsafe-go)
- [x] Add dependencies (anthropic-sdk-go, failsafe-go)
- [x] Add OutputSchema field to LLMRequest for structured outputs
- [x] Rewrite Anthropic provider using anthropic-sdk-go
- [x] Rewrite resilient provider using failsafe-go (retry, circuit breaker, timeout)
- [x] Delete old hand-rolled circuitbreaker.go and retry.go
- [x] Update tests for SDK-based implementation
- [x] Update documentation (doc.go)
