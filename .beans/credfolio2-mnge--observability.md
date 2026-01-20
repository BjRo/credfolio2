---
# credfolio2-mnge
title: Observability
status: draft
type: epic
created_at: 2026-01-20T11:25:45Z
updated_at: 2026-01-20T11:25:45Z
parent: credfolio2-abtx
---

Add monitoring, logging, and debugging capabilities.

## Goals
- Structured logging for debugging
- Metrics for monitoring
- Health checks for orchestration

## Checklist
- [ ] Add structured logging (slog or zerolog)
- [ ] Create Prometheus metrics endpoint
- [ ] Add request/response logging middleware
- [ ] Create health check endpoints (liveness, readiness)
- [ ] Add job queue metrics
- [ ] Add LLM call metrics and logging
- [ ] Create error tracking integration (optional)