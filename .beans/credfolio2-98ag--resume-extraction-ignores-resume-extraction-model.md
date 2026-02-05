---
# credfolio2-98ag
title: Resume extraction ignores RESUME_EXTRACTION_MODEL — all calls route to Anthropic
status: in-progress
type: bug
priority: high
created_at: 2026-02-05T14:35:04Z
updated_at: 2026-02-05T15:32:05Z
parent: credfolio2-2ex3
---

## Observed Behavior

`RESUME_EXTRACTION_MODEL=openai/gpt-4o` is configured, but Braintrust monitoring shows **all** LLM API calls going to Anthropic exclusively. The resume structured data extraction should be using OpenAI.

## Investigation

The configuration flows correctly through parsing and chain creation:

1. `config.go:191` reads `RESUME_EXTRACTION_MODEL` ✓
2. `config.go:39-56` parses `"openai/gpt-4o"` → `("openai", "gpt-4o")` ✓
3. `main.go:300-301` creates `resumeChain = ProviderChain{{Provider: "openai", Model: "gpt-4o"}}` ✓
4. `main.go:310-314` passes the chain to the extractor ✓

The bug is in **`extraction.go:118-127`** — `getProviderForChain()` silently swallows errors:

```go
func (e *DocumentExtractor) getProviderForChain(chain ProviderChain) domain.LLMProvider {
    if len(chain) > 0 && e.config.ProviderRegistry != nil {
        chained, err := NewChainedProvider(e.config.ProviderRegistry, chain)
        if err == nil {
            return chained
        }
        // Fall back to default on error  ← SILENT! No logging!
    }
    return e.defaultProvider  // ← Always Anthropic
}
```

If `NewChainedProvider` fails (e.g., provider "openai" not registered because `OPENAI_API_KEY` is missing, or any validation error), it **silently falls back to the default Anthropic provider** with zero logging. This makes the bug invisible at runtime.

### Likely root cause candidates

1. **Silent error swallowing**: `getProviderForChain` doesn't log when chain creation fails — makes debugging impossible
2. **Missing OpenAI registration**: OpenAI provider is only registered if `OPENAI_API_KEY` is set (`main.go:241`). If the key is missing/empty, "openai" won't be in the registry, chain creation fails silently
3. **No startup validation**: The system doesn't verify at startup that configured chains resolve to registered providers — it only fails at request time, silently

### Secondary issue

`ExtractResumeData` (line 376) still uses the deprecated `DefaultModel` field (always empty) instead of letting the chain handle model selection. While `ChainedProvider.Complete()` at `provider_chain.go:126-128` does override empty models with the chain's model, this is fragile and confusing.

## Checklist

- [x] Add error logging in `getProviderForChain` when chain creation fails (log the error + chain config)
- [x] Add startup validation in `createLLMExtractor` (main.go) that warns if a configured chain references an unregistered provider
- [x] Verify the OpenAI provider is actually being registered (check OPENAI_API_KEY is set in the environment)
- [x] Remove usage of deprecated `DefaultModel` in `ExtractResumeData` — let the chain handle model selection exclusively
- [x] Add an integration/unit test that verifies provider chain routing works correctly
- [ ] Verify the fix in Braintrust — confirm OpenAI calls appear for resume extraction

## Definition of Done
- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [ ] All other checklist items above are completed
- [x] Branch pushed and PR created for human review

## Key Files

- `src/backend/internal/infrastructure/llm/extraction.go:118-127` — silent fallback (primary fix)
- `src/backend/internal/infrastructure/llm/extraction.go:376` — deprecated DefaultModel usage
- `src/backend/internal/infrastructure/llm/provider_chain.go:92-109` — chain creation that can fail
- `src/backend/cmd/server/main.go:241-254` — OpenAI provider registration (conditional on API key)
- `src/backend/cmd/server/main.go:278-318` — extractor setup with chains