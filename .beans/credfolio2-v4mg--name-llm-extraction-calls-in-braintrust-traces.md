---
# credfolio2-v4mg
title: Name LLM extraction calls in Braintrust traces
status: todo
type: task
priority: normal
created_at: 2026-02-05T14:28:32Z
updated_at: 2026-02-05T14:37:24Z
parent: credfolio2-2ex3
---

## Context

The resume extraction pipeline makes two LLM calls (in `src/backend/internal/infrastructure/llm/extraction.go`):

1. **ExtractText** (line ~129): Extracts raw text from PDF/image using LLM vision
2. **ExtractResumeData** (line ~361): Extracts structured JSON from that raw text

Both calls go through Braintrust's auto-instrumentation middleware (HTTP-level spans named `anthropic.messages.create`), but they're **indistinguishable** in the Braintrust dashboard — you can't tell which call is the PDF text extraction vs. the structured data extraction.

## Goal

Wrap each LLM call in a **named OpenTelemetry span** so they appear as recognizable operations in Braintrust:

- `resume_pdf_extraction` — wraps the ExtractText call
- `resume_structured_data_extraction` — wraps the ExtractResumeData call

The auto-instrumented `anthropic.messages.create` span will automatically become a child of the named parent span.

## Approach

The Braintrust SDK uses standard OpenTelemetry. The pattern (from the SDK examples) is:

```go
tracer := otel.Tracer("credfolio")
ctx, span := tracer.Start(ctx, "resume_pdf_extraction")
defer span.End()
// LLM call here — auto-instrumented span becomes a child
resp, err := provider.Complete(ctx, req)
```

Since `otel.SetTracerProvider(tp)` is already called during Braintrust init (in `braintrust.go:55`), we can use `otel.Tracer()` anywhere in the codebase.

## Checklist

- [ ] Add `go.opentelemetry.io/otel` import to `extraction.go`
- [ ] Wrap `ExtractTextWithRequest` LLM call in a `resume_pdf_extraction` span
- [ ] Wrap `ExtractResumeData` LLM call in a `resume_structured_data_extraction` span
- [ ] Add useful span attributes (e.g., content_type for text extraction, text length for structured extraction)
- [ ] Verify spans are nil-safe (work correctly when Braintrust tracing is disabled / no TracerProvider set)
- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review

## Key Files

- `src/backend/internal/infrastructure/llm/extraction.go` — main file to modify
- `src/backend/internal/infrastructure/llm/braintrust.go` — reference for existing tracing setup
- `src/backend/cmd/server/main.go:215-224` — where TracerProvider is initialized