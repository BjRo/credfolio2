---
# credfolio2-hp81
title: Configure CORS for production
status: draft
type: task
created_at: 2026-01-21T15:11:18Z
updated_at: 2026-01-21T15:11:18Z
parent: credfolio2-wxn8
---

Implement proper CORS (Cross-Origin Resource Sharing) configuration for production deployments.

## Context
CORS prevents unauthorized domains from making requests to the API. Without proper configuration, the API could be accessed from malicious websites.

## Acceptance Criteria
- [ ] CORS middleware is configured with allowed origins
- [ ] Production only allows specific frontend domain(s)
- [ ] Development allows localhost origins
- [ ] Credentials, methods, and headers are properly configured
- [ ] Preflight requests are handled correctly

## Implementation Notes
- Use chi CORS middleware or `rs/cors`
- Configure allowed origins from environment variables
- Consider `CORS_ALLOWED_ORIGINS` env var