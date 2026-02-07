---
# credfolio2-fvlb
title: Create review-frontend subagent
status: todo
type: task
priority: normal
created_at: 2026-02-07T15:57:44Z
updated_at: 2026-02-07T16:23:33Z
parent: credfolio2-ynmd
blocking:
    - credfolio2-xflj
---

Convert the `review-frontend` skill into a proper `.claude/agents/` subagent, and preload the frontend-specific review skills.

## Why

Same as review-backend: review output is verbose and should be isolated from the main context. Additionally, the `web-design-guidelines` and `vercel-react-best-practices` skills only matter for frontend review — they should be preloaded into this subagent's context rather than being discoverable by the main conversation.

## What

- Create `.claude/agents/review-frontend.md` with appropriate frontmatter
- Move the system prompt from the skill's SKILL.md into the agent's markdown body
- Preload skills via the `skills` frontmatter field:
  - `web-design-guidelines`
  - `vercel-react-best-practices`
- Set `model: inherit`
- Restrict tools to read-only + Bash (for `gh` PR comment posting)
- Remove the old skill (or keep as thin wrapper, TBD)
- Consider whether `web-design-guidelines` and `vercel-react-best-practices` should have `user-invocable: false` and `disable-model-invocation: true` since they're now only consumed by this subagent

## Definition of Done
- [ ] Subagent file created at `.claude/agents/review-frontend.md`
- [ ] `web-design-guidelines` and `vercel-react-best-practices` skills preloaded
- [ ] Old skill removed or converted to delegate to the subagent
- [ ] Subagent can be invoked and posts review comments to a PR
- [ ] Review output stays out of the main conversation context
- [ ] Skills `web-design-guidelines` and `vercel-react-best-practices` reviewed for invocation settings
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Branch pushed and PR created for human review