---
# credfolio2-h0uy
title: Verify pre-commit hook works in both local and devcontainer environments
status: in-progress
type: task
priority: normal
created_at: 2026-02-07T16:15:18Z
updated_at: 2026-02-08T08:18:57Z
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

## Implementation Plan

### Approach

Audit the existing `pre-commit-check.sh` script for correctness and portability, fix discovered issues, add a unit test script (following the pattern of the two existing hook test files), and perform end-to-end verification in the devcontainer. The fixes are surgical — no architectural changes, just hardening an existing script.

### Issues Found During Codebase Analysis

The following concrete issues were identified by reading the script and testing its regex matching:

#### Issue 1: `set -e` causes silent pass-through on infrastructure errors (CRITICAL)

Line 11 has `set -e`. If `jq` is missing or the JSON input is malformed, `jq` on line 17 exits non-zero, and `set -e` terminates the script with exit code 1. Claude Code PreToolUse hooks only **block** on exit code 2 — exit code 1 is treated as a hook error and the tool call **proceeds**. This means a broken hook silently allows commits.

**Fix**: Replace `set -e` with explicit error handling. Wrap the `jq` call so that if it fails, the script exits 0 (allow) rather than 1, or better yet, default `COMMAND` to empty string on jq failure and let the regex mismatch path handle it.

#### Issue 2: Regex false positives on string literals and comments

The regex on line 20 matches `git commit` anywhere in the command string. Testing confirms it matches:
- `echo "git commit"` — a string literal, not an actual commit
- `# git commit -m "test"` — a comment

These false positives cause unnecessary lint+test runs. They are not dangerous (they do not skip checks) but waste time and could confuse the agent.

**Fix**: Tighten the regex to only match `git commit` at the start of a command (after optional `&&`, `;`, or start-of-string), not inside quoted strings. Follow the approach used by `validate-branch-name.sh` (lines 24-26) which extracts just the first line and uses more targeted matching.

#### Issue 3: `npm exec -- pnpm` indirection is unnecessary

Lines 31 and 44 use `npm exec -- pnpm lint` and `npm exec -- pnpm test`. In both environments (devcontainer and local), `pnpm` is directly on PATH (via mise shims in devcontainer, via corepack or direct install locally). The `npm exec` wrapper adds overhead (~2-3s) and an extra failure point.

**Fix**: Replace `npm exec -- pnpm lint` with `pnpm lint` and `npm exec -- pnpm test` with `pnpm test`.

#### Issue 4: No timeout on lint/test execution

If `pnpm lint` or `pnpm test` hangs (e.g., due to a deadlock in tests or Turborepo cache corruption), the hook blocks indefinitely. The `timeout` command is available in the devcontainer (it is part of coreutils, installed via Debian base image).

**Fix**: Wrap both commands in `timeout`. A 5-minute timeout (`timeout 300`) is reasonable given that turbo's `test` task has `dependsOn: ["build"]`, which means a cold run includes a full build.

#### Issue 5: stdout/stderr mixing

Lines 31 and 44 redirect stderr to stdout (`2>&1`), but the hook communicates blocking decisions via stderr (lines 27, 30, etc.). The lint/test output goes to stdout which may not be visible to the user in Claude Code's PreToolUse context. All output should go to stderr for visibility.

**Fix**: Change `2>&1` to `>&2` so lint/test output goes to stderr alongside the hook's own messages.

#### Issue 6: Missing unit test script

The other two hooks (`validate-bean-dod.sh`, `validate-bean-completion.sh`) each have test scripts in `.claude/hooks/tests/`. The pre-commit hook has none. This makes it hard to verify regex behavior without manual testing.

**Fix**: Create `.claude/hooks/tests/test-pre-commit-check.sh` following the same pattern as the existing test files. Mock `pnpm` instead of `beans`. Test the regex matching (true positives, true negatives, false positive avoidance) and exit code behavior.

### Files to Create/Modify

- `.claude/hooks/pre-commit-check.sh` — Fix issues 1-5 (set -e, regex, npm exec, timeout, stderr)
- `.claude/hooks/tests/test-pre-commit-check.sh` — New file: unit tests for the pre-commit hook (issue 6)

### Steps

1. **Write the unit test script first (TDD)**
   - Create `.claude/hooks/tests/test-pre-commit-check.sh`
   - Follow the pattern from `.claude/hooks/tests/test-validate-bean-dod.sh`: temp dir for mocks, `run_hook` helper, `assert_exit_code`/`assert_stderr_contains`/`assert_stderr_empty` helpers
   - Mock `pnpm` command (not `npm exec -- pnpm`) to return success or failure
   - Mock `jq` to test the infrastructure-error case
   - Test cases for regex matching:
     - `git commit -m "test"` → should trigger checks
     - `git commit --amend` → should trigger checks
     - `git commit` (bare) → should trigger checks
     - `git add . && git commit -m "test"` → should trigger checks
     - `git -C /workspace commit -m "test"` → should trigger checks
     - `echo "git commit"` → should NOT trigger checks (false positive fix)
     - `git log --oneline` → should NOT trigger checks
     - `pnpm lint` → should NOT trigger checks
     - `git merge --commit` → should NOT trigger checks (--commit is a flag, not a subcommand)
   - Test cases for exit behavior:
     - Lint passes + tests pass → exit 0
     - Lint fails → exit 2, stderr contains "COMMIT BLOCKED"
     - Tests fail → exit 2, stderr contains "COMMIT BLOCKED"
     - jq not found / malformed JSON → exit 0 (not exit 1)
   - Make the test script executable: `chmod +x`

