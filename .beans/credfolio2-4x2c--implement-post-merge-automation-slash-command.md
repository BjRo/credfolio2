---
# credfolio2-4x2c
title: Implement post-merge automation slash command
status: in-progress
type: task
created_at: 2026-01-20T11:51:16Z
updated_at: 2026-01-20T11:51:16Z
---

Create a /post-merge slash command that automates cleanup after a PR is merged:
- Verify PR is merged
- Switch to main and pull latest
- Delete local and remote feature branches
- Mark bean as completed
- Commit and push the bean status change

## Checklist
- [x] Create .claude/commands/post-merge.md
- [x] Create .claude/scripts/get-current-bean.sh helper
- [x] Update .claude/skills/dev-workflow/SKILL.md documentation
- [x] Test gh authentication works
- [ ] Verify the workflow