#!/bin/bash
# Tests for start-work.sh helper functions and JSON parsing
# Run: bash .claude/scripts/start-work.test.sh

set -e

PASS=0
FAIL=0

assert_eq() {
    local desc="$1"
    local expected="$2"
    local actual="$3"
    if [ "$expected" = "$actual" ]; then
        echo "  PASS: $desc"
        PASS=$((PASS + 1))
    else
        echo "  FAIL: $desc"
        echo "    expected: '$expected'"
        echo "    actual:   '$actual'"
        FAIL=$((FAIL + 1))
    fi
}

echo "=== start-work.sh tests ==="

# --- JSON parsing tests ---
echo ""
echo "--- JSON parsing (beans query --json output format) ---"

# The actual output format from `beans query --json` has NO .data wrapper
SAMPLE_JSON='{"bean":{"id":"credfolio2-abc1","title":"Add user auth","status":"todo","type":"feature"}}'

TITLE=$(echo "$SAMPLE_JSON" | jq -r '.bean.title // empty')
assert_eq "parses .bean.title correctly" "Add user auth" "$TITLE"

TYPE=$(echo "$SAMPLE_JSON" | jq -r '.bean.type // empty')
assert_eq "parses .bean.type correctly" "feature" "$TYPE"

STATUS=$(echo "$SAMPLE_JSON" | jq -r '.bean.status // empty')
assert_eq "parses .bean.status correctly" "todo" "$STATUS"

# Verify the OLD (buggy) path returns empty
OLD_TITLE=$(echo "$SAMPLE_JSON" | jq -r '.data.bean.title // empty')
assert_eq "old .data.bean.title path returns empty (confirming bug)" "" "$OLD_TITLE"

# --- Null bean test ---
echo ""
echo "--- Null bean handling ---"

NULL_JSON='{"bean":null}'
NULL_TITLE=$(echo "$NULL_JSON" | jq -r '.bean.title // empty')
assert_eq "null bean returns empty title" "" "$NULL_TITLE"

# --- map_type_to_prefix tests ---
echo ""
echo "--- map_type_to_prefix ---"

# Duplicated from start-work.sh (bash scripts can't be easily sourced without side effects)
map_type_to_prefix() {
    case "$1" in
        feature)  echo "feat" ;;
        bug)      echo "fix" ;;
        task)     echo "chore" ;;
        milestone) echo "chore" ;;
        epic)     echo "chore" ;;
        *)        echo "chore" ;;
    esac
}

assert_eq "feature -> feat" "feat" "$(map_type_to_prefix feature)"
assert_eq "bug -> fix" "fix" "$(map_type_to_prefix bug)"
assert_eq "task -> chore" "chore" "$(map_type_to_prefix task)"
assert_eq "milestone -> chore" "chore" "$(map_type_to_prefix milestone)"
assert_eq "epic -> chore" "chore" "$(map_type_to_prefix epic)"
assert_eq "unknown -> chore" "chore" "$(map_type_to_prefix unknown)"

# --- slugify tests ---
echo ""
echo "--- slugify ---"

slugify() {
    echo "$1" \
        | tr '[:upper:]' '[:lower:]' \
        | sed 's/[_ ]/-/g' \
        | sed 's/[^a-z0-9-]//g' \
        | sed 's/-\+/-/g' \
        | sed 's/^-//;s/-$//'
}

assert_eq "basic slugify" "add-user-authentication" "$(slugify "Add user authentication")"
assert_eq "slugify with special chars" "fix-login-bug" "$(slugify "Fix login bug!")"
assert_eq "slugify with underscores" "some-task-name" "$(slugify "some_task_name")"
assert_eq "slugify with multiple spaces" "a-b-c" "$(slugify "a  b  c")"
assert_eq "slugify with uppercase" "hello-world" "$(slugify "HELLO WORLD")"

# --- Summary ---
echo ""
echo "=== Results: $PASS passed, $FAIL failed ==="

if [ $FAIL -gt 0 ]; then
    exit 1
fi
