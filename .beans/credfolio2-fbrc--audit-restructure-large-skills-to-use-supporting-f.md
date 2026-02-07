---
# credfolio2-fbrc
title: Audit & restructure large skills to use supporting files
status: todo
type: task
created_at: 2026-02-07T16:01:38Z
updated_at: 2026-02-07T16:01:38Z
parent: credfolio2-ynmd
---

Audit all skills for size and restructure large ones to use supporting files, keeping SKILL.md focused and under 500 lines.

## Why

Large skills consume significant context when loaded. Claude Code's documentation recommends keeping SKILL.md under 500 lines and moving detailed reference material to separate files that load on-demand. This keeps the skill's core instructions lean while preserving all the detailed content for when it's actually needed.

## Skills to audit

| Skill | Concern |
|-------|---------|
| `vercel-react-best-practices` | 45 rules across 8 categories — likely well over 500 lines |
| `agent-browser` | 60+ commands with detailed usage — very large |
| `web-design-guidelines` | Fetches external content, review size |
| `dev-workflow` | Automated workflow with scripts, review size |
| Others | Quick check, likely fine |

> **Note:** `review-backend` and `review-frontend` have been converted from skills to subagents (`.claude/agents/`). They no longer need to be audited as skills. The `agent-browser` skill is preloaded by the `@qa` subagent but still exists as a standalone skill and should be audited for size.

## What

For each oversized skill:
1. Measure current SKILL.md line count
2. Identify what must stay in SKILL.md (overview, when to use, key instructions)
3. Extract detailed content into supporting files:
   - e.g., `vercel-react-best-practices/rules.md` for the full 45 rules
   - e.g., `agent-browser/commands.md` for the command reference
   - e.g., `agent-browser/examples.md` for usage examples
4. Reference supporting files from SKILL.md so Claude knows when to load them:
   ```markdown
   ## Additional resources
   - For the complete rule set, see [rules.md](rules.md)
   - For command reference, see [commands.md](commands.md)
   ```
5. Verify skills still work correctly after restructuring

## Note

`review-backend`, `review-frontend`, and `agent-browser` have already been restructured as part of the subagent migration. `review-backend` and `review-frontend` are now subagent definitions in `.claude/agents/`. The `agent-browser` skill is preloaded by the `@qa` subagent. Focus restructuring efforts on the remaining skills listed above.

## Definition of Done
- [ ] All skills audited for line count
- [ ] Skills over 500 lines restructured with supporting files
- [ ] SKILL.md files contain overview + references to supporting files
- [ ] Supporting files contain detailed content (rules, commands, examples)
- [ ] All restructured skills tested to verify they still work
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Branch pushed and PR created for human review