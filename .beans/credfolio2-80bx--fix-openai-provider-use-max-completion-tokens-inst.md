---
# credfolio2-80bx
title: 'Fix OpenAI provider: use max_completion_tokens instead of deprecated max_tokens'
status: completed
type: bug
priority: high
created_at: 2026-02-05T16:36:15Z
updated_at: 2026-02-05T16:39:17Z
---

Newer OpenAI models (gpt-5-nano, o1, o3, etc.) reject `max_tokens` with a 400 error: 'Unsupported parameter: max_tokens is not supported with this model. Use max_completion_tokens instead.' The OpenAI provider must use `MaxCompletionTokens` instead of `MaxTokens` in the SDK call.

## Checklist
- [x] Update OpenAI provider to use MaxCompletionTokens instead of MaxTokens
- [x] Update tests to verify the fix
- [x] pnpm lint and pnpm test pass

## Definition of Done
- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] All other checklist items above are completed
- [x] Branch pushed and PR created for human review