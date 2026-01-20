---
# credfolio2-tikg
title: Foundation & First Vertical Slice
status: draft
type: milestone
priority: normal
created_at: 2026-01-20T11:24:04Z
updated_at: 2026-01-20T11:26:54Z
blocking:
    - credfolio2-1fwu
---

Upload a reference letter → extract text via LLM → display raw results. This milestone establishes the core infrastructure and proves out the main user flow end-to-end.

## Development Process

All work in this milestone follows these practices:

1. **Feature Branches**: Create from up-to-date main (`git checkout main && git pull origin main`)
2. **TDD**: Write tests first, then implementation (`/skill tdd`)
3. **PR Review**: Open PRs for human review - do NOT merge yourself
4. **Track Progress**: Update bean checklists as you complete items

See `/skill dev-workflow` for the complete workflow.