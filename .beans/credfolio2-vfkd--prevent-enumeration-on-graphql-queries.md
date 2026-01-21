---
# credfolio2-vfkd
title: Prevent enumeration on GraphQL queries
status: draft
type: task
created_at: 2026-01-21T15:12:32Z
updated_at: 2026-01-21T15:12:32Z
parent: credfolio2-wxn8
---

Prevent attackers from enumerating valid resources (users, files, etc.) through GraphQL query responses.

## Context
Enumeration attacks exploit differences in API responses to discover valid IDs, emails, or other resources. For example, returning "user not found" vs "unauthorized" reveals whether a user exists.

## Acceptance Criteria
- [ ] Queries for non-existent resources return the same response as unauthorized access
- [ ] Error messages don't reveal whether a resource exists
- [ ] Timing differences are minimized between exists/not-exists cases
- [ ] Batch queries don't leak information through partial results

## Implementation Notes
- Return generic "not found or unauthorized" errors
- Consider constant-time comparisons where applicable
- Review all resolver error paths for information leakage
- Test with tools like Burp Suite or custom scripts