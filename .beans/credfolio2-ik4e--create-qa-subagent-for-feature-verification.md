---
# credfolio2-ik4e
title: Create QA subagent for feature verification
status: in-progress
type: task
priority: normal
created_at: 2026-02-07T15:57:50Z
updated_at: 2026-02-07T16:47:41Z
parent: credfolio2-ynmd
blocking:
    - credfolio2-xflj
---

Create a new `.claude/agents/qa.md` subagent dedicated to feature verification using agent-browser.

## Why

Visual verification with agent-browser produces verbose output (screenshots, DOM snapshots, navigation logs) that consumes significant main context. Delegating verification to a QA subagent keeps this isolated. The `agent-browser` skill only matters for this verification workflow — it should be preloaded into the QA subagent rather than being generally available.

## What

- Create `.claude/agents/qa.md` with appropriate frontmatter
- Preload the `agent-browser` skill via the `skills` frontmatter field
- System prompt should instruct the agent to:
  - Start dev servers if not already running
  - Navigate to the relevant page
  - Verify the feature works (screenshots, interaction, error checks)
  - Report a summary of findings (pass/fail with details)
- Set `model: inherit`
- Grant tools: Read, Bash (for agent-browser CLI + dev server management), Glob, Grep
- Consider whether the `agent-browser` skill should have `disable-model-invocation: true` since it's now only consumed via this subagent
- Update CLAUDE.md references to visual verification to point to the QA subagent

## Definition of Done
- [x] Subagent file created at `.claude/agents/qa.md`
- [x] `agent-browser` skill preloaded into the subagent
- [x] Subagent can verify a feature and return a pass/fail summary
- [x] Verification output stays out of the main conversation context
- [x] `agent-browser` skill reviewed for invocation settings
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [ ] Branch pushed and PR created for human review