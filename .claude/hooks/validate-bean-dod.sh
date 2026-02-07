#!/bin/bash
# Claude Code PostToolUse hook: Validates Definition of Done on bean creation
#
# This hook intercepts "beans create" commands and checks that the newly
# created bean includes the required Definition of Done checklist.
#
# Input: JSON via stdin with tool_input.command and tool_response
# Output: Exit 0 with JSON {"decision":"block","reason":"..."} if DoD missing
#         Exit 0 silently if DoD present or command is not beans create

# Read JSON input from stdin
INPUT=$(cat)

# Extract the command being executed
COMMAND=$(echo "$INPUT" | jq -r '.tool_input.command // ""')

# Fast path: only care about "beans create" commands (not "beans query", "beans update", etc.)
if [[ ! "$COMMAND" =~ ^beans[[:space:]]+create[[:space:]] ]]; then
    exit 0
fi

# Extract bean ID from tool_response (format: "Created credfolio2-xxxx ...")
TOOL_RESPONSE=$(echo "$INPUT" | jq -r '.tool_response // ""')
BEAN_ID=$(echo "$TOOL_RESPONSE" | grep -oP 'credfolio2-[a-zA-Z0-9]+' | head -1)

if [ -z "$BEAN_ID" ]; then
    # Could not extract bean ID — don't block, just exit silently
    exit 0
fi

# Query the bean for its type and body
BEAN_JSON=$(beans query "{ bean(id: \"$BEAN_ID\") { type body } }" --json 2>/dev/null)

BEAN_TYPE=$(echo "$BEAN_JSON" | jq -r '.bean.type // ""')
BEAN_BODY=$(echo "$BEAN_JSON" | jq -r '.bean.body // ""')

# Skip validation for epics and milestones (they don't need DoD)
if [ "$BEAN_TYPE" = "epic" ] || [ "$BEAN_TYPE" = "milestone" ]; then
    exit 0
fi

# Find the template file
TEMPLATE_DIR="${CLAUDE_PROJECT_DIR:-.}/.claude/templates"
TEMPLATE_FILE="$TEMPLATE_DIR/definition-of-done.md"

if [ ! -f "$TEMPLATE_FILE" ]; then
    # Template file not found — can't validate, exit silently
    exit 0
fi

# Extract key phrases from the template (strip "- [ ] " prefix from checklist items)
# Check that each required phrase appears in the bean body (case-insensitive)
MISSING_ITEMS=()

while IFS= read -r line; do
    # Extract the text after "- [ ] "
    phrase=$(echo "$line" | sed 's/^- \[ \] //')
    if [ -n "$phrase" ] && [ "$phrase" != "$line" ]; then
        # Check if this phrase (or a key substring) appears in the bean body
        if ! echo "$BEAN_BODY" | grep -qi "$phrase"; then
            MISSING_ITEMS+=("$phrase")
        fi
    fi
done < "$TEMPLATE_FILE"

# If any items are missing, return a block decision
if [ ${#MISSING_ITEMS[@]} -gt 0 ]; then
    REASON="Bean $BEAN_ID is missing the Definition of Done checklist. Read .claude/templates/definition-of-done.md and append it to the bean body using beans update."
    echo "{\"decision\":\"block\",\"reason\":\"$REASON\"}"
    exit 0
fi

# All DoD items present — pass silently
exit 0
