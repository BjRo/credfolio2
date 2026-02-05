---
# credfolio2-wh63
title: Fix HTTP WriteTimeout causing context canceled errors during LLM detection
status: in-progress
type: bug
priority: high
created_at: 2026-02-05T21:59:46Z
updated_at: 2026-02-05T22:00:45Z
---

## Summary

When uploading a document, the detection step calls the Anthropic API via the resilient LLM provider. The HTTP server's WriteTimeout (120s) is shorter than the resilient provider's request timeout (300s), so Go cancels the request context before the LLM response arrives. This surfaces as: `[ERROR] [detection] Failed to extract text from document {"error":"anthropic: context canceled"}`

## Root Cause

In `src/backend/internal/config/config.go`, the server timeouts are:
- `ReadTimeout: 30 * time.Second`
- `WriteTimeout: 120 * time.Second`

In `src/backend/cmd/server/main.go`, the resilient provider config is:
- `RequestTimeout: 300 * time.Second`

The WriteTimeout fires first (120s < 300s), canceling the context.

## Fix

Increase `WriteTimeout` to be longer than `RequestTimeout` (e.g., 360s) so the LLM provider's own timeout governs the deadline, not the HTTP server.

## Checklist

- [x] Increase WriteTimeout in config.go to 360s (longer than provider's 300s)
- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] All checklist items above are completed
- [ ] Branch pushed and PR created for human review