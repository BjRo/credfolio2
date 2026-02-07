# Post-Merge Cleanup

Run this command after a PR has been merged to clean up the branch and complete the bean.

## Arguments

- `$ARGUMENTS` - The bean ID (e.g., `credfolio2-abc1`)

## Instructions

1. If `$ARGUMENTS` is empty, ask the user for the bean ID before proceeding.
2. Run the post-merge cleanup script:
   ```bash
   .claude/scripts/post-merge.sh $ARGUMENTS
   ```
3. Report the results to the user. If the script exits with an error, share the error message and do not proceed.
