#!/bin/bash
# Tests for pre-commit-check.sh PreToolUse hook
#
# Usage: bash .claude/hooks/tests/test-pre-commit-check.sh
#
# These tests use mock `pnpm` and `jq` commands to avoid requiring real build infrastructure.

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
HOOK_SCRIPT="$SCRIPT_DIR/../pre-commit-check.sh"

# Track test results
TESTS_RUN=0
TESTS_PASSED=0
TESTS_FAILED=0

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

# Create temp directory for mock commands
MOCK_DIR=$(mktemp -d)
export PATH="$MOCK_DIR:$PATH"

# Set CLAUDE_PROJECT_DIR for the hook
export CLAUDE_PROJECT_DIR="$(cd "$SCRIPT_DIR/../.." && cd .. && pwd)"

cleanup() {
    rm -rf "$MOCK_DIR"
}
trap cleanup EXIT

# Run the hook and capture stderr
# Returns exit code via echo; stderr stored in file
run_hook() {
    local exit_code=0
    bash "$HOOK_SCRIPT" 2>"$MOCK_DIR/stderr_output" || exit_code=$?
    echo "$exit_code"
}

get_stderr() {
    cat "$MOCK_DIR/stderr_output" 2>/dev/null | grep -v 'setlocale: LC_ALL: cannot change locale' || echo ""
}

assert_exit_code() {
    local test_name="$1"
    local expected_code="$2"
    local actual_code="$3"
    TESTS_RUN=$((TESTS_RUN + 1))
    if [ "$actual_code" -eq "$expected_code" ]; then
        TESTS_PASSED=$((TESTS_PASSED + 1))
        echo -e "  ${GREEN}PASS${NC}: $test_name"
    else
        TESTS_FAILED=$((TESTS_FAILED + 1))
        echo -e "  ${RED}FAIL${NC}: $test_name (expected exit $expected_code, got $actual_code)"
    fi
}

assert_stderr_contains() {
    local test_name="$1"
    local expected="$2"
    local actual
    actual=$(get_stderr)
    TESTS_RUN=$((TESTS_RUN + 1))
    if echo "$actual" | grep -qF "$expected"; then
        TESTS_PASSED=$((TESTS_PASSED + 1))
        echo -e "  ${GREEN}PASS${NC}: $test_name"
    else
        TESTS_FAILED=$((TESTS_FAILED + 1))
        echo -e "  ${RED}FAIL${NC}: $test_name (expected stderr to contain '$expected')"
        echo "    Actual stderr: $actual"
    fi
}

assert_stderr_empty() {
    local test_name="$1"
    local actual
    actual=$(get_stderr)
    TESTS_RUN=$((TESTS_RUN + 1))
    if [ -z "$actual" ]; then
        TESTS_PASSED=$((TESTS_PASSED + 1))
        echo -e "  ${GREEN}PASS${NC}: $test_name"
    else
        TESTS_FAILED=$((TESTS_FAILED + 1))
        echo -e "  ${RED}FAIL${NC}: $test_name (expected empty stderr, got '$actual')"
    fi
}

# Helper: create PreToolUse JSON input
make_input() {
    local command="$1"
    jq -n \
        --arg cmd "$command" \
        '{
            session_id: "test-session",
            transcript_path: "/tmp/transcript.jsonl",
            cwd: "/workspace",
            permission_mode: "default",
            hook_event_name: "PreToolUse",
            tool_name: "Bash",
            tool_input: {command: $cmd},
            tool_use_id: "toolu_test123"
        }'
}

# Helper: create mock pnpm that succeeds
setup_mock_pnpm_success() {
    cat > "$MOCK_DIR/pnpm" <<'MOCKEOF'
#!/bin/bash
exit 0
MOCKEOF
    chmod +x "$MOCK_DIR/pnpm"
}

# Helper: create mock pnpm where lint fails
setup_mock_pnpm_lint_fails() {
    cat > "$MOCK_DIR/pnpm" <<'MOCKEOF'
#!/bin/bash
if [[ "$1" == "lint" ]]; then
    echo "Error: Lint failed" >&2
    exit 1
fi
exit 0
MOCKEOF
    chmod +x "$MOCK_DIR/pnpm"
}

# Helper: create mock pnpm where test fails
setup_mock_pnpm_test_fails() {
    cat > "$MOCK_DIR/pnpm" <<'MOCKEOF'
#!/bin/bash
if [[ "$1" == "test" ]]; then
    echo "Error: Tests failed" >&2
    exit 1
fi
exit 0
MOCKEOF
    chmod +x "$MOCK_DIR/pnpm"
}

