---
# credfolio2-z8v3
title: Instrument LLM calls with Braintrust
status: todo
type: feature
priority: normal
created_at: 2026-02-03T10:33:23Z
updated_at: 2026-02-03T10:54:10Z
parent: credfolio2-2ex3
---

## Summary
Add Braintrust SDK instrumentation to all LLM calls in the application for observability and debugging.

## Context
Braintrust provides LLM observability, allowing inspection of prompts, responses, latency, and costs. This will help debug and improve the extraction quality.

## Requirements
- Integrate https://github.com/braintrustdata/braintrust-sdk-go
- Instrument all existing LLM calls (resume extraction, reference letter extraction, etc.)
- Ensure traces are visible in Braintrust dashboard

## Technical Decisions
- **API Key**: Will be provided via `BRAINTRUST_API_KEY` environment variable (needs to be configured)
- Instrumentation should be transparent - wrap existing LLM client without changing call sites

## Checklist
- [ ] Configure `BRAINTRUST_API_KEY` in environment (docker-compose, devcontainer, etc.)
- [ ] Add braintrust-sdk-go dependency to go.mod
- [ ] Review braintrust-sdk-go documentation for integration pattern
- [ ] Create Braintrust client wrapper in backend
- [ ] Identify all LLM call sites in codebase
- [ ] Wrap LLM calls with Braintrust tracing
- [ ] Add project/experiment naming for organized traces
- [ ] Verify traces appear in Braintrust dashboard
- [ ] Document how to view traces in Braintrust

## Definition of Done
- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] All checklist items above are completed
- [ ] Branch pushed and PR created for human review
