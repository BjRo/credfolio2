---
# credfolio2-tsc2
title: Add SessionStart hook for dynamic context injection
status: completed
type: task
priority: normal
created_at: 2026-02-07T16:08:12Z
updated_at: 2026-02-07T18:31:59Z
parent: credfolio2-ynmd
---

Add a SessionStart hook that injects dynamic work context on every session start, resume, and compaction.

## Why

Claude starts each session (or resumes after compaction) without knowing the current work state — which branch you're on, which bean is active, or what happened recently. This leads to repeated exploration and questions. A SessionStart hook can inject this context automatically so Claude is oriented from the first prompt.

## What

Create a hook script at `.claude/hooks/session-context.sh` that outputs:
- **Current git branch** — so Claude knows if it's on main or a feature branch
- **Active bean** — extracted from the branch name using the existing `.claude/scripts/get-current-bean.sh` pattern (branch format: `<type>/<bean-id>-<description>`)
- **Bean details** — if an active bean is found, fetch its title, status, and checklist via `beans graphql`
- **Recent commits** — last 5 commits on the current branch for continuity

Register the hook for all SessionStart triggers (startup, resume, compact) so context is always fresh.

## Hook configuration

The new hook must be **added alongside** the existing `beans prime` hook in `.claude/settings.json`, not replace it. The complete merged `SessionStart` section should look like:

```json
"SessionStart": [
  {
    "hooks": [
      {
        "type": "command",
        "command": "beans prime"
      }
    ]
  },
  {
    "hooks": [
      {
        "type": "command",
        "command": "$CLAUDE_PROJECT_DIR/.claude/hooks/session-context.sh"
      }
    ]
  }
]
```

Note: The existing `PreCompact` hook runs `beans prime` before compaction. `SessionStart` fires on startup, resume, **and** after compaction — so registering the context hook only under `SessionStart` is sufficient. No `PreCompact` registration is needed.

## Example output

The script's stdout gets added as context for Claude:

```
## Current Work Context
- Branch: feat/credfolio2-abc1-add-user-auth
- Active bean: credfolio2-abc1 — "Add user authentication" (in-progress)
- Unchecked items: 3 of 7
- Recent commits:
  - abc1234 feat: add login endpoint
  - def5678 feat: add user model and migration
  - ghi9012 chore: create feature branch
```

When on main with no active bean:

```
## Current Work Context
- Branch: main
- No active bean (on main branch)
- Recent commits:
  - 8ff77f0 chore: Mark credfolio2-g1gz as completed
  - 536f34e chore: Add branch naming script and validation hook
  - ed536f7 chore: Add bean for post-merge cleanup consolidation
```

When in detached HEAD state:

```
## Current Work Context
- Branch: (detached HEAD)
- No active bean
- Recent commits:
  - <latest commits>
```

## Notes

- The hook should be fast (< 2 seconds) since it runs on every session start
- If on `main` with no active bean, output a short note like "On main branch, no active bean"
- Use the existing `get-current-bean.sh` logic for bean ID extraction
- The hook already registered for `SessionStart` with `beans prime` — add this alongside it

## Implementation Plan

### Approach

Create a single bash script that gracefully gathers git and beans context, never failing even when data is unavailable. The script inlines the bean ID extraction logic (rather than calling `get-current-bean.sh`) to avoid its `set -e` / `exit 1` behavior which would kill the hook on main. Register the hook alongside the existing `beans prime` SessionStart hook.

### Files to Create/Modify

- `.claude/hooks/session-context.sh` — **Create**: The main hook script
- `.claude/settings.json` — **Modify**: Add the new hook entry to the existing `SessionStart` array

### Steps

