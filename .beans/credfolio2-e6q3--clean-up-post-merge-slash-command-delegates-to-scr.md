---
# credfolio2-e6q3
title: 'Clean up post-merge: slash command delegates to script'
status: todo
type: task
priority: normal
created_at: 2026-02-07T17:36:16Z
updated_at: 2026-02-07T17:36:20Z
parent: credfolio2-ynmd
---

The `/post-merge` slash command (`.claude/commands/post-merge.md`) and `scripts/post-merge.sh` duplicate the same post-merge cleanup logic. Consolidate so the slash command delegates to the script, and move the script into `.claude/scripts/`.

## Context
- `.claude/commands/post-merge.md` — Claude Code skill with step-by-step instructions
- `scripts/post-merge.sh` — standalone bash script with identical logic

## Checklist
- [ ] Move `scripts/post-merge.sh` to `.claude/scripts/post-merge.sh`
- [ ] Update `.claude/commands/post-merge.md` to instruct Claude to call `.claude/scripts/post-merge.sh <bean-id>` instead of manually executing each step
- [ ] Ensure the slash command still validates the `$ARGUMENTS` bean ID before calling the script
- [ ] Remove duplicated step-by-step instructions from the slash command
- [ ] Remove empty `scripts/` directory if nothing else remains in it

## Definition of Done
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review
- [ ] Automated code review passed via `@review-backend` and/or `@review-frontend` subagents (via Task tool)