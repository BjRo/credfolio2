---
name: dev-workflow
description: Development workflow for feature implementation. Use when starting work on any bean/task to ensure proper git hygiene, TDD, and PR-based review.
---

# Development Workflow

Follow this workflow for all feature development. It ensures clean git history, test coverage, and human review.

## Starting Work on a Bean

### Quick Start (Recommended)

Use the automated script to set up your branch and mark the bean in-progress:

```bash
.claude/scripts/start-work.sh <bean-id>

# Examples:
.claude/scripts/start-work.sh credfolio2-kdtx
.claude/scripts/start-work.sh credfolio2-abc1
.claude/scripts/start-work.sh credfolio2-xyz9
```

This script automatically:
- Ensures main is up-to-date
- Queries the bean for its type and title
- Derives the branch prefix (feature->feat, bug->fix, task->chore) and slugifies the title
- Creates a properly named feature branch (e.g., `feat/credfolio2-abc1-add-user-auth`)
- Marks the bean as in-progress
- Commits the bean status change

**Note**: A PreToolUse hook validates all branch-creation commands against the naming convention. Manual `git checkout -b` with non-conforming names will be blocked.

### Manual Steps (Reference)

If you need to do this manually:

#### 1. Ensure Main is Up-to-Date

```bash
git checkout main
git pull origin main
```

#### 2. Create a Feature Branch

Branch naming convention: `<type>/<bean-id>-<short-description>`

Types:
- `feat/` - New features
- `fix/` - Bug fixes
- `refactor/` - Code improvements without behavior change
- `chore/` - Build, tooling, dependencies
- `docs/` - Documentation only

#### 3. Mark Bean as In-Progress

```bash
beans update <bean-id> --status in-progress
git add .beans/
git commit -m "chore: start work on <bean-id>"
```

### 4. Develop Using TDD

Follow the TDD skill (`/skill tdd`):

1. **RED**: Write a failing test
2. **GREEN**: Write minimum code to pass
3. **REFACTOR**: Clean up while green

Commit frequently with meaningful messages.

### 5. Visual Verification via QA Subagent

After implementation is complete, verify the feature works end-to-end by launching the `@qa` subagent. This keeps verbose browser automation output out of the main conversation context.

**Launch the QA subagent using the Task tool** (use exactly these parameters):

```
Task tool call:
  subagent_type: "qa"
  description: "Visual verification of <feature>"
  prompt: "Verify <feature description>. Start dev servers if needed, navigate to <URL>, test <interactions>, check for errors, and report pass/fail."
```

**IMPORTANT**: Always launch as a subagent via the Task tool. The `@qa` subagent handles dev server management, browser automation, and error checking automatically.

**After the QA subagent returns:**
- Read the summary to confirm the feature works
- If it reports failures, fix the issues and re-run
- If it reports console errors, investigate and address them

**For backend-only changes**, verify the API by including API endpoint checks in the QA subagent prompt:
```
prompt: "Verify the health endpoint. Start dev servers if needed, navigate to http://localhost:8080/health, confirm it returns a valid response."
```

This self-verification catches integration issues before human review.

### 6. Update Bean Checklist as You Go

After completing each checklist item in the bean:

1. Edit the bean file: change `- [ ]` to `- [x]`
2. Include bean file in your commits:

```bash
git add src/... .beans/<bean-file>.md
git commit -m "feat: implement X

- Adds Y functionality
- Includes tests for Z"
```

### 7. Push and Open Pull Request

When the bean checklist is complete:

```bash
git push -u origin <branch-name>
```

Create a PR using GitHub CLI:

```bash
gh pr create --title "<type>: <description>" --body "$(cat <<'EOF'
## Summary
Brief description of what this PR does.

## Bean
Closes beans-<id>

## Checklist
- [ ] Tests pass (`pnpm test`)
- [ ] Build succeeds (`pnpm build`)
- [ ] TDD followed (tests written first)
- [ ] Smoke tested with browser automation

## Test Plan
How to verify this works.

🤖 Generated with [Claude Code](https://claude.com/claude-code)
EOF
)"
```

### 8. Run Automated Code Reviews

