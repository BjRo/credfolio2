#!/bin/bash
set -e

# Post-merge cleanup: verify merge, delete branches, complete bean
# Usage: ./scripts/post-merge.sh <bean-id>
# Example: ./scripts/post-merge.sh credfolio2-abc1

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

usage() {
    echo "Usage: $0 <bean-id>"
    echo ""
    echo "Arguments:"
    echo "  bean-id  The bean ID (e.g., credfolio2-abc1)"
    echo ""
    echo "This script should be run after your PR has been merged."
    echo "It will:"
    echo "  - Verify the PR is merged"
    echo "  - Switch to main and pull latest"
    echo "  - Delete local and remote feature branches"
    echo "  - Mark the bean as completed"
    echo "  - Commit and push the bean status change"
    exit 1
}

# Validate arguments
if [ $# -lt 1 ]; then
    usage
fi

BEAN_ID="$1"

echo -e "${YELLOW}Running post-merge cleanup for ${BEAN_ID}...${NC}"

# 1. Get current branch
CURRENT_BRANCH=$(git branch --show-current)
echo -e "\n${GREEN}[1/7]${NC} Current branch: ${CURRENT_BRANCH}"

if [ "$CURRENT_BRANCH" = "main" ]; then
    echo -e "${RED}Error: Already on main branch.${NC}"
    echo "Please run this from your feature branch, or specify the branch manually."
    exit 1
fi

# 2. Check PR status
echo -e "\n${GREEN}[2/7]${NC} Checking PR status..."
PR_INFO=$(gh pr view --json state,headRefName,mergedAt 2>/dev/null || echo '{"error": true}')

if echo "$PR_INFO" | grep -q '"error"'; then
    echo -e "${RED}Error: No PR found for branch '${CURRENT_BRANCH}'${NC}"
    echo "Please create and merge a PR before running this script."
    exit 1
fi

PR_STATE=$(echo "$PR_INFO" | jq -r '.state')
if [ "$PR_STATE" != "MERGED" ]; then
    echo -e "${RED}Error: PR is not merged (state: ${PR_STATE})${NC}"
    echo "Please wait for the PR to be merged before running this script."
    exit 1
fi

echo -e "PR state: ${GREEN}MERGED${NC}"

# 3. Switch to main and pull
echo -e "\n${GREEN}[3/7]${NC} Switching to main and pulling latest..."
git checkout main
git pull origin main

# 4. Delete local branch
echo -e "\n${GREEN}[4/7]${NC} Deleting local branch '${CURRENT_BRANCH}'..."
if ! git branch -d "$CURRENT_BRANCH" 2>/dev/null; then
    echo -e "${YELLOW}Warning: Branch has unmerged changes. Force deleting...${NC}"
    git branch -D "$CURRENT_BRANCH"
fi

# 5. Delete remote branch
echo -e "\n${GREEN}[5/7]${NC} Deleting remote branch..."
git push origin --delete "$CURRENT_BRANCH" 2>/dev/null || echo "Remote branch already deleted or doesn't exist"

# 6. Mark bean as completed
echo -e "\n${GREEN}[6/7]${NC} Marking bean as completed..."
beans update "$BEAN_ID" --status completed

# 7. Commit and push
echo -e "\n${GREEN}[7/7]${NC} Committing and pushing bean status change..."
git add .beans/
git commit -m "chore: Mark ${BEAN_ID} as completed"
git push origin main

echo -e "\n${GREEN}✓ Post-merge cleanup complete!${NC}"
echo ""
echo "Summary:"
echo "  - Branch '${CURRENT_BRANCH}' deleted (local and remote)"
echo "  - Bean '${BEAN_ID}' marked as completed"
echo "  - Changes pushed to main"
