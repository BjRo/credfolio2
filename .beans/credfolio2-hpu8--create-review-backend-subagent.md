---
# credfolio2-hpu8
title: Create review-backend subagent
status: todo
type: task
priority: normal
created_at: 2026-02-07T15:57:37Z
updated_at: 2026-02-07T16:23:33Z
parent: credfolio2-ynmd
blocking:
    - credfolio2-xflj
---

Convert the `review-backend` skill into a proper `.claude/agents/` subagent.

## Why

The review-backend skill currently runs in the main conversation context. Review output is verbose (PR comments, code analysis) and pollutes the main context window. Running as a subagent isolates this output — only a summary returns to the caller.

## What

- Create `.claude/agents/review-backend.md` with appropriate frontmatter
- Move the system prompt from the skill's SKILL.md into the agent's markdown body
- Set `model: inherit` (reviews need the best model available)
- Restrict tools to read-only + Bash (for `gh` PR comment posting)
- Remove the old skill (or keep as a thin wrapper that delegates, TBD)

## Definition of Done
- [ ] Subagent file created at `.claude/agents/review-backend.md`
- [ ] Old skill removed or converted to delegate to the subagent
- [ ] Subagent can be invoked and posts review comments to a PR
- [ ] Review output stays out of the main conversation context
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Branch pushed and PR created for human review