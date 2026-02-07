---
# credfolio2-sxqs
title: Create refine-agent for bean implementation planning
status: in-progress
type: task
priority: normal
created_at: 2026-02-07T15:57:52Z
updated_at: 2026-02-07T16:39:29Z
parent: credfolio2-ynmd
---

Create a new `.claude/agents/refine.md` subagent that takes a bean and develops a detailed implementation plan for it.

## Why

Planning work before implementation is critical for quality. Currently this happens ad-hoc in the main conversation, consuming context and lacking structure. A dedicated refine agent can:
- Explore the codebase thoroughly without polluting main context
- Ask clarifying questions to the user
- Produce a structured plan directly in the bean

## What

- Create `.claude/agents/refine.md` with appropriate frontmatter
- The agent's description should indicate it's used for developing implementation plans for beans
- System prompt should instruct the agent to:
  1. Read the bean (using `beans query`)
  2. Explore the relevant parts of the codebase to understand the current state
  3. Think hard, take things step by step
  4. Ask all necessary questions via AskUserQuestion to ensure a high quality user experience and high quality technical solution
  5. Document the implementation plan in the bean body (using `beans update`)
- Set `model: inherit` (planning needs the best model)
- Grant tools: Read, Bash (for `beans` CLI), Glob, Grep, AskUserQuestion
- The agent should be user-invocable (e.g. "Use the refine agent on bean X")

## Usage Example

```
Use the refine agent on credfolio2-abc1
```

The agent would then read the bean, explore the codebase, ask clarifying questions, and update the bean with a detailed implementation plan.

## Definition of Done
- [ ] Subagent file created at `.claude/agents/refine.md`
- [ ] Agent can read a bean, explore the codebase, and ask clarifying questions
- [ ] Agent updates the bean body with a structured implementation plan
- [ ] Agent works end-to-end on a test bean
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Branch pushed and PR created for human review