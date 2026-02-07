#!/bin/bash
# Tests for validate-bean-dod.sh PostToolUse hook
#
# Usage: bash .claude/hooks/tests/test-validate-bean-dod.sh
#
# These tests use a mock `beans` command to avoid requiring real bean infrastructure.

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
HOOK_SCRIPT="$SCRIPT_DIR/../validate-bean-dod.sh"
TEMPLATE_FILE="$SCRIPT_DIR/../../templates/definition-of-done.md"

# Track test results
TESTS_RUN=0
TESTS_PASSED=0
TESTS_FAILED=0

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

# Create temp directory for mock beans command and test artifacts
MOCK_DIR=$(mktemp -d)
MOCK_BEANS="$MOCK_DIR/beans"
export PATH="$MOCK_DIR:$PATH"

# Set CLAUDE_PROJECT_DIR for the hook to find the template
# From hooks/tests -> hooks -> .claude -> project root
export CLAUDE_PROJECT_DIR="$(cd "$SCRIPT_DIR/../.." && cd .. && pwd)"

cleanup() {
    rm -rf "$MOCK_DIR"
}
trap cleanup EXIT

# Run the hook and capture only stdout (discard stderr to avoid locale warnings etc.)
run_hook() {
    bash "$HOOK_SCRIPT" 2>/dev/null
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

assert_output_contains() {
    local test_name="$1"
    local expected="$2"
    local actual="$3"
    TESTS_RUN=$((TESTS_RUN + 1))
    if echo "$actual" | grep -qF "$expected"; then
        TESTS_PASSED=$((TESTS_PASSED + 1))
        echo -e "  ${GREEN}PASS${NC}: $test_name"
    else
        TESTS_FAILED=$((TESTS_FAILED + 1))
        echo -e "  ${RED}FAIL${NC}: $test_name (expected output to contain '$expected')"
        echo "    Actual output: $actual"
    fi
}

assert_output_empty() {
    local test_name="$1"
    local actual="$2"
    TESTS_RUN=$((TESTS_RUN + 1))
    if [ -z "$actual" ]; then
        TESTS_PASSED=$((TESTS_PASSED + 1))
        echo -e "  ${GREEN}PASS${NC}: $test_name"
    else
        TESTS_FAILED=$((TESTS_FAILED + 1))
        echo -e "  ${RED}FAIL${NC}: $test_name (expected empty output, got '$actual')"
    fi
}

# Helper: create mock beans that returns a specific type and body
setup_mock_beans() {
    local bean_type="$1"
    local bean_body="$2"
    # Write the JSON response to a file to avoid escaping issues
    local response_file="$MOCK_DIR/mock_response.json"
    # Use jq to construct valid JSON
    jq -n --arg type "$bean_type" --arg body "$bean_body" \
        '{"bean":{"type":$type,"body":$body}}' > "$response_file"

    cat > "$MOCK_BEANS" <<MOCKEOF
#!/bin/bash
if [[ "\$1" == "query" ]]; then
    cat "$response_file"
fi
MOCKEOF
    chmod +x "$MOCK_BEANS"
}

# Helper: create PostToolUse JSON input using jq for proper escaping
make_input() {
    local command="$1"
    local tool_response="$2"
    jq -n \
        --arg cmd "$command" \
        --arg resp "$tool_response" \
        '{
            session_id: "test-session",
            transcript_path: "/tmp/transcript.jsonl",
            cwd: "/workspace",
            permission_mode: "default",
            hook_event_name: "PostToolUse",
            tool_name: "Bash",
            tool_input: {command: $cmd},
            tool_response: $resp,
            tool_use_id: "toolu_test123"
        }'
}