1. **Create the hook script at `.claude/hooks/session-context.sh`**

   The script should follow these patterns from existing hooks:
   - Use `#!/bin/bash` shebang (consistent with `pre-commit-check.sh`, `validate-branch-name.sh`)
   - Use `${CLAUDE_PROJECT_DIR:-/workspace}` for project directory (pattern from `pre-commit-check.sh` line 25)
   - Output context to **stdout** (SessionStart hooks inject stdout as system context)
   - Send any debug/error messages to **stderr** (they will not be seen by Claude)
   - **Never exit with non-zero code** — use `|| true` on fallible commands and avoid `set -e`

   Script structure:

   ```bash
   #!/bin/bash
   # Claude Code SessionStart hook: Injects current work context
   # Output goes to stdout and is injected as system context for Claude.

   PROJECT_DIR="${CLAUDE_PROJECT_DIR:-/workspace}"
   cd "$PROJECT_DIR" || exit 0

   echo "## Current Work Context"

   # --- Git branch ---
   BRANCH=$(git branch --show-current 2>/dev/null)
   if [ -z "$BRANCH" ]; then
       echo "- Branch: (detached HEAD)"
   else
       echo "- Branch: $BRANCH"
   fi

   # --- Active bean (inline extraction, do NOT call get-current-bean.sh) ---
   BEAN_ID=""
   if [ -n "$BRANCH" ] && [ "$BRANCH" != "main" ] && [ "$BRANCH" != "master" ]; then
       if [[ "$BRANCH" =~ ^[a-z]+/(credfolio2-[a-zA-Z0-9]+)-.* ]]; then
           BEAN_ID="${BASH_REMATCH[1]}"
       elif [[ "$BRANCH" =~ ^[a-z]+/(beans-[a-zA-Z0-9]+)-.* ]]; then
           BEAN_ID="${BASH_REMATCH[1]}"
       fi
   fi

   # --- Bean details (only if we have a bean ID) ---
   if [ -n "$BEAN_ID" ]; then
       BEAN_JSON=$(timeout 2s beans graphql "{ bean(id: \"${BEAN_ID}\") { title status body } }" --json 2>/dev/null || echo "")
       if [ -n "$BEAN_JSON" ]; then
           BEAN_TITLE=$(echo "$BEAN_JSON" | jq -r '.bean.title // empty' 2>/dev/null)
           BEAN_STATUS=$(echo "$BEAN_JSON" | jq -r '.bean.status // empty' 2>/dev/null)
           BEAN_BODY=$(echo "$BEAN_JSON" | jq -r '.bean.body // empty' 2>/dev/null)

           if [ -n "$BEAN_TITLE" ]; then
               echo "- Active bean: ${BEAN_ID} — \"${BEAN_TITLE}\" (${BEAN_STATUS})"

               # Count checklist items
               TOTAL=$(echo "$BEAN_BODY" | grep -c '^\- \[[ x]\]' 2>/dev/null || echo "0")
               UNCHECKED=$(echo "$BEAN_BODY" | grep -c '^\- \[ \]' 2>/dev/null || echo "0")
               if [ "$TOTAL" -gt 0 ]; then
                   echo "- Unchecked items: ${UNCHECKED} of ${TOTAL}"
               fi
           else
               echo "- Active bean: ${BEAN_ID} (could not fetch details)"
           fi
       else
           echo "- Active bean: ${BEAN_ID} (beans query timed out)"
       fi
   else
       if [ "$BRANCH" = "main" ] || [ "$BRANCH" = "master" ]; then
           echo "- No active bean (on ${BRANCH} branch)"
       elif [ -z "$BRANCH" ]; then
           echo "- No active bean"
       else
           echo "- No active bean (branch does not follow naming convention)"
       fi
   fi

   # --- Recent commits ---
   COMMITS=$(git log --oneline -5 2>/dev/null || echo "")
   if [ -n "$COMMITS" ]; then
       echo "- Recent commits:"
       echo "$COMMITS" | while IFS= read -r line; do
           echo "  - $line"
       done
   fi
   ```

   Key design decisions in this script:
   - **No `set -e`**: Unlike `get-current-bean.sh`, this script must never abort. Every external command uses `|| true`, `|| echo ""`, or `2>/dev/null`.
   - **Inline bean extraction**: Copies the regex logic from `get-current-bean.sh` (lines 17-20) rather than calling it, because that script uses `set -e` and `exit 1` on main/unrecognized branches.
   - **`timeout 2s`** on the beans query: Prevents the hook from blocking if beans CLI hangs.
   - **`--json` flag** on `beans graphql`: Without it, output contains ANSI color codes that are unparseable by `jq`.
   - **`jq` path is `.bean.title`** (not `.data.bean.title`): Verified against actual `beans graphql --json` output structure.
   - **Checklist counting**: Uses `grep -c '^\- \[[ x]\]'` for total items and `grep -c '^\- \[ \]'` for unchecked. The `^` anchor avoids matching checklist items inside code blocks at deeper indentation, though nested lists (e.g., `  - [ ]`) will not be counted — this is acceptable since top-level checklist items are the primary tracking mechanism.

