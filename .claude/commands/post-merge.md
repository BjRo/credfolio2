# Post-Merge Cleanup

Run this command after a PR has been merged to clean up the branch and complete the bean.

## Arguments

- `$ARGUMENTS` - The bean ID (e.g., `credfolio2-abc1`)

## Instructions

Perform the following steps in order:

### 1. Validate Arguments

If no bean ID is provided in `$ARGUMENTS`, ask the user for the bean ID before proceeding.

### 2. Get Current Branch Info

```bash
git branch --show-current
```

Save the current branch name - you'll need it for cleanup.

### 3. Check PR Status

```bash
gh pr view --json state,headRefName,mergedAt
```

- If state is not "MERGED", inform the user and stop
- If no PR exists for this branch, inform the user and stop

### 4. Switch to Main and Pull Latest

```bash
git checkout main
git pull origin main
```

### 5. Delete the Local Feature Branch

```bash
git branch -d <branch-name>
```

Use the branch name from step 2. If deletion fails (unmerged changes), use `-D` flag after confirming with user.

### 6. Delete the Remote Branch (if exists)

```bash
git push origin --delete <branch-name> 2>/dev/null || echo "Remote branch already deleted or doesn't exist"
```

### 7. Mark Bean as Completed

```bash
beans update $ARGUMENTS --status completed
```

### 8. Commit and Push the Bean Status Change

```bash
git add .beans/
git commit -m "chore: Mark $ARGUMENTS as completed"
git push origin main
```

### 9. Report Success

Summarize what was done:
- Branch deleted (local and remote)
- Bean marked as completed
- Changes pushed to main