2. **Run the tests to confirm they fail** (red phase)
   - Run `bash .claude/hooks/tests/test-pre-commit-check.sh`
   - The regex false-positive tests and jq-failure test should fail against the current script

3. **Fix `set -e` and add error handling for jq** (issue 1)
   - Remove `set -e` from line 11
   - Change line 17 to: `COMMAND=$(echo "$INPUT" | jq -r '.tool_input.command // ""' 2>/dev/null || echo "")`
   - This way, if jq is missing or input is malformed, COMMAND defaults to empty string, the regex does not match, and the script exits 0 (allow)

4. **Tighten the regex to avoid false positives** (issue 2)
   - Replace the current regex on line 20 with a more precise pattern that matches `git commit` only as a command, not inside string literals
   - Approach: check that `git` appears at a command boundary (start of string, after `&&`, after `;`, or after `|`)
   - Suggested pattern:
     ```bash
     if ! echo "$COMMAND" | grep -qP '(^|&&\s*|;\s*|\|\s*)git\s+(-\S+\s+)*commit(\s|$)'; then
         exit 0
     fi
     ```
   - This uses `grep -P` (Perl regex, available in the devcontainer via grep) for consistency with `validate-branch-name.sh` which also uses `grep -qP`
   - The `(\s|$)` anchor prevents matching `git committer` or similar

5. **Replace `npm exec -- pnpm` with direct `pnpm`** (issue 3)
   - Line 31: change `npm exec -- pnpm lint 2>&1` to `pnpm lint >&2`
   - Line 44: change `npm exec -- pnpm test 2>&1` to `pnpm test >&2`
   - Also fixes issue 5 (stderr redirect) in the same change

6. **Add timeout to lint and test commands** (issue 4)
   - Line 31: `timeout 300 pnpm lint >&2`
   - Line 44: `timeout 300 pnpm test >&2`
   - Add a check for timeout exit code (124) with a distinct error message:
     ```bash
     lint_exit=$?
     if [ "$lint_exit" -eq 124 ]; then
         echo "COMMIT BLOCKED: Linter timed out after 5 minutes" >&2
         exit 2
     elif [ "$lint_exit" -ne 0 ]; then
         echo "COMMIT BLOCKED: Linter failed" >&2
         exit 2
     fi
     ```

7. **Run the unit tests again** (green phase)
   - Run `bash .claude/hooks/tests/test-pre-commit-check.sh`
   - All tests should pass

8. **End-to-end verification in devcontainer**
   - **Test lint failure blocks commit**: Create a temporary file with a lint error (e.g., add `var x = 1;` to a .ts file), stage it, then simulate the hook by piping JSON with a `git commit` command to the hook script and verifying exit code 2
   - **Test test failure blocks commit**: Temporarily break a test, stage it, run the hook, verify exit code 2
   - **Test clean code allows commit**: With no staged changes that would break lint/tests, run the hook with a `git commit` command and verify exit code 0
   - **Test hook output visibility**: Verify that stderr messages ("COMMIT BLOCKED", "All checks passed!") appear in the output

9. **Verify hook registration**
   - Confirm `.claude/settings.json` has the PreToolUse hook registered (it does — lines 12-19)
   - Confirm `.claude/settings.local.json` does NOT need the hook (it only has permissions overrides)
   - Confirm `~/.claude/settings.json` (which in the devcontainer is the `claude-user-settings.json` copied to `/home/node/.claude/settings.json`) does NOT need the hook
   - Document: hooks in `.claude/settings.json` are project-level and apply to all environments. The `bypassPermissions` mode in the devcontainer only affects permission prompts, not hook execution.

10. **Run `pnpm lint` and `pnpm test` at the project root**
    - Verify the full project passes lint and tests with the modified hook script

11. **Push branch and create PR**

### Testing Strategy

- **Unit tests**: `.claude/hooks/tests/test-pre-commit-check.sh` — covers regex matching (true/false positives/negatives), exit code behavior (lint fail, test fail, jq fail, timeout), and output content
- **End-to-end tests**: Manual verification in devcontainer by piping JSON to the hook script and checking exit codes and stderr
- **Regression**: Run existing hook tests (`test-validate-bean-dod.sh`, `test-validate-bean-completion.sh`) to verify no accidental breakage

### Open Questions

- **Turbo caching and speed**: `pnpm test` triggers `dependsOn: ["build"]` in turbo.json, which means cold runs include a full build. With Turbo caching this is usually fast, but consider whether the 5-minute timeout is sufficient for truly cold builds (no cache). If this proves too slow in practice, a follow-up bean could explore running only affected workspace tests.
- **`git merge --commit` false positive**: The current regex matches `git merge --commit` because `commit` appears as a substring. The tightened regex (step 4) fixes this by requiring `commit` to appear as a git subcommand (after `git` and optional flags), not as a flag value. Verify this with the unit tests.

## Definition of Done
- [x] Hook registration verified in both environments
- [x] Script reviewed for portability issues
- [x] Hook input parsing verified for all commit command variants
- [x] End-to-end test: lint failure blocks commit in devcontainer
- [x] End-to-end test: test failure blocks commit in devcontainer
- [x] End-to-end test: clean code allows commit in devcontainer
- [x] Any issues found are fixed
- [x] Failure modes documented or addressed (e.g., timeout added)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] Branch pushed and PR created for human review
