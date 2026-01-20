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

### 5. Update Bean Checklist as You Go

After completing each checklist item in the bean:

1. Edit the bean file: change `- [ ]` to `- [x]`
2. Include bean file in your commits:

```bash
git add src/... .beans/<bean-file>.md
git commit -m "feat: implement X

- Adds Y functionality
- Includes tests for Z"
```

### 6. Push and Open Pull Request

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

## Test Plan
How to verify this works.

🤖 Generated with [Claude Code](https://claude.com/claude-code)
EOF
)"
```

### 7. Wait for Human Review

**IMPORTANT**: Do NOT merge the PR yourself.

- PRs require human review before merging
- Address review feedback with additional commits
- Wait for explicit approval

### 8. After Merge: Complete the Bean

Once the PR is merged by a human:

```bash
git checkout main
git pull origin main
beans update <bean-id> --status completed
```

## Quick Reference

```bash
# Start work
git checkout main && git pull origin main
git checkout -b feat/<bean-id>-<description>
beans update <bean-id> --status in-progress

# During work (repeat)
# 1. Write failing test
# 2. Make it pass
# 3. Refactor
# 4. Update bean checklist
# 5. Commit with bean file

# Finish work
git push -u origin <branch-name>
gh pr create --title "..." --body "..."
# WAIT for human review and merge

# After merge
git checkout main && git pull origin main
beans update <bean-id> --status completed
```

## Rules

1. **Never commit directly to main** - Always use feature branches
2. **Never merge your own PRs** - Wait for human review
3. **Always pull main before branching** - Avoid merge conflicts
4. **Always use TDD** - Tests before implementation
5. **Always update bean checklists** - Track progress persistently
6. **Include bean files in commits** - Keep state synchronized
