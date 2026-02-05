---
# credfolio2-hyzh
title: Remove global LLM_PROVIDER in favor of use-case-specific model config
status: in-progress
type: task
priority: normal
created_at: 2026-02-05T13:42:40Z
updated_at: 2026-02-05T15:53:47Z
parent: credfolio2-2ex3
---

## Problem

The `.env.example` still lists `LLM_PROVIDER=anthropic` as a global setting. Instead, the LLM provider should be chosen on a per-use-case basis (e.g., `RESUME_EXTRACTION_MODEL`), not as a single global provider.

## Solution

Rework the configuration so that `LLM_PROVIDER` is no longer needed. Each use case that requires an LLM should specify its own model/provider configuration independently.

## Checklist
- [x] Audit all places that reference `LLM_PROVIDER` (env files, config code, backend logic)
- [x] Replace global provider with use-case-specific config (e.g., `RESUME_EXTRACTION_MODEL`)
- [x] Update `.env.example` to remove `LLM_PROVIDER` and document use-case-specific vars
- [x] Update backend config parsing to reflect the new approach
- [x] Ensure existing functionality (resume extraction) still works with new config shape

## Definition of Done
- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review