# Helper: remove mock jq (to simulate jq not found)
setup_missing_jq() {
    # Create a mock jq that fails
    cat > "$MOCK_DIR/jq" <<'MOCKEOF'
#!/bin/bash
exit 1
MOCKEOF
    chmod +x "$MOCK_DIR/jq"
}

# Helper: restore jq by removing mock
restore_jq() {
    rm -f "$MOCK_DIR/jq"
}

# Clean up mocks between tests
reset_mocks() {
    rm -f "$MOCK_DIR/pnpm" "$MOCK_DIR/jq" "$MOCK_DIR/timeout" "$MOCK_DIR/stderr_output"
}


# ============================================================
echo "=== Test: Regex matching - true positives ==="
echo "=== (commands that SHOULD trigger checks) ==="
# ============================================================

# For regex-only tests, mock pnpm to succeed so we can confirm the hook runs checks
setup_mock_pnpm_success

EXIT_CODE=$(make_input 'git commit -m "test"' | run_hook)
assert_exit_code "git commit -m triggers checks (exit 0 = checks ran and passed)" 0 "$EXIT_CODE"
assert_stderr_contains "git commit -m shows 'All checks passed'" "All checks passed"

reset_mocks
setup_mock_pnpm_success
EXIT_CODE=$(make_input 'git commit --amend' | run_hook)
assert_exit_code "git commit --amend triggers checks" 0 "$EXIT_CODE"
assert_stderr_contains "git commit --amend shows 'All checks passed'" "All checks passed"

reset_mocks
setup_mock_pnpm_success
EXIT_CODE=$(make_input 'git commit' | run_hook)
assert_exit_code "bare git commit triggers checks" 0 "$EXIT_CODE"
assert_stderr_contains "bare git commit shows 'All checks passed'" "All checks passed"

reset_mocks
setup_mock_pnpm_success
EXIT_CODE=$(make_input 'git add . && git commit -m "test"' | run_hook)
assert_exit_code "git add && git commit triggers checks" 0 "$EXIT_CODE"
assert_stderr_contains "git add && git commit shows 'All checks passed'" "All checks passed"

reset_mocks
setup_mock_pnpm_success
EXIT_CODE=$(make_input 'git -C /workspace commit -m "test"' | run_hook)
assert_exit_code "git -C <dir> commit triggers checks" 0 "$EXIT_CODE"
assert_stderr_contains "git -C <dir> commit shows 'All checks passed'" "All checks passed"


# ============================================================
echo ""
echo "=== Test: Regex matching - true negatives ==="
echo "=== (commands that should NOT trigger checks) ==="
# ============================================================

reset_mocks
EXIT_CODE=$(make_input 'git log --oneline' | run_hook)
assert_exit_code "git log exits 0" 0 "$EXIT_CODE"
assert_stderr_empty "git log produces no stderr"

reset_mocks
EXIT_CODE=$(make_input 'pnpm lint' | run_hook)
assert_exit_code "pnpm lint exits 0" 0 "$EXIT_CODE"
assert_stderr_empty "pnpm lint produces no stderr"

reset_mocks
EXIT_CODE=$(make_input 'git status' | run_hook)
assert_exit_code "git status exits 0" 0 "$EXIT_CODE"
assert_stderr_empty "git status produces no stderr"

reset_mocks
EXIT_CODE=$(make_input 'git push origin main' | run_hook)
assert_exit_code "git push exits 0" 0 "$EXIT_CODE"
assert_stderr_empty "git push produces no stderr"


# ============================================================
echo ""
echo "=== Test: Regex matching - false positive avoidance ==="
echo "=== (commands that look like commit but are not) ==="
# ============================================================

reset_mocks
EXIT_CODE=$(make_input 'echo "git commit"' | run_hook)
assert_exit_code "echo git commit is NOT a commit (exit 0)" 0 "$EXIT_CODE"
assert_stderr_empty "echo git commit produces no stderr"

reset_mocks
EXIT_CODE=$(make_input 'git merge --commit' | run_hook)
assert_exit_code "git merge --commit is NOT a commit (exit 0)" 0 "$EXIT_CODE"
assert_stderr_empty "git merge --commit produces no stderr"


# ============================================================
echo ""
echo "=== Test: Exit behavior - lint failure ==="
# ============================================================

reset_mocks
setup_mock_pnpm_lint_fails
EXIT_CODE=$(make_input 'git commit -m "test"' | run_hook)
assert_exit_code "lint failure blocks commit (exit 2)" 2 "$EXIT_CODE"
assert_stderr_contains "lint failure mentions COMMIT BLOCKED" "COMMIT BLOCKED"
assert_stderr_contains "lint failure mentions Linter failed" "Linter failed"


