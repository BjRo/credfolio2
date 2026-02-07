---
name: implement
description: Launch the implement subagent to execute a refined bean's implementation plan. Use with a bean ID argument, e.g. /implement credfolio2-abc1
metadata:
  argument-hint: <bean-id>
---

# Implement a Bean

Launch the `@implement` subagent to execute a bean's implementation plan in an isolated context.

## Prerequisites

Before launching, verify the bean is ready:

1. **Bean must have an implementation plan** — Run `/refine <bean-id>` first if it doesn't
2. **Bean must have a checklist** — The agent works through checklist items in order

## How to Launch

Use the Task tool to launch the implement agent:

```
Task tool call:
  subagent_type: "implement"
  description: "Implement <bean-title>"
  prompt: "Implement bean <BEAN_ID>. Read the bean, follow its implementation plan using TDD, commit as you go, and update the bean checklist. Report what was completed when done."
```

### Running in Background

For long implementations, run in the background:

```
Task tool call:
  subagent_type: "implement"
  description: "Implement <bean-title>"
  prompt: "Implement bean <BEAN_ID>. Read the bean, follow its implementation plan using TDD, commit as you go, and update the bean checklist. Report what was completed when done."
  run_in_background: true
```

Then check progress by reading the output file.

### Adding Extra Context

If the user has provided additional guidance, include it in the prompt:

```
prompt: "Implement bean <BEAN_ID>. <Additional context from user>. Read the bean, follow its implementation plan using TDD, commit as you go, and update the bean checklist."
```

## After the Agent Completes

The implement agent handles code + tests + commits. You still need to:

1. **Run QA** — Launch `@qa` subagent for visual verification (if UI changes)
2. **Push and create PR** — `git push -u origin <branch>` then `gh pr create`
3. **Run reviews** — Launch `@review-backend` and/or `@review-frontend`
4. **Check off remaining DoD items** — Update the bean's Definition of Done

## What the Agent Does

- Reads the bean and its implementation plan
- Sets up the feature branch (via `start-work.sh`) if not already on one
- Follows TDD: RED → GREEN → REFACTOR for each step
- Commits frequently with meaningful messages
- Updates bean checklist items as they're completed
- Runs `pnpm lint` and `pnpm test` at the end
- Reports a summary of what was done

## What the Agent Does NOT Do

- Create PRs or push branches
- Launch QA or review agents
- Mark the bean as completed
- Merge anything into main
