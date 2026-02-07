#!/bin/bash
# Claude Code PreToolUse hook: Validates branch names on creation
#
# This hook intercepts branch-creation commands and ensures the branch name
# follows the pattern: <type>/<bean-id>-<description>
#   Types: feat, fix, refactor, chore, docs
#   Bean ID: credfolio2-XXXX or beans-XXXX
#
# Input: JSON via stdin with tool_input.command
# Output: Exit 0 to allow, Exit 2 to block (message via stderr)

set -e

# Read JSON input from stdin
INPUT=$(cat)

# Extract the command being executed
COMMAND=$(echo "$INPUT" | jq -r '.tool_input.command // ""')

# Quick exit: only inspect commands that start with "git" (possibly after &&/;)
# This avoids false positives from git subcommands appearing inside heredocs,
# string arguments, or non-git commands like "gh pr create --body '...git checkout -b...'"
# We extract only the first command token (before any heredoc/string body)
FIRST_LINE=$(echo "$COMMAND" | head -1)
if ! echo "$FIRST_LINE" | grep -qP '(^|&&\s*|;\s*)git\s+(checkout|switch)\s'; then
    exit 0
fi

# Extract a branch name from branch-creation commands.
# Returns the branch name if found, empty string otherwise.
extract_branch_name() {
    local cmd="$1"

    # Only check the first line to avoid matching inside heredocs or string bodies
    local first_line
    first_line=$(echo "$cmd" | head -1)

    # git checkout -b/-B <name> (branch names: alphanumeric, hyphens, slashes, dots, underscores)
    if echo "$first_line" | grep -qP '(^|&&\s*|;\s*)git\s+checkout\s+(-\S+\s+)*-[bB]\s+'; then
        echo "$first_line" | grep -oP '(^|&&\s*|;\s*)git\s+checkout\s+(-\S+\s+)*-[bB]\s+\K[a-zA-Z0-9/_.-]+' | tail -1
        return
    fi

    # git switch -c/-C <name>
    if echo "$first_line" | grep -qP '(^|&&\s*|;\s*)git\s+switch\s+(-\S+\s+)*-[cC]\s+'; then
        echo "$first_line" | grep -oP '(^|&&\s*|;\s*)git\s+switch\s+(-\S+\s+)*-[cC]\s+\K[a-zA-Z0-9/_.-]+' | tail -1
        return
    fi

    echo ""
}

BRANCH_NAME=$(extract_branch_name "$COMMAND")

# Not a branch-creation command, allow it
if [ -z "$BRANCH_NAME" ]; then
    exit 0
fi

# Strip surrounding quotes if present
BRANCH_NAME=$(echo "$BRANCH_NAME" | sed "s/^[\"']//;s/[\"']$//")

# Validate the branch name pattern: <type>/<bean-id>-<description>
VALID_PATTERN='^(feat|fix|refactor|chore|docs)/(credfolio2-[a-zA-Z0-9]+|beans-[a-zA-Z0-9]+)-.+$'

if echo "$BRANCH_NAME" | grep -qP "$VALID_PATTERN"; then
    exit 0
fi

# Block with helpful error message
echo "" >&2
echo "========================================" >&2
echo "BLOCKED: Invalid branch name" >&2
echo "========================================" >&2
echo "" >&2
echo "Branch name '${BRANCH_NAME}' does not match the required pattern." >&2
echo "" >&2
echo "Expected: <type>/<bean-id>-<description>" >&2
echo "  Types: feat, fix, refactor, chore, docs" >&2
echo "  Example: feat/credfolio2-abc1-add-user-auth" >&2
echo "" >&2
echo "Use the start-work script instead:" >&2
echo "  .claude/scripts/start-work.sh <bean-id>" >&2
exit 2
