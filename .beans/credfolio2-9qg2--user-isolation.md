---
# credfolio2-9qg2
title: User Isolation
status: draft
type: epic
created_at: 2026-01-20T11:25:24Z
updated_at: 2026-01-20T11:25:24Z
parent: credfolio2-fq2i
---

Ensure data is properly scoped to users with public profile access.

## Goals
- All data belongs to a specific user
- Public profile URLs for sharing
- Proper authorization checks

## Checklist
- [ ] Add user_id to all relevant tables
- [ ] Update queries to filter by user
- [ ] Add authorization middleware
- [ ] Create public profile route (/u/username)
- [ ] Differentiate owner vs viewer permissions
- [ ] Add profile visibility settings (public/private)
- [ ] Implement username selection