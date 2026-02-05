---
# credfolio2-1i3d
title: Debug Braintrust instrumentation - no data appearing
status: completed
type: bug
priority: normal
created_at: 2026-02-05T12:50:14Z
updated_at: 2026-02-05T14:06:09Z
---

Braintrust integration was added but no data appears in Braintrust dashboard. API key is confirmed present.

## Root Cause Analysis

**Primary cause: Firewall blocking outbound OTLP traffic.**
The devcontainer has a strict whitelist-based firewall. `api.braintrust.dev` and `www.braintrust.dev` were not in the whitelist, so the OTLP HTTP exporter's requests were silently rejected by iptables.

**Secondary cause: Zero observability on export failures.**
The OpenTelemetry batch span processor swallows export errors internally. No OTel error handler was configured, and the Braintrust SDK had no logger — so the silent firewall rejections produced no log output whatsoever.

## Checklist
- [x] Add `api.braintrust.dev` and `www.braintrust.dev` to devcontainer firewall whitelist
- [x] Add OpenTelemetry error handler to surface OTLP export failures in server logs
- [x] Pass application logger to Braintrust SDK for login/tracing visibility
- [x] Run `go mod tidy` to fix `braintrust-sdk-go` being marked as indirect
- [x] Update tests for new adapter types
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` — LLM tests pass (database tests fail due to missing postgres, pre-existing)

## Definition of Done
- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures (LLM package; DB tests are pre-existing failures)
- [x] All other checklist items above are completed
- [x] Branch pushed and PR created for human review