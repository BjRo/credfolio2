---
# credfolio2-3hz3
title: Add code review agent skills and integrate into dev-workflow
status: in-progress
type: feature
created_at: 2026-02-06T08:24:44Z
updated_at: 2026-02-06T08:24:44Z
---

## Summary

Create two specialized code review agent skills and integrate them into the dev-workflow:

1. **review-backend**: A staff-level Go/Backend developer agent that reviews backend code (Go, GraphQL API) focusing on maintainability, design, performance, and security. Posts findings as PR comments via `gh`.

2. **review-frontend**: A staff-level React/Next.js Frontend developer agent that reviews frontend code focusing on frontend best practices. Posts findings as PR comments via `gh`.

3. **Update dev-workflow**: Add a code review step after PR creation that invokes both agents.

4. **Update bean template**: Add code review checkpoint to the Definition of Done.

## Checklist
- [x] Create `review-backend` skill (`.claude/skills/review-backend/SKILL.md`)
- [x] Create `review-frontend` skill (`.claude/skills/review-frontend/SKILL.md`)
- [x] Update `dev-workflow` skill to include code review step
- [x] Update bean "Definition of Done" template in dev-workflow and CLAUDE.md

## Definition of Done
- [x] All checklist items above are completed
- [ ] Branch pushed and PR created for human review