# ============================================================
echo ""
echo "=== Test: Exit behavior - test failure ==="
# ============================================================

reset_mocks
setup_mock_pnpm_test_fails
EXIT_CODE=$(make_input 'git commit -m "test"' | run_hook)
assert_exit_code "test failure blocks commit (exit 2)" 2 "$EXIT_CODE"
assert_stderr_contains "test failure mentions COMMIT BLOCKED" "COMMIT BLOCKED"
assert_stderr_contains "test failure mentions Tests failed" "Tests failed"


# ============================================================
echo ""
echo "=== Test: Exit behavior - both lint and test pass ==="
# ============================================================

reset_mocks
setup_mock_pnpm_success
EXIT_CODE=$(make_input 'git commit -m "test"' | run_hook)
assert_exit_code "all checks pass allows commit (exit 0)" 0 "$EXIT_CODE"
assert_stderr_contains "all checks pass shows success message" "All checks passed"


# ============================================================
echo ""
echo "=== Test: Exit behavior - jq failure (infrastructure error) ==="
# ============================================================

reset_mocks
setup_missing_jq
EXIT_CODE=$(make_input 'git commit -m "test"' | run_hook)
assert_exit_code "jq failure allows commit (exit 0, not exit 1)" 0 "$EXIT_CODE"
# When jq fails, the hook should gracefully exit 0 (not exit 1 from set -e)
restore_jq


# ============================================================
echo ""
echo "=== Test: Exit behavior - pnpm called directly (not via npm exec) ==="
# ============================================================

# This test verifies that the hook calls `pnpm` directly, not `npm exec -- pnpm`.
# We create a mock pnpm that records it was called, and no mock npm.
# If the hook uses `npm exec -- pnpm`, it would call the real npm (or fail),
# not our mock pnpm.
reset_mocks
cat > "$MOCK_DIR/pnpm" <<'MOCKEOF'
#!/bin/bash
echo "MOCK_PNPM_CALLED" >> /tmp/test-pre-commit-pnpm-calls.txt
exit 0
MOCKEOF
chmod +x "$MOCK_DIR/pnpm"
rm -f /tmp/test-pre-commit-pnpm-calls.txt

EXIT_CODE=$(make_input 'git commit -m "test"' | run_hook)
PNPM_CALLS=""
if [ -f /tmp/test-pre-commit-pnpm-calls.txt ]; then
    PNPM_CALLS=$(cat /tmp/test-pre-commit-pnpm-calls.txt)
fi
rm -f /tmp/test-pre-commit-pnpm-calls.txt

TESTS_RUN=$((TESTS_RUN + 1))
if echo "$PNPM_CALLS" | grep -q "MOCK_PNPM_CALLED"; then
    TESTS_PASSED=$((TESTS_PASSED + 1))
    echo -e "  ${GREEN}PASS${NC}: pnpm called directly (not via npm exec)"
else
    TESTS_FAILED=$((TESTS_FAILED + 1))
    echo -e "  ${RED}FAIL${NC}: pnpm called directly (not via npm exec) - pnpm mock was not invoked"
fi


# ============================================================
echo ""
echo "=== Test: Output goes to stderr (not stdout) ==="
# ============================================================

reset_mocks
setup_mock_pnpm_lint_fails
# Capture stdout separately from stderr
STDOUT_OUTPUT=""
EXIT_CODE=0
STDOUT_OUTPUT=$(make_input 'git commit -m "test"' | bash "$HOOK_SCRIPT" 2>"$MOCK_DIR/stderr_output") || EXIT_CODE=$?

TESTS_RUN=$((TESTS_RUN + 1))
if [ -z "$STDOUT_OUTPUT" ]; then
    TESTS_PASSED=$((TESTS_PASSED + 1))
    echo -e "  ${GREEN}PASS${NC}: hook produces no stdout (all output goes to stderr)"
else
    TESTS_FAILED=$((TESTS_FAILED + 1))
    echo -e "  ${RED}FAIL${NC}: hook produces no stdout (all output goes to stderr) - got stdout: '$STDOUT_OUTPUT'"
fi


# ============================================================
echo ""
echo "==============================="
echo "Results: $TESTS_PASSED/$TESTS_RUN passed, $TESTS_FAILED failed"
echo "==============================="

if [ "$TESTS_FAILED" -gt 0 ]; then
    exit 1
fi
exit 0