After creating the PR, launch the two code review agents as **separate subagents** so they don't pollute the main conversation context. You MUST use the Task tool to run them.

**Determine which reviewers to run** based on changed files:

```bash
# Check which areas have changes
gh pr diff --name-only
```

- If `src/backend/` files changed → launch **@review-backend**
- If `src/frontend/` files changed → launch **@review-frontend**
- If both changed → launch **both in parallel**

**Launch the reviewers using the Task tool** (use exactly these parameters):

For **@review-backend**:
```
Task tool call:
  subagent_type: "review-backend"
  description: "Backend code review"
  prompt: "Review the current PR. Post your findings as PR comments using the gh CLI."
```

For **@review-frontend**:
```
Task tool call:
  subagent_type: "review-frontend"
  description: "Frontend code review"
  prompt: "Review the current PR. Post your findings as PR comments using the gh CLI."
```

**IMPORTANT**: Always launch these as subagents via the Task tool. Never invoke review agents directly in the main conversation — that defeats the purpose of keeping the context clean. Both `@review-backend` and `@review-frontend` are named agents in `.claude/agents/`.

**After the reviews complete:**
- Read the review summaries returned by the subagents
- If either finds `🔴 CRITICAL` issues, address them before requesting human review
- `🟡 WARNING` items should generally be addressed but use your judgment
- `🔵 SUGGESTION` items are optional improvements
- Report the review outcomes to the user

### 9. Wait for Human Review

**IMPORTANT**: Do NOT merge the PR yourself.

- PRs require human review before merging
- Address review feedback (both automated and human) with additional commits
- Wait for explicit approval

### 10. After Merge: Complete the Bean

Once the PR is merged by a human, use the automated post-merge script:

```bash
.claude/scripts/post-merge.sh <bean-id>

# Example:
.claude/scripts/post-merge.sh credfolio2-abc1
```

This script automatically:
- Verifies the PR is merged
- Switches to main and pulls latest
- Deletes the local and remote feature branches
- Marks the bean as completed
- Commits and pushes the bean status change

**Note**: Run this from your feature branch (not main).

**Manual alternative** (if needed):

```bash
git checkout main && git pull origin main
git branch -d <branch-name>
beans update <bean-id> --status completed
git add .beans/ && git commit -m "chore: Mark <bean-id> as completed" && git push
```

## Quick Reference

```bash
# Start work (automated)
.claude/scripts/start-work.sh <bean-id>

# During work (repeat TDD cycle)
# 1. Write failing test
# 2. Make it pass
# 3. Refactor
# 4. Update bean checklist
# 5. Commit with bean file

# After implementation - visual verification
# Launch @qa subagent via Task tool to verify feature works

# Finish work
git push -u origin <branch-name>
gh pr create --title "..." --body "..."

# Run automated code reviews (as parallel subagents)
# Launch @review-backend and @review-frontend via Task tool
# Address any CRITICAL findings, then WAIT for human review and merge

# After merge (from feature branch)
.claude/scripts/post-merge.sh <bean-id>
```

## Rules

1. **Never commit directly to main** - Always use feature branches
2. **Never merge your own PRs** - Wait for human review
3. **Always pull main before branching** - Avoid merge conflicts
4. **Always use TDD** - Tests before implementation
5. **Always smoke test** - Verify features work in browser before PR
6. **Always update bean checklists** - Track progress persistently
7. **Include bean files in commits** - Keep state synchronized

## Mandatory Definition of Done

**Every bean MUST include a "Definition of Done" checklist at the end of its body.** The canonical template is at `.claude/templates/definition-of-done.md`. A PostToolUse hook automatically validates this on `beans create` commands and will prompt you to add it if missing.

**You CANNOT mark a bean as completed while it has unchecked items.** This structurally enforces compliance.

### Before Marking Work Complete

Run this verification sequence:

```bash
# 1. Lint
pnpm lint

# 2. Test
pnpm test

# 3. Visual verification (for UI changes)
# Launch @qa subagent via Task tool to verify

# 4. Check off all Definition of Done items in the bean
# 5. Only then: beans update <bean-id> --status completed
```

**DO NOT skip these steps. DO NOT tell the user "you can verify by running tests" — run them yourself.**
