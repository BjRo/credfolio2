---
# credfolio2-g1gz
title: Add branch naming script and validation hook
status: in-progress
type: task
priority: normal
created_at: 2026-02-07T16:13:03Z
updated_at: 2026-02-07T18:10:08Z
parent: credfolio2-ynmd
---

Create a script that generates consistent branch names from beans, and a PreToolUse hook that validates branch names on creation. Also migrate the existing `scripts/start-work.sh` into `.claude/scripts/` to consolidate Claude workflow scripts in one place.

## Why

The dev-workflow requires branch names to follow the pattern `<type>/<bean-id>-<description>` (e.g., `feat/credfolio2-abc1-add-user-auth`). Currently this is only enforced by instructions in the dev-workflow skill — Claude sometimes creates inconsistent branch names. A script + hook combination makes this deterministic.

Additionally, the existing `scripts/start-work.sh` is a Claude workflow script that belongs in `.claude/scripts/` alongside `get-current-bean.sh`, not in the root `scripts/` directory (which should be reserved for infrastructure scripts like `init-db.sh`).

## What

### 1. Move and refactor start-work.sh into .claude/scripts/

Move `scripts/start-work.sh` to `.claude/scripts/start-work.sh` and refactor it so that:
- It takes only a bean ID as argument (not type + description)
- It auto-derives the branch type from the bean's type field (feature->feat, bug->fix, task->chore, etc.)
- It auto-generates the slugified description from the bean's title
- It retains all existing functionality (ensure main up-to-date, verify bean exists, create branch, mark in-progress, commit status change)

### 2. PreToolUse validation hook

Create `.claude/hooks/validate-branch-name.sh` that:
- Intercepts `git checkout -b`, `git switch -c`, and `git branch` commands
- Extracts the proposed branch name
- Validates it matches the pattern: `<type>/<bean-id>-<description>` where type is one of: feat, fix, refactor, chore, docs
- Exits with code 2 (blocks) if the pattern doesn't match, with a helpful error message pointing to the start-work script
- Allows other git commands to pass through unblocked

Register as a PreToolUse hook in `.claude/settings.json`.

### 3. Integration

- Update the dev-workflow skill to reference `.claude/scripts/start-work.sh <bean-id>` (new single-argument interface)
- Remove old `scripts/start-work.sh`
- Remove `scripts/` directory if only `init-db.sh` remains (it stays because it is referenced by docker-compose)

## Example flow

```bash
# Claude runs:
.claude/scripts/start-work.sh credfolio2-abc1

# Script queries bean, derives type=feature title="Add user authentication", outputs:
# Creating branch: feat/credfolio2-abc1-add-user-authentication
# Switched to new branch 'feat/credfolio2-abc1-add-user-authentication'

# If Claude tries to create a branch manually with wrong format:
git checkout -b my-feature
# Hook blocks: "Branch name 'my-feature' doesn't match required pattern.
#   Use: .claude/scripts/start-work.sh <bean-id>"
```

## Note

The validation hook will share the PreToolUse(Bash) matcher with the existing pre-commit hook. Both hooks will run in parallel for Bash commands — ensure they don't conflict. The validation hook should only inspect branch-creation commands and exit 0 for everything else.

## Relationship to other beans

- `credfolio2-e6q3` ("Clean up post-merge: slash command delegates to script") handles moving `scripts/post-merge.sh` to `.claude/scripts/`. That is out of scope for this bean.
- `scripts/init-db.sh` is an infrastructure script referenced by docker-compose and should NOT be moved.
- After both this bean and `credfolio2-e6q3` are done, the `scripts/` directory should contain only `init-db.sh`.

## Checklist
- [ ] Move `scripts/start-work.sh` to `.claude/scripts/start-work.sh`
- [ ] Refactor to accept only a bean ID (auto-derive type and description from bean metadata)
- [ ] Create `.claude/hooks/validate-branch-name.sh`
- [ ] Register the validation hook in `.claude/settings.json`
- [ ] Update dev-workflow skill (`/workspace/.claude/skills/dev-workflow/SKILL.md`) to reference the new script location and interface
- [ ] Tested: valid branch names pass, invalid ones are blocked with helpful message
- [ ] Existing pre-commit hook still works alongside the new hook

## Implementation Plan

### Approach

Refactor the existing `scripts/start-work.sh` into `.claude/scripts/start-work.sh` with a simplified interface (single bean-id argument), then add a PreToolUse validation hook that blocks malformed branch names. The branch type mapping uses conventional git prefixes (`feat`, `fix`, `chore`, `docs`, `refactor`) derived from bean type fields (`feature`, `bug`, `task`, `milestone`, `epic`).

### Files to Create/Modify

- `.claude/scripts/start-work.sh` — **CREATE** (moved from `scripts/start-work.sh` and refactored)
- `.claude/hooks/validate-branch-name.sh` — **CREATE** (new PreToolUse hook)
- `.claude/settings.json` — **MODIFY** (register the new hook)
- `.claude/skills/dev-workflow/SKILL.md` — **MODIFY** (update script references and examples)
- `scripts/start-work.sh` — **DELETE** (moved to `.claude/scripts/`)

