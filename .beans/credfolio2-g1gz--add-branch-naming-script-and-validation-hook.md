---
# credfolio2-g1gz
title: Add branch naming script and validation hook
status: todo
type: task
created_at: 2026-02-07T16:13:03Z
updated_at: 2026-02-07T16:13:03Z
parent: credfolio2-ynmd
---

Create a script that generates consistent branch names from beans, and a PreToolUse hook that validates branch names on creation.

## Why

The dev-workflow requires branch names to follow the pattern `<type>/<bean-id>-<description>` (e.g., `feature/credfolio2-abc1-add-user-auth`). Currently this is only enforced by instructions in the dev-workflow skill — Claude sometimes creates inconsistent branch names. A script + hook combination makes this deterministic.

## What

### 1. Branch name creation script

Create `.claude/scripts/create-branch.sh` that:
- Takes a bean ID as argument
- Queries the bean for its type and title via `beans query`
- Generates a slugified branch name: `<type>/<bean-id>-<slugified-title>`
- Creates and switches to the branch
- Example: bean "Add user authentication" (type: feature, id: credfolio2-abc1) → `feature/credfolio2-abc1-add-user-authentication`

Slugification rules:
- Lowercase
- Spaces → hyphens
- Strip special characters
- Truncate to reasonable length (e.g., 60 chars total)

### 2. PreToolUse validation hook

Create `.claude/hooks/validate-branch-name.sh` that:
- Intercepts `git checkout -b`, `git switch -c`, and `git branch` commands
- Extracts the proposed branch name
- Validates it matches the pattern: `<type>/<bean-id>-<description>` where type is one of: feature, bug, task, milestone, epic
- Exits with code 2 (blocks) if the pattern doesn't match, with a helpful error message pointing to the create-branch script
- Allows other git commands to pass through unblocked

Register as a PreToolUse hook in `.claude/settings.json`:
```json
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Bash",
        "hooks": [
          {
            "type": "command",
            "command": "\"$CLAUDE_PROJECT_DIR\"/.claude/hooks/validate-branch-name.sh"
          }
        ]
      }
    ]
  }
}
```

### 3. Integration

- Update the dev-workflow skill to reference the create-branch script
- Claude should use `.claude/scripts/create-branch.sh <bean-id>` instead of manually constructing branch names

## Example flow

```bash
# Claude runs:
.claude/scripts/create-branch.sh credfolio2-abc1

# Script queries bean, outputs:
# Creating branch: feature/credfolio2-abc1-add-user-authentication
# Switched to new branch 'feature/credfolio2-abc1-add-user-authentication'

# If Claude tries to create a branch manually with wrong format:
git checkout -b my-feature
# Hook blocks: "Branch name 'my-feature' doesn't match required pattern.
#   Use: .claude/scripts/create-branch.sh <bean-id>"
```

## Note

The validation hook will share the PreToolUse(Bash) matcher with the existing pre-commit hook. Both hooks will run in parallel for Bash commands — ensure they don't conflict. The validation hook should only inspect branch-creation commands and exit 0 for everything else.

## Definition of Done
- [ ] `.claude/scripts/create-branch.sh` created and working
- [ ] `.claude/hooks/validate-branch-name.sh` created and working
- [ ] Hook registered in `.claude/settings.json`
- [ ] Dev-workflow skill updated to reference the script
- [ ] Existing pre-commit hook still works alongside the new hook
- [ ] Tested: valid branch names pass, invalid ones are blocked with helpful message
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Branch pushed and PR created for human review