#!/bin/bash
# Claude Code PreToolUse hook: Validates lint and tests pass before git commit
#
# This hook intercepts "git commit" commands and ensures:
# 1. Linter passes (pnpm lint)
# 2. All tests pass (pnpm test)
#
# Input: JSON via stdin with tool_input.command
# Output: Exit 0 to allow, Exit 2 to block (message via stderr)

set -e

# Read JSON input from stdin
INPUT=$(cat)

# Extract the command being executed
COMMAND=$(echo "$INPUT" | jq -r '.tool_input.command // ""')

# Only intercept git commit commands
if [[ ! "$COMMAND" =~ ^git\ commit ]] && [[ ! "$COMMAND" =~ \&\&\ *git\ commit ]]; then
    exit 0  # Not a commit command, allow it
fi

# Change to project directory
cd "${CLAUDE_PROJECT_DIR:-/workspace}"

echo "Pre-commit hook: Running lint and tests before commit..." >&2

# Run linter
echo "Running linter..." >&2
if ! pnpm lint 2>&1; then
    echo "" >&2
    echo "========================================" >&2
    echo "COMMIT BLOCKED: Linter failed" >&2
    echo "========================================" >&2
    echo "Please fix linting errors before committing." >&2
    echo "Run 'pnpm lint' to see issues." >&2
    echo "Run 'pnpm lint:fix' in frontend or backend to auto-fix." >&2
    exit 2
fi

# Run tests
echo "Running tests..." >&2
if ! pnpm test 2>&1; then
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