### Steps

#### 1. Create `.claude/scripts/start-work.sh` (refactored from `scripts/start-work.sh`)

Start from the existing script at `/workspace/scripts/start-work.sh` (88 lines), but change the interface and add auto-derivation:

**Interface change:**
- Old: `./scripts/start-work.sh <bean-id> <type> <short-description>` (3 args)
- New: `.claude/scripts/start-work.sh <bean-id>` (1 arg)

**New logic to add:**

1. Query bean metadata: `beans query '{ bean(id: "<bean-id>") { id title status type } }' --json`
2. Extract `type` and `title` fields using `jq`
3. Map bean type to branch prefix:
   - `feature` -> `feat`
   - `bug` -> `fix`
   - `task` -> `chore`
   - `milestone` -> `chore`
   - `epic` -> `chore`
   - Any unrecognized type -> `chore` (safe fallback)
4. Slugify the title:
   - Lowercase: `tr '[:upper:]' '[:lower:]'`
   - Replace spaces and underscores with hyphens
   - Strip non-alphanumeric/hyphen characters: `sed 's/[^a-z0-9-]//g'`
   - Collapse multiple hyphens: `sed 's/-\+/-/g'`
   - Strip leading/trailing hyphens: `sed 's/^-//;s/-$//'`
   - Truncate slug so total branch name (`<prefix>/<bean-id>-<slug>`) stays under 72 characters
5. Construct branch name: `${PREFIX}/${BEAN_ID}-${SLUG}`
6. Keep all existing steps: checkout main, pull, verify bean, create branch, mark in-progress, commit

**Pattern to follow:** Look at `/workspace/.claude/scripts/get-current-bean.sh` for the style conventions used in `.claude/scripts/` (shebang, `set -e`, comments describing the script purpose and usage).

**Important details:**
- Use `"${CLAUDE_PROJECT_DIR:-/workspace}"` for any path references, matching the pattern in `pre-commit-check.sh`
- The `--json` flag on beans query returns raw JSON suitable for `jq` parsing
- Validate that the bean status is not `completed` (refuse to start work on completed beans)
- The `git commit` inside this script will trigger the pre-commit hook. Since we are only committing `.beans/` files (no lint/test-relevant changes), ensure the pre-commit hook does not get stuck on this. Looking at `/workspace/.claude/hooks/pre-commit-check.sh`, it runs `pnpm lint` and `pnpm test` for every commit -- this is fine, those should pass since we are not changing source code.

#### 2. Create `.claude/hooks/validate-branch-name.sh`

Create a new PreToolUse hook at `/workspace/.claude/hooks/validate-branch-name.sh`.

**Input contract** (same as `pre-commit-check.sh`): JSON on stdin with `tool_input.command` field.

**Logic:**

1. Read JSON input: `INPUT=$(cat)`, `COMMAND=$(echo "$INPUT" | jq -r '.tool_input.command // ""')`
2. Check if the command creates a new branch. Match these patterns:
   - `git checkout -b <name>` (also with flags between: `git checkout -B <name>`)
   - `git switch -c <name>` (also `-C`)
   - `git branch <name>` (but NOT `git branch -d`, `git branch -D`, `git branch --delete`, `git branch -a`, `git branch -l`, `git branch --list`, `git branch --show-current`, etc.)
3. If not a branch-creation command: `exit 0` (allow)
4. Extract the proposed branch name from the command
5. Validate it matches: `^(feat|fix|refactor|chore|docs)/(credfolio2-[a-zA-Z0-9]+|beans-[a-zA-Z0-9]+)-.+$`
   - Note: include `beans-` prefix as well, since `get-current-bean.sh` already supports that pattern
6. If valid: `exit 0`
7. If invalid: print helpful error to stderr and `exit 2`

**Error message format:**
```
Branch name '<name>' does not match the required pattern.
Expected: <type>/<bean-id>-<description>
  Types: feat, fix, refactor, chore, docs
  Example: feat/credfolio2-abc1-add-user-auth

Use the start-work script instead:
  .claude/scripts/start-work.sh <bean-id>
```

**Edge cases to handle:**
- Commands with `&&` chaining: the branch creation might appear after `&&` (e.g., `git checkout main && git checkout -b my-branch`). The regex should handle this.
- Quoted branch names: `git checkout -b "my branch"` -- unlikely but handle gracefully
- The `git branch` command without `-b` flag can also create branches, but it is ambiguous (e.g., `git branch` alone just lists branches). Only match when there is a name argument that is not a flag.

**Pattern to follow:** Mirror the structure of `/workspace/.claude/hooks/pre-commit-check.sh`:
- Same header comment style explaining purpose, input/output contract
- Same `set -e`, `INPUT=$(cat)`, `jq` extraction pattern
- Same exit code conventions (0 = allow, 2 = block)

