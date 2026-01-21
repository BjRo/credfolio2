---
# credfolio2-rztd
title: Implement GraphQL query complexity limiting
status: draft
type: task
created_at: 2026-01-21T15:11:23Z
updated_at: 2026-01-21T15:11:23Z
parent: credfolio2-wxn8
---

Add query complexity analysis to prevent expensive or deeply nested GraphQL queries from overwhelming the server.

## Context
Without complexity limits, attackers can craft queries that consume excessive resources (e.g., deeply nested queries, queries requesting many relations).

## Acceptance Criteria
- [ ] Maximum query depth is enforced
- [ ] Query complexity score is calculated and limited
- [ ] Appropriate error messages for rejected queries
- [ ] Limits are configurable via environment variables

## Implementation Notes
- gqlgen supports complexity limiting via extensions
- Consider both depth limiting and cost-based complexity
- Document the limits for API consumers