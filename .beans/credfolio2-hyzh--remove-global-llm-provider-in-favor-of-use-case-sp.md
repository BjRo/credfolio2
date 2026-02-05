---
# credfolio2-hyzh
title: Remove global LLM_PROVIDER in favor of use-case-specific model config
status: todo
type: task
priority: normal
created_at: 2026-02-05T13:42:40Z
updated_at: 2026-02-05T14:08:20Z
parent: credfolio2-2ex3
---

## Problem

The `.env.example` still lists `LLM_PROVIDER=anthropic` as a global setting. Instead, the LLM provider should be chosen on a per-use-case basis (e.g., `RESUME_EXTRACTION_MODEL`), not as a single global provider.

## Solution

Rework the configuration so that `LLM_PROVIDER` is no longer needed. Each use case that requires an LLM should specify its own model/provider configuration independently.

## Checklist
- [ ] Audit all places that reference `LLM_PROVIDER` (env files, config code, backend logic)
- [ ] Replace global provider with use-case-specific config (e.g., `RESUME_EXTRACTION_MODEL`)
- [ ] Update `.env.example` to remove `LLM_PROVIDER` and document use-case-specific vars
- [ ] Update backend config parsing to reflect the new approach
- [ ] Ensure existing functionality (resume extraction) still works with new config shape

## Definition of Done
- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review