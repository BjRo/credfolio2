---
name: dev-workflow
description: Development workflow for feature implementation. Use when starting work on any bean/task to ensure proper git hygiene, TDD, and PR-based review.
---

# Development Workflow

Follow this workflow for all feature development. It ensures clean git history, test coverage, and human review.

## Starting Work on a Bean

### 1. Ensure Main is Up-to-Date

Before creating a branch, always sync with remote:

```bash
git checkout main
git pull origin main
```

### 2. Create a Feature Branch

Branch naming convention: `<type>/<bean-id>-<short-description>`

```bash
# Examples:
git checkout -b feat/credfolio2-kdtx-docker-compose
git checkout -b fix/credfolio2-abc1-upload-validation
git checkout -b refactor/credfolio2-xyz9-clean-handlers
```

Types:
- `feat/` - New features
- `fix/` - Bug fixes
- `refactor/` - Code improvements without behavior change
- `chore/` - Build, tooling, dependencies
- `docs/` - Documentation only

### 3. Mark Bean as In-Progress

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

### 5. Smoke Test with Browser Automation

After implementation is complete, verify the feature works end-to-end:

1. **Start the dev servers** (if not already running):
   ```bash
   pnpm dev  # Starts frontend on :3000, backend on :8080
   ```

2. **Run browser smoke tests** using `agent-browser`:
   ```bash
   # Open the app
   agent-browser open http://localhost:3000

   # Get interactive elements
   agent-browser snapshot -i

   # Interact and verify (example)
   agent-browser click @e1
   agent-browser wait --load networkidle
   agent-browser snapshot -i  # Check result

   # Take screenshot as evidence (optional)
   agent-browser screenshot ./smoke-test.png

   # Close when done
   agent-browser close
   ```

3. **What to verify**:
   - Page loads without errors
   - Key elements render correctly
   - User interactions work as expected
   - Backend integration functions (API calls succeed)
   - No console errors (`agent-browser errors`)

4. **For backend-only changes**, verify the API:
   ```bash
   agent-browser open http://localhost:8080/health
   agent-browser snapshot
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

### 8. Wait for Human Review

**IMPORTANT**: Do NOT merge the PR yourself.

- PRs require human review before merging
- Address review feedback with additional commits
- Wait for explicit approval

### 9. After Merge: Complete the Bean

Once the PR is merged by a human, use the automated post-merge command:

```bash
/post-merge <bean-id>
```

This automatically:
- Verifies the PR is merged
- Switches to main and pulls latest
- Deletes the local and remote feature branches
- Marks the bean as completed
- Commits and pushes the bean status change

**Manual alternative** (if needed):

```bash
git checkout main && git pull origin main
git branch -d <branch-name>
beans update <bean-id> --status completed
git add .beans/ && git commit -m "chore: Mark <bean-id> as completed" && git push
```

## Quick Reference

```bash
# Start work
git checkout main && git pull origin main
git checkout -b feat/<bean-id>-<description>
beans update <bean-id> --status in-progress

# During work (repeat TDD cycle)
# 1. Write failing test
# 2. Make it pass
# 3. Refactor
# 4. Update bean checklist
# 5. Commit with bean file

# After implementation - smoke test
pnpm dev  # Start servers if not running
agent-browser open http://localhost:3000
agent-browser snapshot -i
# Interact and verify feature works
agent-browser close

# Finish work
git push -u origin <branch-name>
gh pr create --title "..." --body "..."
# WAIT for human review and merge

# After merge (automated)
/post-merge <bean-id>
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

**Every bean MUST include a "Definition of Done" checklist at the end of its body.** Add this when creating the bean:

```markdown
## Definition of Done
- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Visual verification with agent-browser (for UI changes)
- [ ] All other checklist items above are completed
```

**You CANNOT mark a bean as completed while it has unchecked items.** This structurally enforces compliance.

### Before Marking Work Complete

Run this verification sequence:

```bash
# 1. Lint
pnpm lint

# 2. Test
pnpm test

# 3. Visual verification (for UI changes)
pnpm dev  # if not running
# Then use /skill agent-browser to verify

# 4. Check off all Definition of Done items in the bean
# 5. Only then: beans update <bean-id> --status completed
```

**DO NOT skip these steps. DO NOT tell the user "you can verify by running tests" — run them yourself.**
