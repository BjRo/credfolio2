---
# credfolio2-h0uy
title: Verify pre-commit hook works in both local and devcontainer environments
status: todo
type: task
created_at: 2026-02-07T16:15:18Z
updated_at: 2026-02-07T16:15:18Z
parent: credfolio2-ynmd
---

Verify the existing pre-commit-check hook works reliably in both local and devcontainer environments.

## Why

The pre-commit hook (`.claude/hooks/pre-commit-check.sh`) runs lint and tests before every git commit. If it silently fails or doesn't trigger in one environment, broken code can get committed. We need confidence that this enforcement is consistent everywhere.

## What to investigate

### 1. Hook registration
- Is the hook registered in `.claude/settings.json` (project-level, shared)?
- Does it also need to be in `.claude/settings.local.json` or `~/.claude/settings.json`?
- Does the devcontainer's `bypassPermissions` mode affect hook execution? (It shouldn't — hooks are separate from permissions — but verify)

### 2. Script portability
- Read `.claude/hooks/pre-commit-check.sh` and check for:
  - Hardcoded paths that might differ between environments
  - Assumptions about shell (bash vs zsh)
  - Dependencies (`pnpm`, `jq`, etc.) available in both environments
  - Correct shebang line
  - Executable bit set (`chmod +x`)

### 3. Hook input parsing
- The hook receives JSON on stdin with the Bash command
- Verify the matcher correctly catches all commit variants: `git commit`, `git commit -m`, `git commit --amend`, etc.
- Check edge cases: `git commit` inside a pipe, with `&&`, or via aliases

### 4. Actual testing
- Test in the devcontainer:
  - Stage a file with a lint error → attempt commit → verify it's blocked
  - Stage clean code → attempt commit → verify it passes
  - Check hook output is visible (stderr messages)
- Test locally (if applicable):
  - Same checks as above

### 5. Failure modes
- What happens if `pnpm lint` or `pnpm test` hangs? Is there a timeout?
- What happens if the hook script itself errors (e.g., jq not found)? Does the commit proceed or block?
- Non-zero exit codes other than 2 — do they silently allow the commit?

## Definition of Done
- [ ] Hook registration verified in both environments
- [ ] Script reviewed for portability issues
- [ ] Hook input parsing verified for all commit command variants
- [ ] End-to-end test: lint failure blocks commit in devcontainer
- [ ] End-to-end test: test failure blocks commit in devcontainer
- [ ] End-to-end test: clean code allows commit in devcontainer
- [ ] Any issues found are fixed
- [ ] Failure modes documented or addressed (e.g., timeout added)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Branch pushed and PR created for human review