DOD_BODY="## Definition of Done
- [ ] Tests written (TDD: write tests before implementation)
- [ ] \`pnpm lint\` passes with no errors
- [ ] \`pnpm test\` passes with no failures
- [ ] Visual verification via \`@qa\` subagent (via Task tool, for UI changes)
- [ ] ADR written via \`/decision\` skill (if new dependencies, patterns, or architectural changes were introduced)
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review
- [ ] Automated code review passed via \`@review-backend\` and/or \`@review-frontend\` subagents (via Task tool)"

echo "=== Test: Non-beans-create commands ==="

# Test: ls command should be ignored
OUTPUT=$(make_input "ls -la" "" | run_hook)
assert_exit_code "ls command exits 0" 0 $?
assert_output_empty "ls command produces no output" "$OUTPUT"

# Test: git commit should be ignored
OUTPUT=$(make_input "git commit -m 'test'" "" | run_hook)
assert_exit_code "git commit exits 0" 0 $?
assert_output_empty "git commit produces no output" "$OUTPUT"

# Test: beans query should be ignored
OUTPUT=$(make_input "beans query '{ beans { id } }'" "" | run_hook)
assert_exit_code "beans query exits 0" 0 $?
assert_output_empty "beans query produces no output" "$OUTPUT"

# Test: beans update should be ignored
OUTPUT=$(make_input "beans update credfolio2-abc1 --status completed" "" | run_hook)
assert_exit_code "beans update exits 0" 0 $?
assert_output_empty "beans update produces no output" "$OUTPUT"

echo ""
echo "=== Test: Epic creation (should skip validation) ==="

setup_mock_beans "epic" "An epic for grouping features"
OUTPUT=$(make_input 'beans create "My Epic" -t epic -d "Epic desc"' "Created credfolio2-abc1 credfolio2-abc1--my-epic.md" | run_hook)
assert_exit_code "epic creation exits 0" 0 $?
assert_output_empty "epic creation produces no output" "$OUTPUT"

echo ""
echo "=== Test: Milestone creation (should skip validation) ==="

setup_mock_beans "milestone" "A milestone for release"
OUTPUT=$(make_input 'beans create "Release v1" -t milestone -d "Release desc"' "Created credfolio2-xyz9 credfolio2-xyz9--release-v1.md" | run_hook)
assert_exit_code "milestone creation exits 0" 0 $?
assert_output_empty "milestone creation produces no output" "$OUTPUT"

echo ""
echo "=== Test: Task creation without DoD ==="

setup_mock_beans "task" "A simple task description without DoD"
OUTPUT=$(make_input 'beans create "Fix something" -t task -d "A simple task"' "Created credfolio2-def2 credfolio2-def2--fix-something.md" | run_hook)
assert_exit_code "task without DoD exits 0" 0 $?
assert_output_contains "task without DoD returns block decision" '"decision"' "$OUTPUT"
assert_output_contains "task without DoD mentions bean ID" 'credfolio2-def2' "$OUTPUT"
assert_output_contains "task without DoD mentions template file" 'definition-of-done.md' "$OUTPUT"

echo ""
echo "=== Test: Bug creation without DoD ==="

setup_mock_beans "bug" "A bug report"
OUTPUT=$(make_input 'beans create "Login broken" -t bug -d "Login is broken"' "Created credfolio2-ghi3 credfolio2-ghi3--login-broken.md" | run_hook)
assert_exit_code "bug without DoD exits 0" 0 $?
assert_output_contains "bug without DoD returns block decision" '"decision"' "$OUTPUT"

echo ""
echo "=== Test: Feature creation without DoD ==="

setup_mock_beans "feature" "A new feature"
OUTPUT=$(make_input 'beans create "Add dark mode" -t feature -d "Add dark mode support"' "Created credfolio2-jkl4 credfolio2-jkl4--add-dark-mode.md" | run_hook)
assert_exit_code "feature without DoD exits 0" 0 $?
assert_output_contains "feature without DoD returns block decision" '"decision"' "$OUTPUT"

echo ""
echo "=== Test: Draft bean creation without DoD ==="

setup_mock_beans "task" "A draft task"
OUTPUT=$(make_input 'beans create "Draft task" -t task -s draft -d "Some draft"' "Created credfolio2-mno5 credfolio2-mno5--draft-task.md" | run_hook)
assert_exit_code "draft without DoD exits 0" 0 $?
assert_output_contains "draft without DoD returns block decision" '"decision"' "$OUTPUT"

echo ""
echo "=== Test: Task creation WITH DoD ==="

setup_mock_beans "task" "$DOD_BODY"
OUTPUT=$(make_input 'beans create "Good task" -t task -d "A good task"' "Created credfolio2-pqr6 credfolio2-pqr6--good-task.md" | run_hook)
assert_exit_code "task with DoD exits 0" 0 $?
assert_output_empty "task with DoD produces no output" "$OUTPUT"

echo ""
echo "=== Test: Bean creation via beans query mutation ==="

# GraphQL mutations should NOT be caught (only CLI `beans create`)
setup_mock_beans "task" "No DoD here"
OUTPUT=$(make_input "beans query 'mutation { createBean(input: { title: \"Test\" }) { id } }'" '{"createBean":{"id":"credfolio2-zzz9"}}' | run_hook)
assert_exit_code "GraphQL mutation exits 0" 0 $?
assert_output_empty "GraphQL mutation produces no output" "$OUTPUT"

echo ""
echo "=== Test: Bean ID extraction from different response formats ==="

setup_mock_beans "task" "No DoD"
OUTPUT=$(make_input 'beans create "Test" -t task' "Created credfolio2-Ab12 credfolio2-Ab12--test.md" | run_hook)
assert_exit_code "mixed-case bean ID exits 0" 0 $?
assert_output_contains "mixed-case bean ID found" 'credfolio2-Ab12' "$OUTPUT"

echo ""
echo "==============================="
echo "Results: $TESTS_PASSED/$TESTS_RUN passed, $TESTS_FAILED failed"
echo "==============================="

if [ "$TESTS_FAILED" -gt 0 ]; then
    exit 1
fi
exit 0
