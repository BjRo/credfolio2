---
# credfolio2-tdjg
title: Add TaskCompleted hook to enforce bean checklist completion
status: in-progress
type: task
priority: normal
created_at: 2026-02-07T16:19:41Z
updated_at: 2026-02-07T21:11:13Z
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
  - [ ] Visual verification via `@qa` subagent (via Task tool)
  - [ ] Branch pushed and PR created
  - [ ] Automated code review passed
"
exit 2  # blocks completion

# Claude sees the feedback and works on remaining items
```

## Implementation Plan

### Approach

Create a bash hook script that uses a multi-strategy bean ID extraction approach: first try to extract a bean ID from the TaskCompleted hook's `task_subject` and `task_description` fields (which may contain the bean ID directly), then fall back to extracting it from the current git branch name (which follows the `<type>/<bean-id>-<description>` convention). Once the bean is identified, query its body and scan for unchecked checklist items (`- [ ]`). If any remain, exit with code 2 to block task completion and print the remaining items to stderr so the agent sees feedback.

The script follows the same patterns as the existing hooks (`validate-bean-dod.sh`, `pre-commit-check.sh`, `validate-branch-name.sh`): read JSON from stdin via `jq`, use `beans query` for bean data, and follow the exit code conventions (0=allow, 2=block with stderr feedback).

A comprehensive test script follows the pattern of the existing `test-validate-bean-dod.sh`.

### Files to Create/Modify

- `.claude/hooks/validate-bean-completion.sh` (CREATE) — The main hook script
- `.claude/hooks/tests/test-validate-bean-completion.sh` (CREATE) — Test script following the pattern of `test-validate-bean-dod.sh`
- `.claude/settings.json` (MODIFY) — Add `TaskCompleted` hook registration

### Steps

1. **Create the hook script at `.claude/hooks/validate-bean-completion.sh`**

   The script should:

   a. Read JSON from stdin into a variable (same as `validate-bean-dod.sh` line 12: `INPUT=$(cat)`).

   b. Extract fields from the JSON using `jq`:
      - `task_id` (string)
      - `task_subject` (string)
      - `task_description` (string, optional)

   c. **Bean ID extraction — Strategy 1: From hook input fields.** Search `task_subject` and `task_description` for a bean ID pattern (`credfolio2-[a-zA-Z0-9]+`). Use `grep -oP` to extract. This covers cases where the task was created with the bean ID in its subject (e.g., "Implement bean credfolio2-abc1" or "credfolio2-abc1 — Add user auth").

   d. **Bean ID extraction — Strategy 2: From git branch name.** If Strategy 1 finds nothing, extract from the current branch using the same regex as `get-current-bean.sh` and `session-context.sh`:
      ```bash
      BRANCH=$(git branch --show-current 2>/dev/null)
      if [[ "$BRANCH" =~ ^[a-z]+/(credfolio2-[a-zA-Z0-9]+)-.* ]]; then
          BEAN_ID="${BASH_REMATCH[1]}"
      fi
      ```

   e. **No bean found — exit 0.** If neither strategy produces a bean ID, allow completion silently. This handles non-bean tasks gracefully.

   f. **Query the bean body** using `beans query`:
      ```bash
      BEAN_JSON=$(timeout 5s beans query "{ bean(id: \"$BEAN_ID\") { id title status body } }" --json 2>/dev/null)
      ```
      Use `timeout 5s` to prevent hanging (same defensive pattern as `session-context.sh` line 30). If the query fails or returns empty, exit 0 (don't block on infrastructure failures).

   g. **Extract the bean body** from the JSON response:
      ```bash
      BEAN_BODY=$(echo "$BEAN_JSON" | jq -r '.bean.body // empty')
      ```

   h. **Scan for unchecked items.** Use `grep` to find lines matching `^- \[ \]`:
      ```bash
      UNCHECKED=$(echo "$BEAN_BODY" | grep -P '^\- \[ \] ' || true)
      ```
      If no unchecked items found (empty result), exit 0.

   i. **Count unchecked items and format error message.** If unchecked items exist:
      ```bash
      COUNT=$(echo "$UNCHECKED" | wc -l)
      echo "" >&2
      echo "========================================" >&2
      echo "TASK COMPLETION BLOCKED" >&2
      echo "========================================" >&2
      echo "" >&2
      echo "Bean $BEAN_ID has $COUNT unchecked checklist item(s):" >&2
      echo "$UNCHECKED" >&2
      echo "" >&2
      echo "Complete all checklist items before finishing this task." >&2
      exit 2
      ```

   j. **Make the script executable**: `chmod +x .claude/hooks/validate-bean-completion.sh`

   Key design decisions:
   - Do NOT use `set -e` (a failing `grep` or `git` should not abort the script — we want graceful fallback to exit 0)
   - Use `timeout` on `beans query` to avoid hanging
   - Check both `task_subject` AND `task_description` for bean ID (the bean might be referenced in either)
   - Use `grep -P '^\- \[ \] '` with the leading space after `]` to avoid false positives on non-checklist content

2. **Register the hook in `.claude/settings.json`**

   Add a `TaskCompleted` entry to the existing `hooks` object. The current file has `PreToolUse`, `SessionStart`, `PostToolUse`, and `PreCompact`. Add after `PreCompact`:

   ```json
   "TaskCompleted": [
     {
       "hooks": [
         {
           "type": "command",
           "command": "$CLAUDE_PROJECT_DIR/.claude/hooks/validate-bean-completion.sh"
         }
       ]
     }
   ]
   ```

   Note: TaskCompleted does NOT support matchers (confirmed from source code analysis — it fires on all task completions, and matcher matching is skipped for this event type). So no `matcher` field is needed.

3. **Create the test script at `.claude/hooks/tests/test-validate-bean-completion.sh`**

   Follow the exact pattern of `test-validate-bean-dod.sh`:

   a. **Test infrastructure**: Same `assert_exit_code`, `assert_output_empty`, `assert_output_contains` helpers. Same mock `beans` command approach using `$MOCK_DIR` on `PATH`. Same `CLAUDE_PROJECT_DIR` setup.

   b. **Helper: `make_input`** — Create TaskCompleted JSON input:
      ```bash
      make_input() {
          local task_id="$1"
          local task_subject="$2"
          local task_description="${3:-}"
          jq -n \
              --arg id "$task_id" \
              --arg subject "$task_subject" \
              --arg desc "$task_description" \
              '{
                  session_id: "test-session",
                  transcript_path: "/tmp/transcript.jsonl",
                  cwd: "/workspace",
                  hook_event_name: "TaskCompleted",
                  task_id: $id,
                  task_subject: $subject,
                  task_description: $desc
              }'
      }
      ```

   c. **Helper: `setup_mock_beans`** — Mock beans CLI to return specific body content (same pattern as existing test).

   d. **Helper: `setup_mock_git`** — Mock `git` to return a specific branch name:
      ```bash
      setup_mock_git() {
          local branch="$1"
          cat > "$MOCK_DIR/git" <<MOCKEOF
      #!/bin/bash
      if [[ "\$1" == "branch" ]] && [[ "\$2" == "--show-current" ]]; then
          echo "$branch"
      fi
      MOCKEOF
          chmod +x "$MOCK_DIR/git"
      }
      ```

   e. **Test cases**:
      - **Bean ID in task_subject, unchecked items** — Should block (exit 2), stderr lists items
      - **Bean ID in task_subject, all items checked** — Should allow (exit 0)
      - **Bean ID in task_description only** — Should find it and check
      - **Bean ID in neither field, found via git branch** — Should find via branch fallback
      - **No bean ID anywhere (main branch, no bean reference)** — Should allow (exit 0)
      - **Bean query fails/times out** — Should allow (exit 0, graceful degradation)
      - **Bean has no checklist at all** — Should allow (exit 0)
      - **Bean has mixed checked/unchecked items** — Should block, listing only unchecked
      - **Non-bean task (no bean context)** — Should allow (exit 0)

4. **Run the test script to verify**

   ```bash
   bash /workspace/.claude/hooks/tests/test-validate-bean-completion.sh
   ```

### Testing Strategy

- **Automated**: The test script (`test-validate-bean-completion.sh`) uses mocked `beans` and `git` commands to test all scenarios without requiring real bean infrastructure.
- **Manual verification**: After implementation, test against a real bean:
  1. Create a test bean with unchecked items
  2. Be on a feature branch for that bean
  3. Verify the hook blocks completion (simulate by piping TaskCompleted JSON to the script)
  4. Check all items in the bean
  5. Verify the hook allows completion
- **Integration**: The `pnpm lint` and `pnpm test` commands should pass (this is a shell script + JSON config change, so no frontend/backend code is affected).

### Open Questions

None — the approach is well-defined in the bean and confirmed by source code analysis of Claude Code's hook system. The TaskCompleted hook input schema is: `{ session_id, transcript_path, cwd, permission_mode?, hook_event_name: "TaskCompleted", task_id, task_subject, task_description?, teammate_name?, team_name? }`. Exit code 2 blocks completion and shows stderr to the model. Exit code 0 allows.

## Definition of Done
- [x] Hook script created at `.claude/hooks/validate-bean-completion.sh`
- [x] Script reads bean body and detects unchecked checklist items
- [x] Hook registered in `.claude/settings.json`
- [x] Tested: bean with unchecked items → completion blocked with helpful message
- [x] Tested: bean with all items checked → completion allowed
- [x] Tested: non-bean task → completion allowed (graceful fallback)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Branch pushed and PR created for human review