#### 3. Register the hook in `.claude/settings.json`

Modify `/workspace/.claude/settings.json` to add the new hook. The current file has a single PreToolUse entry with matcher "Bash" running `pre-commit-check.sh`. Add the validation hook as a second entry in the same hooks array.

**Current structure:**
```json
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Bash",
        "hooks": [
          {
            "type": "command",
            "command": "$CLAUDE_PROJECT_DIR/.claude/hooks/pre-commit-check.sh"
          }
        ]
      }
    ]
  }
}
```

**Target structure:** Add a second PreToolUse entry with the same "Bash" matcher:
```json
{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Bash",
        "hooks": [
          {
            "type": "command",
            "command": "$CLAUDE_PROJECT_DIR/.claude/hooks/pre-commit-check.sh"
          }
        ]
      },
      {
        "matcher": "Bash",
        "hooks": [
          {
            "type": "command",
            "command": "$CLAUDE_PROJECT_DIR/.claude/hooks/validate-branch-name.sh"
          }
        ]
      }
    ]
  }
}
```

Both hooks run for every Bash command but each one returns early (exit 0) for commands it does not care about, so they will not conflict.

#### 4. Update the dev-workflow skill

Modify `/workspace/.claude/skills/dev-workflow/SKILL.md`:

**Section "Quick Start (Recommended)"** (lines 16-30):
- Change from: `./scripts/start-work.sh <bean-id> <type> <short-description>`
- Change to: `.claude/scripts/start-work.sh <bean-id>`
- Update the examples to show the single-argument interface
- Update the description to mention auto-derivation of type and description

**Section "Manual Steps (Reference)"** (lines 32-59):
- Keep this section but add a note that the script is the preferred approach
- Ensure the branch type list is consistent (`feat`, `fix`, `refactor`, `chore`, `docs`)

**Section "After Merge: Complete the Bean"** (lines 198-214):
- Keep references to `./scripts/post-merge.sh` as-is (that migration is handled by `credfolio2-e6q3`)

**Quick Reference section** (lines 225-251):
- Update the start-work line to match the new script path and interface

#### 5. Delete old `scripts/start-work.sh`

Remove `/workspace/scripts/start-work.sh`. The `scripts/` directory will still contain `init-db.sh` and `post-merge.sh` (until `credfolio2-e6q3` is implemented), so do NOT remove the directory.

#### 6. Make scripts executable

Ensure both new scripts have the executable bit set:
```bash
chmod +x .claude/scripts/start-work.sh
chmod +x .claude/hooks/validate-branch-name.sh
```

### Testing Strategy

**Automated tests (shell script verification):**

Since these are bash scripts (not Go or TypeScript), testing is done via manual execution. No unit test framework applies, but the following verification steps should be performed:

1. **start-work.sh - happy path:**
   - Pick an existing todo bean (e.g., `credfolio2-duun` "Fix Makefile help output showing wrong target names", type: task)
   - Run `.claude/scripts/start-work.sh credfolio2-duun`
   - Verify it creates branch `chore/credfolio2-duun-fix-makefile-help-output-showing-wrong-target-names` (or truncated version)
   - Verify bean is marked in-progress
   - Clean up: `git checkout main && git branch -D <branch> && beans update credfolio2-duun --status todo && git add .beans/ && git commit --amend --no-edit`

2. **start-work.sh - error cases:**
   - Run with no arguments: should print usage and exit 1
   - Run with nonexistent bean: should print error and exit 1

3. **validate-branch-name.sh - valid names pass:**
   - Feed it JSON input simulating `git checkout -b feat/credfolio2-abc1-my-feature`
   - Verify exit code 0

4. **validate-branch-name.sh - invalid names blocked:**
   - Feed it JSON input simulating `git checkout -b my-feature`
   - Verify exit code 2 and helpful error message on stderr

5. **validate-branch-name.sh - non-branch commands pass:**
   - Feed it JSON input simulating `git status`, `git commit -m "test"`, `pnpm build`
   - Verify exit code 0 for all

6. **Both hooks coexist:**
   - Verify `pnpm lint` still passes
   - Verify `pnpm test` still passes
   - Verify the pre-commit hook still fires on `git commit` commands

### Open Questions

None -- all ambiguities have been resolved by the user's clarification.

## Definition of Done
- [ ] `.claude/scripts/start-work.sh` created (moved and refactored from `scripts/start-work.sh`)
- [ ] `.claude/hooks/validate-branch-name.sh` created and working
- [ ] Hook registered in `.claude/settings.json`
- [ ] Dev-workflow skill updated to reference the new script location and interface
- [ ] Old `scripts/start-work.sh` deleted
- [ ] Existing pre-commit hook still works alongside the new hook
- [ ] Tested: valid branch names pass, invalid ones are blocked with helpful message
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Branch pushed and PR created for human review
- [ ] Automated code review passed via `@review-backend` and/or `@review-frontend` subagents (via Task tool)
