#!/bin/bash
# Tests for validate-bean-completion.sh TaskCompleted hook
#
# Usage: bash .claude/hooks/tests/test-validate-bean-completion.sh
#
# These tests use mock `beans` and `git` commands to avoid requiring real infrastructure.

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
HOOK_SCRIPT="$SCRIPT_DIR/../validate-bean-completion.sh"

# Track test results
TESTS_RUN=0
TESTS_PASSED=0
TESTS_FAILED=0

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

# Create temp directory for mock commands and test artifacts
MOCK_DIR=$(mktemp -d)
export PATH="$MOCK_DIR:$PATH"

# Set CLAUDE_PROJECT_DIR for the hook
export CLAUDE_PROJECT_DIR="$(cd "$SCRIPT_DIR/../.." && cd .. && pwd)"

cleanup() {
    rm -rf "$MOCK_DIR"
}
trap cleanup EXIT

# Run the hook and capture stderr (where the blocking messages go)
# Returns exit code; stdout and stderr captured separately
run_hook() {
    local exit_code=0
    bash "$HOOK_SCRIPT" 2>"$MOCK_DIR/stderr_output" || exit_code=$?
    echo "$exit_code"
}

get_stderr() {
    # Filter out locale warnings that appear in devcontainer environments
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

# Helper: create TaskCompleted JSON input
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

# Helper: create mock beans command that returns specific body content
setup_mock_beans() {
    local bean_body="$1"
    local response_file="$MOCK_DIR/mock_response.json"
    jq -n --arg body "$bean_body" \
        '{"bean":{"id":"credfolio2-test","title":"Test Bean","status":"in-progress","body":$body}}' > "$response_file"

    cat > "$MOCK_DIR/beans" <<MOCKEOF
#!/bin/bash
if [[ "\$1" == "query" ]]; then
    cat "$response_file"
fi
MOCKEOF
    chmod +x "$MOCK_DIR/beans"
}

# Helper: create mock beans command that fails (simulates timeout/error)
setup_mock_beans_fail() {
    cat > "$MOCK_DIR/beans" <<'MOCKEOF'
#!/bin/bash
exit 1
MOCKEOF
    chmod +x "$MOCK_DIR/beans"
}

# Helper: create mock git command that returns a specific branch name
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

# Helper: remove mock git (so real git is used, or no git available)
remove_mock_git() {
    rm -f "$MOCK_DIR/git"
}

# Clean up mocks between tests
reset_mocks() {
    rm -f "$MOCK_DIR/beans" "$MOCK_DIR/git" "$MOCK_DIR/mock_response.json" "$MOCK_DIR/stderr_output"
}


# ============================================================
echo "=== Test: Bean ID in task_subject, unchecked items ==="
# ============================================================

reset_mocks
BODY="## Checklist
- [x] Write the code
- [ ] Run tests
- [ ] Update docs

## Definition of Done
- [x] Tests written
- [ ] pnpm lint passes
- [ ] pnpm test passes"

setup_mock_beans "$BODY"

EXIT_CODE=$(make_input "task-1" "Implement bean credfolio2-abc1" "" | run_hook)
assert_exit_code "blocks completion (exit 2)" 2 "$EXIT_CODE"
assert_stderr_contains "mentions bean ID" "credfolio2-abc1"
assert_stderr_contains "mentions BLOCKED" "BLOCKED"
assert_stderr_contains "lists unchecked item: Run tests" "Run tests"
assert_stderr_contains "lists unchecked item: pnpm lint" "pnpm lint passes"


# ============================================================
echo ""
echo "=== Test: Bean ID in task_subject, all items checked ==="
# ============================================================

reset_mocks
BODY="## Checklist
- [x] Write the code
- [x] Run tests
- [x] Update docs

## Definition of Done
- [x] Tests written
- [x] pnpm lint passes
- [x] pnpm test passes"

setup_mock_beans "$BODY"

EXIT_CODE=$(make_input "task-2" "Implement bean credfolio2-def2" "" | run_hook)
assert_exit_code "allows completion (exit 0)" 0 "$EXIT_CODE"
assert_stderr_empty "no stderr output when all checked"


# ============================================================
echo ""
echo "=== Test: Bean ID in task_description only ==="
# ============================================================

reset_mocks
BODY="## Checklist
- [ ] Remaining work"

setup_mock_beans "$BODY"

EXIT_CODE=$(make_input "task-3" "Do some work" "Working on credfolio2-ghi3 implementation" | run_hook)
assert_exit_code "finds bean ID in description, blocks (exit 2)" 2 "$EXIT_CODE"
assert_stderr_contains "found bean from description" "credfolio2-ghi3"


# ============================================================
echo ""
echo "=== Test: Bean ID from git branch fallback ==="
# ============================================================

reset_mocks
BODY="## Checklist
- [ ] Still todo"

setup_mock_beans "$BODY"
setup_mock_git "feat/credfolio2-jkl4-some-feature"

EXIT_CODE=$(make_input "task-4" "Do some task" "No bean reference here" | run_hook)
assert_exit_code "finds bean via git branch, blocks (exit 2)" 2 "$EXIT_CODE"
assert_stderr_contains "found bean from branch" "credfolio2-jkl4"


# ============================================================
echo ""
echo "=== Test: No bean ID anywhere (main branch, no reference) ==="
# ============================================================

reset_mocks
setup_mock_git "main"

EXIT_CODE=$(make_input "task-5" "Random task" "No bean context at all" | run_hook)
assert_exit_code "allows completion (exit 0)" 0 "$EXIT_CODE"
assert_stderr_empty "no stderr output for non-bean task"


# ============================================================
echo ""
echo "=== Test: Bean query fails/times out ==="
# ============================================================

reset_mocks
setup_mock_beans_fail

EXIT_CODE=$(make_input "task-6" "Implement credfolio2-mno5" "" | run_hook)
assert_exit_code "allows completion on query failure (exit 0)" 0 "$EXIT_CODE"
assert_stderr_empty "no stderr when query fails"


# ============================================================
echo ""
echo "=== Test: Bean has no checklist at all ==="
# ============================================================

reset_mocks
BODY="## Description
This bean is just a description with no checklist items whatsoever.

## Notes
Some additional context."

setup_mock_beans "$BODY"

EXIT_CODE=$(make_input "task-7" "Work on credfolio2-pqr6" "" | run_hook)
assert_exit_code "allows completion with no checklist (exit 0)" 0 "$EXIT_CODE"
assert_stderr_empty "no stderr when no checklist"


# ============================================================
echo ""
echo "=== Test: Bean has mixed checked/unchecked items ==="
# ============================================================

reset_mocks
BODY="## Checklist
- [x] First item done
- [ ] Second item pending
- [x] Third item done
- [ ] Fourth item pending
- [x] Fifth item done"

setup_mock_beans "$BODY"

EXIT_CODE=$(make_input "task-8" "credfolio2-stu7 mixed items" "" | run_hook)
assert_exit_code "blocks completion with mixed items (exit 2)" 2 "$EXIT_CODE"
assert_stderr_contains "counts unchecked items" "2 unchecked"
assert_stderr_contains "lists Second item" "Second item pending"
assert_stderr_contains "lists Fourth item" "Fourth item pending"


# ============================================================
echo ""
echo "=== Test: Non-bean task (no bean context, no branch) ==="
# ============================================================

reset_mocks
# No mock git set up — remove any existing one
remove_mock_git
# Use a real git that will return whatever branch we're actually on
# But the task has no bean reference, so strategy 1 fails
# Strategy 2 might or might not find a bean depending on current branch
# To make this deterministic, mock git to return a non-matching branch
setup_mock_git "some-random-branch"

EXIT_CODE=$(make_input "task-9" "General maintenance task" "No beans involved" | run_hook)
assert_exit_code "allows completion for non-bean task (exit 0)" 0 "$EXIT_CODE"
assert_stderr_empty "no stderr for non-bean task"


# ============================================================
echo ""
echo "==============================="
echo "Results: $TESTS_PASSED/$TESTS_RUN passed, $TESTS_FAILED failed"
echo "==============================="

if [ "$TESTS_FAILED" -gt 0 ]; then
    exit 1
fi
exit 0
