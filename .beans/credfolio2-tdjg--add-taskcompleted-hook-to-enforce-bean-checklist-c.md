---
# credfolio2-tdjg
title: Add TaskCompleted hook to enforce bean checklist completion
status: todo
type: task
created_at: 2026-02-07T16:19:41Z
updated_at: 2026-02-07T16:19:41Z
parent: credfolio2-ynmd
---

Add a TaskCompleted hook that prevents beans from being marked complete when they still have unchecked checklist items.

## Why

The CLAUDE.md and beans system instructions say "DO NOT mark a bean complete if it has unchecked checklist items" — but this is only enforced by instructions. Claude can (and sometimes does) ignore this. A TaskCompleted hook makes this deterministic: the hook reads the bean, checks for unchecked items, and blocks completion if any remain.

## What

### 1. Hook script

Create `.claude/hooks/validate-bean-completion.sh` that:
- Reads the hook input JSON from stdin (contains `task_id`, `task_subject`, etc.)
- Extracts the bean ID from the task context (task_subject or task_description may contain it)
- Reads the bean body via `beans query`
- Scans for unchecked checklist items (`- [ ]`)
- If unchecked items exist: exits with code 2 and lists the remaining items on stderr
- If all items checked (or no checklist): exits with code 0

### 2. Hook registration

In `.claude/settings.json`:
```json
{
  "hooks": {
    "TaskCompleted": [
      {
        "hooks": [
          {
            "type": "command",
            "command": "\"$CLAUDE_PROJECT_DIR\"/.claude/hooks/validate-bean-completion.sh"
          }
        ]
      }
    ]
  }
}
```

Note: TaskCompleted doesn't support matchers — it fires on every task completion.

### 3. Mapping task to bean

The challenge is mapping the TaskCompleted hook input to a bean. Options to investigate:
- The `task_subject` may contain the bean title
- The current branch name contains the bean ID (use `get-current-bean.sh`)
- The `task_description` may reference the bean
- Fall back to current branch → bean ID if direct mapping isn't available

If no bean can be identified, exit 0 (allow completion — don't block non-bean tasks).

## Example flow

```
# Claude tries to mark a bean complete
# Hook reads bean, finds unchecked items:

stderr: "Cannot complete bean credfolio2-abc1 — 3 unchecked items remain:
  - [ ] Visual verification with agent-browser
  - [ ] Branch pushed and PR created
  - [ ] Automated code review passed
"
exit 2  # blocks completion

# Claude sees the feedback and works on remaining items
```

## Definition of Done
- [ ] Hook script created at `.claude/hooks/validate-bean-completion.sh`
- [ ] Script reads bean body and detects unchecked checklist items
- [ ] Hook registered in `.claude/settings.json`
- [ ] Tested: bean with unchecked items → completion blocked with helpful message
- [ ] Tested: bean with all items checked → completion allowed
- [ ] Tested: non-bean task → completion allowed (graceful fallback)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Branch pushed and PR created for human review