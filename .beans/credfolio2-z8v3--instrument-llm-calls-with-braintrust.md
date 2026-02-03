---
# credfolio2-z8v3
title: Instrument LLM calls with Braintrust
status: in-progress
type: feature
priority: normal
created_at: 2026-02-03T10:33:23Z
updated_at: 2026-02-03T17:16:14Z
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
- [x] Configure `BRAINTRUST_API_KEY` in environment (docker-compose, devcontainer, etc.) - Added to .env.example
- [x] Add braintrust-sdk-go dependency to go.mod
- [x] Review braintrust-sdk-go documentation for integration pattern
- [x] Create Braintrust initialization module in backend (internal/infrastructure/llm/braintrust.go)
- [x] Identify all LLM call sites in codebase (anthropic.go, openai.go providers)
- [x] Modify LLM providers to accept middleware options (AnthropicConfig.Middleware, OpenAIConfig.Middleware)
- [x] Initialize Braintrust and inject tracing middleware at startup (createProviderRegistry in main.go)
- [x] Add project naming for organized traces (BRAINTRUST_PROJECT env var, defaults to "credfolio")
- [ ] Verify traces appear in Braintrust dashboard (requires BRAINTRUST_API_KEY to test)
- [x] Document how to view traces in Braintrust (see README section below)

## Definition of Done
- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [ ] All checklist items above are completed
- [ ] Branch pushed and PR created for human review

## How to View Traces in Braintrust

1. **Get API Key**: Go to https://www.braintrust.dev/app/settings and create an API key
2. **Configure Environment**: Set `BRAINTRUST_API_KEY=<your-key>` in your `.env` file
3. **Optional - Set Project Name**: Set `BRAINTRUST_PROJECT=<project-name>` (defaults to "credfolio")
4. **Start Server**: Run `pnpm dev` from the backend directory
5. **View Traces**: Go to https://www.braintrust.dev/app and navigate to your project to see LLM traces

When Braintrust tracing is enabled, all LLM calls (resume extraction, reference letter extraction, document text extraction) will automatically appear in the Braintrust dashboard with:
- Full prompt/response content
- Token usage
- Latency metrics
- Model information
