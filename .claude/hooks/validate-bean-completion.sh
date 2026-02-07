#!/bin/bash
# Claude Code TaskCompleted hook: Validates bean checklist completion
#
# This hook fires when a task completes and checks whether the associated
# bean has any unchecked checklist items. If so, it blocks completion with
# exit code 2 and prints the remaining items to stderr.
#
# Input: JSON via stdin with task_id, task_subject, task_description, etc.
# Output: Exit 0 to allow completion
#         Exit 2 with stderr message to block completion
#
# Bean ID extraction strategy:
#   1. Search task_subject and task_description for credfolio2-XXXX pattern
#   2. Fall back to extracting from current git branch name

# Do NOT use set -e — we want graceful fallback to exit 0 on any failure

# Read JSON input from stdin
INPUT=$(cat)

# Extract fields from the JSON
TASK_SUBJECT=$(echo "$INPUT" | jq -r '.task_subject // ""')
TASK_DESCRIPTION=$(echo "$INPUT" | jq -r '.task_description // ""')

# --- Strategy 1: Extract bean ID from task_subject or task_description ---
BEAN_ID=""

# Check task_subject first
if [ -n "$TASK_SUBJECT" ]; then
    BEAN_ID=$(echo "$TASK_SUBJECT" | grep -oP 'credfolio2-[a-zA-Z0-9]+' | head -1 || true)
fi

# If not found, check task_description
if [ -z "$BEAN_ID" ] && [ -n "$TASK_DESCRIPTION" ]; then
    BEAN_ID=$(echo "$TASK_DESCRIPTION" | grep -oP 'credfolio2-[a-zA-Z0-9]+' | head -1 || true)
fi

# --- Strategy 2: Extract from git branch name ---
if [ -z "$BEAN_ID" ]; then
    BRANCH=$(git branch --show-current 2>/dev/null || true)
    if [ -n "$BRANCH" ] && [[ "$BRANCH" =~ ^[a-z]+/(credfolio2-[a-zA-Z0-9]+)-.* ]]; then
        BEAN_ID="${BASH_REMATCH[1]}"
    fi
fi

# --- No bean found — allow completion silently ---
if [ -z "$BEAN_ID" ]; then
    exit 0
fi

# --- Query the bean body ---
BEAN_JSON=$(timeout 5s beans query "{ bean(id: \"$BEAN_ID\") { id title status body } }" --json 2>/dev/null || true)

if [ -z "$BEAN_JSON" ]; then
    # Query failed or timed out — don't block on infrastructure failures
    exit 0
fi

# Extract the bean body
BEAN_BODY=$(echo "$BEAN_JSON" | jq -r '.bean.body // empty' 2>/dev/null || true)

if [ -z "$BEAN_BODY" ]; then
    # No body or extraction failed — allow completion
    exit 0
fi

# --- Scan for unchecked checklist items ---
UNCHECKED=$(echo "$BEAN_BODY" | grep -P '^\- \[ \] ' || true)

if [ -z "$UNCHECKED" ]; then
    # No unchecked items — allow completion
    exit 0
fi

# --- Block completion: unchecked items remain ---
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
