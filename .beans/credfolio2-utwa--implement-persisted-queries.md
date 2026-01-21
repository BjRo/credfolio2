---
# credfolio2-utwa
title: Implement persisted queries
status: draft
type: task
created_at: 2026-01-21T15:11:27Z
updated_at: 2026-01-21T15:11:27Z
parent: credfolio2-wxn8
---

Implement persisted (or automatic persisted) queries to improve security and performance.

## Context
Persisted queries allow only pre-approved query shapes to be executed, preventing arbitrary query execution. This also improves performance by reducing request payload sizes.

## Acceptance Criteria
- [ ] Clients can send query hashes instead of full query strings
- [ ] Unknown query hashes are rejected in production
- [ ] Development mode allows query registration
- [ ] Query hash extraction supports Apollo-style APQ format

## Implementation Notes
- Consider Automatic Persisted Queries (APQ) for easier adoption
- gqlgen has APQ support via extensions
- Decide between allowlist (strict) vs APQ (automatic registration)
- May need Redis or similar for distributed query storage