2. **Make the script executable**

   ```bash
   chmod +x .claude/hooks/session-context.sh
   ```

3. **Update `.claude/settings.json` to register the hook**

   Add a new entry to the existing `SessionStart` array. The complete hooks section should become:

   ```json
   {
     "permissions": { ... },
     "hooks": {
       "PreToolUse": [ ... ],
       "SessionStart": [
         {
           "hooks": [
             {
               "type": "command",
               "command": "beans prime"
             }
           ]
         },
         {
           "hooks": [
             {
               "type": "command",
               "command": "$CLAUDE_PROJECT_DIR/.claude/hooks/session-context.sh"
             }
           ]
         }
       ],
       "PreCompact": [ ... ]
     }
   }
   ```

   Do NOT modify the `PreCompact` section — `SessionStart` already fires after compaction.

4. **Manual verification**

   - Start a new Claude Code session and verify the context block appears
   - Switch to a feature branch and restart to verify bean details appear
   - Test on main branch to verify the "no active bean" fallback
   - Time the script: `time .claude/hooks/session-context.sh` (must be under 2 seconds)

### Edge Cases to Handle

| Scenario | Expected behavior |
|---|---|
| On `main` branch | Output "No active bean (on main branch)" |
| On feature branch with valid bean | Output bean title, status, and checklist counts |
| On feature branch with non-existent bean ID | Output "Active bean: <id> (could not fetch details)" |
| Detached HEAD (during rebase) | Output "(detached HEAD)" and "No active bean" |
| `beans` CLI not installed | `timeout` or command failure caught, bean section degraded |
| `jq` not installed | `jq` errors caught by `2>/dev/null`, bean details skipped |
| Branch name does not follow convention | Output "No active bean (branch does not follow naming convention)" |
| Very large bean body (slow grep) | Negligible — grep on a string in memory is fast |

### Testing Strategy

- **Manual test on main**: Switch to main, run `.claude/hooks/session-context.sh`, verify output
- **Manual test on feature branch**: Create/switch to a feature branch, run script, verify bean details
- **Manual test with detached HEAD**: `git checkout HEAD~1`, run script, verify graceful handling
- **Performance test**: `time .claude/hooks/session-context.sh` — verify under 2 seconds
- **Integration test**: Start a new Claude Code session and verify the context appears in the system prompt
- No unit tests needed — this is a bash script with simple output; manual verification and the existing CI (lint/test) are sufficient

### Open Questions

None — the plan is complete and all ambiguities have been resolved through codebase analysis.

## Definition of Done
- [x] Hook script created at `.claude/hooks/session-context.sh`
- [x] Script is executable (`chmod +x`)
- [x] Script extracts branch, bean, and recent commits
- [x] Script handles edge cases: main branch, detached HEAD, missing bean, naming convention mismatch
- [x] Script never exits non-zero (graceful degradation)
- [x] Hook registered in `.claude/settings.json` for SessionStart (alongside existing `beans prime` hook)
- [x] Verified context appears on session start, resume, and after compaction
- [x] Script runs in under 2 seconds
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] Branch pushed and PR created for human review
- [x] Automated code review passed via `@review-backend` and/or `@review-frontend` subagents (via Task tool)
