#!/bin/bash
# Claude Code PreToolUse hook: Validates lint and tests pass before git commit
#
# This hook intercepts "git commit" commands and ensures:
# 1. Linter passes (pnpm lint)
# 2. All tests pass (pnpm test)
#
# Input: JSON via stdin with tool_input.command
# Output: Exit 0 to allow, Exit 2 to block (message via stderr)
#
# Note: No `set -e` — we use explicit error handling to avoid exit code 1
# which Claude Code treats as a hook error (allowing the tool call to proceed).

# Read JSON input from stdin
INPUT=$(cat)

# Extract the command being executed
# If jq fails (missing, malformed input), default to empty string
COMMAND=$(echo "$INPUT" | jq -r '.tool_input.command // ""' 2>/dev/null || echo "")

# Only inspect the first line to avoid matching inside heredocs or string bodies
FIRST_LINE=$(echo "$COMMAND" | head -1)

# Match git commit as a command, not inside quoted strings or as a flag value.
# The pattern matches:
#   - git commit (at start of line)
#   - git <flags> commit (with flags like -C /workspace before commit)
#   - <cmd> && git commit (chained commands)
#   - <cmd> ; git commit (sequential commands)
# It does NOT match:
#   - echo "git commit" (inside string literals)
#   - git merge --commit (commit as a flag value, not a subcommand)
if ! echo "$FIRST_LINE" | grep -qP '(^|&&\s*|;\s*)git\s+(-\S+\s+\S+\s+)*commit(\s|$)'; then
    exit 0  # Not a commit command, allow it
fi

# Change to project directory
cd "${CLAUDE_PROJECT_DIR:-/workspace}"

echo "Pre-commit hook: Running lint and tests before commit..." >&2

# Run linter with timeout (5 minutes)
echo "Running linter..." >&2
timeout 300 pnpm lint >&2
lint_exit=$?
if [ "$lint_exit" -eq 124 ]; then
    echo "" >&2
    echo "========================================" >&2
    echo "COMMIT BLOCKED: Linter timed out after 5 minutes" >&2
    echo "========================================" >&2
    echo "The linter did not complete within 5 minutes." >&2
    echo "Check for Turborepo cache corruption: rm -rf .turbo" >&2
    exit 2
elif [ "$lint_exit" -ne 0 ]; then
    echo "" >&2
    echo "========================================" >&2
    echo "COMMIT BLOCKED: Linter failed" >&2
    echo "========================================" >&2
    echo "Please fix linting errors before committing." >&2
    echo "Run 'pnpm lint' to see issues." >&2
    echo "Run 'pnpm lint:fix' in frontend or backend to auto-fix." >&2
    exit 2
fi

# Run tests with timeout (5 minutes)
echo "Running tests..." >&2
timeout 300 pnpm test >&2
test_exit=$?
if [ "$test_exit" -eq 124 ]; then
    echo "" >&2
    echo "========================================" >&2
    echo "COMMIT BLOCKED: Tests timed out after 5 minutes" >&2
    echo "========================================" >&2
    echo "Tests did not complete within 5 minutes." >&2
    echo "Check for Turborepo cache corruption: rm -rf .turbo" >&2
    exit 2
elif [ "$test_exit" -ne 0 ]; then
    echo "" >&2
    echo "========================================" >&2
    echo "COMMIT BLOCKED: Tests failed" >&2
    echo "========================================" >&2
    echo "Please fix failing tests before committing." >&2
    echo "Run 'pnpm test' to see failures." >&2
    exit 2
fi

echo "Pre-commit hook: All checks passed!" >&2
exit 0
