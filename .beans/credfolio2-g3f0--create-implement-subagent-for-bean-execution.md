---
# credfolio2-g3f0
title: Create implement subagent for bean execution
status: in-progress
type: feature
priority: normal
created_at: 2026-02-07T19:59:16Z
updated_at: 2026-02-07T20:03:47Z
parent: credfolio2-ynmd
---

Create a dedicated subagent that executes the implementation of a well-refined bean, keeping the main conversation context clean for high-level steering.

## Motivation

Implementation is the most context-hungry phase of the dev workflow (reading files, writing code, running tests, iterating on failures). Currently this all happens in the main conversation, which:
- Consumes context window with low-level details
- Makes it harder to maintain strategic oversight
- Mixes implementation noise with decision-making

We already have subagents for every other phase:
- `refine` — plans the work
- `qa` — verifies UI changes
- `review-backend` / `review-frontend` — reviews code

An "implement" agent fills the gap: **refine → implement → qa → review**.

## Design Considerations

- **Input**: Bean ID (to read the bean body/checklist), branch name, and any additional context
- **Agent type**: `general-purpose` (needs Edit, Write, Bash, Read, Glob, Grep, AskUserQuestion)
- **Scope**: The agent should follow the bean checklist, write code, run tests, and update checklist items as it goes
- **Boundaries**: Should NOT create PRs, run review agents, or mark beans complete — those remain in the main context
- **Background mode**: Should support `run_in_background: true` for non-blocking execution
- **Failure handling**: If blocked, should use AskUserQuestion to surface issues to the user

## Checklist

- [x] Research existing subagent configurations (refine, qa, review-backend, review-frontend) for patterns
- [x] Design the implement agent prompt template
- [x] Create the agent configuration in `.claude/agents/`
- [x] Create a skill (`/implement`) that invokes the agent with a bean ID
- [ ] Test with a real bean implementation
- [ ] Document usage in CLAUDE.md or relevant rules file

## Definition of Done
- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Visual verification via `@qa` subagent (via Task tool, for UI changes)
- [ ] ADR written via `/decision` skill (if new dependencies, patterns, or architectural changes were introduced)
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review
- [ ] Automated code review passed via `@review-backend` and/or `@review-frontend` subagents (via Task tool)