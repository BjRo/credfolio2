---
# credfolio2-7s2l
title: Enforce required environment variables with no hardcoded defaults
status: draft
type: task
created_at: 2026-01-30T15:57:07Z
updated_at: 2026-01-30T15:57:07Z
parent: credfolio2-abtx
---

Remove all hardcoded default values for credentials and hosts. Introduce a helper function that reads from environment variables and fails hard if required variables are missing.

## Rationale

Hardcoded defaults for credentials/hosts are a security risk and can lead to:
- Accidentally connecting to wrong services in production
- Leaking default credentials into version control
- Silent fallback to insecure defaults when env vars are misconfigured

The app should fail fast and loud if required configuration is missing.

## Checklist

- [ ] Audit codebase for hardcoded credential/host defaults
- [ ] Create `RequireEnv(key string) string` helper that panics if env var is missing
- [ ] Create `RequireEnvOr(key, fallback string) string` for truly optional vars with safe defaults
- [ ] Replace all credential/host lookups to use `RequireEnv`
- [ ] Document all required environment variables in README or .env.example
- [ ] Update devcontainer/docker-compose with required env vars
- [ ] Write tests for the env helper functions

## Definition of Done
- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review