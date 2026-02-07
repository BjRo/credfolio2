---
name: refine
description: Launch the refine subagent to create a detailed implementation plan for a bean. Use with a bean ID argument, e.g. /refine credfolio2-abc1
metadata:
  argument-hint: <bean-id>
---

# Refine a Bean

Launch the `@refine` subagent to develop a detailed, actionable implementation plan for a bean in an isolated context.

## How to Launch

Use the Task tool to launch the refine agent:

```
Task tool call:
  subagent_type: "refine"
  description: "Refine <bean-title>"
  prompt: "Refine bean <BEAN_ID>. Read the bean, explore the codebase, and develop a detailed implementation plan. Ask clarifying questions if needed. Update the bean body with the plan."
```

### Adding Extra Context

If the user has provided additional guidance, include it in the prompt:

```
prompt: "Refine bean <BEAN_ID>. <Additional context from user>. Read the bean, explore the codebase, and develop a detailed implementation plan."
```

## What the Agent Does

- Reads the bean and understands the requirements
- Explores the codebase to understand current state, patterns, and dependencies
- Asks clarifying questions via `AskUserQuestion` when requirements are ambiguous
- Writes a structured implementation plan (approach, files, steps, testing strategy)
- Updates the bean body with the plan (preserving existing content)

## What the Agent Does NOT Do

- Modify source code — it is a planning agent only
- Mark the bean as completed
- Change bean status (except todo → in-progress)

## After the Agent Completes

Review the implementation plan it produced, then:

1. **Approve the plan** — proceed to `/implement <bean-id>`
2. **Adjust the plan** — edit the bean body or re-run `/refine` with additional context
3. **Split the bean** — if the plan reveals the scope is too large, break it into child beans
