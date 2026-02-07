---
# credfolio2-e6q3
title: 'Clean up post-merge: slash command delegates to script'
status: in-progress
type: task
priority: normal
created_at: 2026-02-07T17:36:16Z
updated_at: 2026-02-07T19:51:26Z
parent: credfolio2-ynmd
---

The `/post-merge` slash command (`.claude/commands/post-merge.md`) and `scripts/post-merge.sh` duplicate the same post-merge cleanup logic. Consolidate so the slash command delegates to the script, and move the script into `.claude/scripts/`.

## Context
- `.claude/commands/post-merge.md` — Claude Code skill with step-by-step instructions
- `scripts/post-merge.sh` — standalone bash script with identical logic
- `.claude/skills/dev-workflow/SKILL.md` — references `./scripts/post-merge.sh` in 3 places (lines 204, 207, 253)
- `.claude/scripts/start-work.sh` — sister script already in the target directory; follow its conventions

## Checklist
- [ ] Move `scripts/post-merge.sh` to `.claude/scripts/post-merge.sh`
- [ ] Update `.claude/commands/post-merge.md` to instruct Claude to call `.claude/scripts/post-merge.sh <bean-id>` instead of manually executing each step
- [ ] Ensure the slash command still validates the `$ARGUMENTS` bean ID before calling the script
- [ ] Remove duplicated step-by-step instructions from the slash command
- [ ] Update `.claude/skills/dev-workflow/SKILL.md` — change all 3 references from `./scripts/post-merge.sh` to `.claude/scripts/post-merge.sh` (lines 204, 207, 253)
- [ ] Add `--no-gpg-sign` flag to `git commit` in `post-merge.sh` (align with project convention and `start-work.sh`)
- [ ] Add `PROJECT_DIR` variable to `post-merge.sh` (follow `start-work.sh` pattern: `PROJECT_DIR="${CLAUDE_PROJECT_DIR:-/workspace}"`)
- [ ] Verify `scripts/` directory still contains `init-db.sh` — do NOT remove the directory (it is used by `docker-compose.yml`)

## Definition of Done
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review
- [ ] Automated code review passed via `@review-backend` and/or `@review-frontend` subagents (via Task tool)

## Implementation Plan

### Approach

Move the post-merge script into `.claude/scripts/` alongside `start-work.sh`, apply minor fixes to align with project conventions, then rewrite the slash command to be a thin delegation wrapper. Update the dev-workflow skill to reference the new path. The `scripts/` directory stays because `init-db.sh` is still used by Docker.

### Files to Modify

- `scripts/post-merge.sh` -> `.claude/scripts/post-merge.sh` — Move and apply convention fixes
- `.claude/commands/post-merge.md` — Rewrite to delegate to the script
- `.claude/skills/dev-workflow/SKILL.md` — Update 3 path references

### Steps

1. **Move `scripts/post-merge.sh` to `.claude/scripts/post-merge.sh`**
   - `git mv scripts/post-merge.sh .claude/scripts/post-merge.sh`
   - Verify the file retains its executable permission (`chmod +x` if needed)

2. **Apply convention fixes to `.claude/scripts/post-merge.sh`**
   - Add `PROJECT_DIR="${CLAUDE_PROJECT_DIR:-/workspace}"` near the top (after the color variable definitions), matching `start-work.sh` line 16
   - Add `cd "$PROJECT_DIR"` before the `git add .beans/` line, matching `start-work.sh` line 121
   - Change the `git commit` on line 90 from `git commit -m "chore: Mark ${BEAN_ID} as completed"` to `git commit --no-gpg-sign -m "chore: Mark ${BEAN_ID} as completed"`, matching `start-work.sh` line 123 and the CLAUDE.md convention
   - Update the usage comment at the top from `# Usage: ./scripts/post-merge.sh <bean-id>` to `# Usage: .claude/scripts/post-merge.sh <bean-id>` (and same for the Example line)

3. **Rewrite `.claude/commands/post-merge.md` as a thin wrapper**
   - Keep the title ("Post-Merge Cleanup") and the description of what it does
   - Keep the `$ARGUMENTS` section documenting the bean ID parameter
   - Replace the 9 step-by-step instruction sections with a simple instruction block:
     ```markdown
     ## Instructions
     
     1. If `$ARGUMENTS` is empty, ask the user for the bean ID before proceeding.
     2. Run the post-merge cleanup script:
        ```bash
        .claude/scripts/post-merge.sh $ARGUMENTS
        ```
     3. Report the results to the user. If the script exits with an error, share the error message and do not proceed.
     ```
   - This follows the principle of keeping the slash command as a delegation point, not a duplication of logic

4. **Update `.claude/skills/dev-workflow/SKILL.md`**
   - Line 204: Change `./scripts/post-merge.sh <bean-id>` to `.claude/scripts/post-merge.sh <bean-id>`
   - Line 207: Change `./scripts/post-merge.sh credfolio2-abc1` to `.claude/scripts/post-merge.sh credfolio2-abc1`
   - Line 253: Change `./scripts/post-merge.sh <bean-id>` to `.claude/scripts/post-merge.sh <bean-id>`

5. **Verify `scripts/` directory state**
   - Confirm `scripts/init-db.sh` still exists
   - Confirm `scripts/post-merge.sh` is gone
   - Do NOT remove `scripts/` — it is referenced by `docker-compose.yml` line 13

6. **Run lint and tests**
   - `pnpm lint` — should pass (no source code changes)
   - `pnpm test` — should pass (no source code changes)

### Testing Strategy

- Manual: Run `.claude/scripts/post-merge.sh` with no arguments and verify it prints usage info and exits with code 1
- Manual: Run `.claude/scripts/post-merge.sh some-bean-id` from `main` branch and verify it errors with "Already on main branch"
- Manual: Verify `docker-compose.yml` still works (the `init-db.sh` reference is unchanged)
- Verify all 3 updated paths in `SKILL.md` are correct

### Notes

- `scripts/init-db.sh` is a Docker entrypoint script mounted by `docker-compose.yml` and has nothing to do with Claude tooling. It stays in `scripts/`.
- Bean `credfolio2-g1gz` explicitly deferred this migration to this bean. Its historical references are accurate and do not need updating.
- No test file for `post-merge.sh` is required for this bean — the script's logic is straightforward (no slug/prefix functions to unit test like `start-work.sh` has). A test could be added as a follow-up if desired.
