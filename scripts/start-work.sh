#!/bin/bash
set -e

# Start work on a bean: create branch, mark in-progress, commit
# Usage: ./scripts/start-work.sh <bean-id> <type> <short-description>
# Example: ./scripts/start-work.sh credfolio2-abc1 feat docker-compose

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

usage() {
    echo "Usage: $0 <bean-id> <type> <short-description>"
    echo ""
    echo "Arguments:"
    echo "  bean-id           The bean ID (e.g., credfolio2-abc1)"
    echo "  type              Branch type: feat, fix, refactor, chore, docs"
    echo "  short-description Brief description for branch name (use hyphens)"
    echo ""
    echo "Example:"
    echo "  $0 credfolio2-abc1 feat docker-compose"
    echo "  $0 credfolio2-xyz9 fix upload-validation"
    exit 1
}

# Validate arguments
if [ $# -lt 3 ]; then
    usage
fi

BEAN_ID="$1"
TYPE="$2"
DESCRIPTION="$3"

# Validate type
case "$TYPE" in
    feat|fix|refactor|chore|docs)
        ;;
    *)
        echo -e "${RED}Error: Invalid type '$TYPE'${NC}"
        echo "Valid types: feat, fix, refactor, chore, docs"
        exit 1
        ;;
esac

BRANCH_NAME="${TYPE}/${BEAN_ID}-${DESCRIPTION}"

echo -e "${YELLOW}Starting work on ${BEAN_ID}...${NC}"

# 1. Ensure we're on main and up-to-date
echo -e "\n${GREEN}[1/5]${NC} Ensuring main is up-to-date..."
git checkout main
git pull origin main

# 2. Check if bean exists
echo -e "\n${GREEN}[2/5]${NC} Verifying bean exists..."
if ! beans query "{ bean(id: \"${BEAN_ID}\") { id title status } }" --json | grep -q "\"id\""; then
    echo -e "${RED}Error: Bean '${BEAN_ID}' not found${NC}"
    exit 1
fi

# Show bean info
beans query "{ bean(id: \"${BEAN_ID}\") { id title status type } }"

# 3. Create feature branch
echo -e "\n${GREEN}[3/5]${NC} Creating branch '${BRANCH_NAME}'..."
git checkout -b "$BRANCH_NAME"

# 4. Mark bean as in-progress
echo -e "\n${GREEN}[4/5]${NC} Marking bean as in-progress..."
beans update "$BEAN_ID" --status in-progress

# 5. Commit the bean status change
echo -e "\n${GREEN}[5/5]${NC} Committing bean status change..."
git add .beans/
git commit -m "chore: Start work on ${BEAN_ID}"

echo -e "\n${GREEN}✓ Ready to work!${NC}"
echo -e "Branch: ${YELLOW}${BRANCH_NAME}${NC}"
echo -e "Bean:   ${YELLOW}${BEAN_ID}${NC} (in-progress)"
echo ""
echo "Next steps:"
echo "  1. Implement the feature using TDD"
echo "  2. Update bean checklist as you go"
echo "  3. Run: pnpm lint && pnpm test"
echo "  4. Push and create PR: git push -u origin ${BRANCH_NAME}"
