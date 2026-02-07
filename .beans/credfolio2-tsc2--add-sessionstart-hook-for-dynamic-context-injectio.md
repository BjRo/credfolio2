---
# credfolio2-tsc2
title: Add SessionStart hook for dynamic context injection
status: todo
type: task
created_at: 2026-02-07T16:08:12Z
updated_at: 2026-02-07T16:08:12Z
parent: credfolio2-ynmd
---

Add a SessionStart hook that injects dynamic work context on every session start, resume, and compaction.

## Why

Claude starts each session (or resumes after compaction) without knowing the current work state — which branch you're on, which bean is active, or what happened recently. This leads to repeated exploration and questions. A SessionStart hook can inject this context automatically so Claude is oriented from the first prompt.

## What

Create a hook script at `.claude/hooks/session-context.sh` that outputs:
- **Current git branch** — so Claude knows if it's on main or a feature branch
- **Active bean** — extracted from the branch name using the existing `.claude/scripts/get-current-bean.sh` pattern (branch format: `<type>/<bean-id>-<description>`)
- **Bean details** — if an active bean is found, fetch its title, status, and checklist via `beans query`
- **Recent commits** — last 5 commits on the current branch for continuity

Register the hook for all SessionStart triggers (startup, resume, compact) so context is always fresh.

## Hook configuration

In `.claude/settings.json`:
```json
{
  "hooks": {
    "SessionStart": [
      {
        "matcher": "",
        "hooks": [
          {
            "type": "command",
            "command": "\"$CLAUDE_PROJECT_DIR\"/.claude/hooks/session-context.sh"
          }
        ]
      }
    ]
  }
}
```

## Example output

The script's stdout gets added as context for Claude:

```
## Current Work Context
- Branch: feature/credfolio2-abc1-add-user-auth
- Active bean: credfolio2-abc1 — "Add user authentication" (in-progress)
- Unchecked items: 3 of 7
- Recent commits:
  - abc1234 feat: add login endpoint
  - def5678 feat: add user model and migration
  - ghi9012 chore: create feature branch
```

## Notes

- The hook should be fast (< 2 seconds) since it runs on every session start
- If on `main` with no active bean, output a short note like "On main branch, no active bean"
- Use the existing `get-current-bean.sh` logic for bean ID extraction
- The hook already registered for `SessionStart` with `beans prime` — add this alongside it

## Definition of Done
- [ ] Hook script created at `.claude/hooks/session-context.sh`
- [ ] Script extracts branch, bean, and recent commits
- [ ] Hook registered in `.claude/settings.json` for all SessionStart events
- [ ] Verified context appears on session start, resume, and after compaction
- [ ] Script runs in under 2 seconds
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Branch pushed and PR created for human review