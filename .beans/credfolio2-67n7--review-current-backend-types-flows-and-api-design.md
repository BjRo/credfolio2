---
# credfolio2-67n7
title: Review current backend types, flows, and API design
status: todo
type: task
priority: normal
created_at: 2026-02-06T12:03:21Z
updated_at: 2026-02-06T12:04:12Z
---

Conduct a thorough review of the backend codebase to assess clarity, consistency, and correctness of types, flows, and the GraphQL API surface.

## Approach

Use the `@review-backend` subagent (via Task tool) to perform the review. This subagent provides staff-level Go/Backend code review covering maintainability, design, performance, and security. Run it against the backend codebase (not tied to a PR) to get a comprehensive assessment.

## Goals

1. **Type review** — Audit Go types (models, DTOs, GraphQL schema types) for clarity of naming, purpose, and freedom from redundancy. Are there overlapping types that could be consolidated? Are type boundaries between layers (DB models, domain types, GraphQL types) clean?

2. **Flow review** — Trace the key processing flows end-to-end (document upload, detection, extraction, import, reference letter handling, etc.) and assess whether each step has a clear responsibility. Look for unnecessary indirection, duplicated logic, or confusing handoffs between components.

3. **API review** — Evaluate the GraphQL schema (queries, mutations, input types, enums) for consistency, naming conventions, and whether the API surface accurately represents the domain. Are there mutations that could be consolidated or queries that expose internal implementation details?

4. **Interaction graph** — Document the current system interactions as a graph (services, workers, storage, LLM calls, database) showing how data flows through the system. Identify any flows that seem overly complex or that don't make sense.

## Key areas to examine

- `src/backend/internal/` — domain types, services, resolvers
- `src/backend/graph/` — GraphQL schema and resolver implementations
- `src/backend/cmd/server/` — server setup and wiring
- `src/backend/migrations/` — database schema (for understanding the data model)

## Checklist
- [ ] Run `@review-backend` subagent (via Task tool) against the backend codebase
- [ ] Audit Go types across layers (models, services, GraphQL) for clarity and redundancy
- [ ] Trace and document each major flow (upload, detect, extract, import, reference letters)
- [ ] Review GraphQL schema for consistency, naming, and API design
- [ ] Produce an interaction graph showing services, data flows, and external dependencies
- [ ] Document findings: what's clean, what's confusing, what could be improved
- [ ] Summarize actionable recommendations (if any) as follow-up beans