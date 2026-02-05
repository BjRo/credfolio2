---
# credfolio2-f2na
title: Wrap all registered LLM providers with resilience (timeout, retries, circuit breaker)
status: in-progress
type: bug
priority: high
created_at: 2026-02-05T16:48:07Z
updated_at: 2026-02-05T16:48:07Z
---

Chain-resolved providers (used for resume and reference extraction) bypass the ResilientProvider wrapper. Only the default/fallback provider gets resilience. This causes context deadline exceeded errors when chain-resolved providers (e.g. OpenAI for reference extraction) take longer than Go's default HTTP timeout.

Fix: wrap each provider with ResilientProvider at registration time so all chain-resolved providers get timeout (120s), retries, and circuit breaker protection.

## Checklist
- [x] Wrap providers with ResilientProvider at registration time in createProviderRegistry
- [x] Remove separate resilient wrapping of defaultProvider in createLLMExtractor
- [x] Tests pass
- [x] Lint passes

## Definition of Done
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review