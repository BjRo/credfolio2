---
# credfolio2-2o0e
title: Introduce logger abstraction for backend
status: completed
type: feature
priority: normal
created_at: 2026-01-22T10:35:06Z
updated_at: 2026-01-22T11:07:45Z
---

Create a logger abstraction that:
- Has mandatory log message (string)
- Has mandatory severity (enum: debug, info, warning, error, critical)
- Has optional feature (string) for categorization
- Has optional data parameter (deep key-value, JSON compatible)

Currently logs to stdout with nice formatting. Designed for future extensibility (Datadog, Sentry, etc.)

## Checklist
- [x] Design and implement logger interface
- [x] Implement stdout logger
- [x] Instrument main error-prone points in the application