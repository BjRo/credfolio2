---
# credfolio2-6ebh
title: Create workflow automation scripts
status: in-progress
type: task
created_at: 2026-02-03T17:04:45Z
updated_at: 2026-02-03T17:04:45Z
---

Create shell scripts to automate repetitive dev workflow tasks:

## Scope
1. `scripts/start-work.sh` - Automates starting work on a bean:
   - Ensure main is up-to-date
   - Create feature branch with proper naming
   - Mark bean as in-progress
   - Commit the bean status change

2. `scripts/post-merge.sh` - Automates post-merge cleanup:
   - Verify PR is merged
   - Switch to main and pull
   - Delete local and remote branches
   - Mark bean as completed
   - Commit and push

3. Update dev-workflow skill to reference these scripts
4. Update CLAUDE.md with script usage

## Checklist
- [x] Create start-work.sh script
- [x] Create post-merge.sh script
- [x] Update dev-workflow SKILL.md
- [x] Test scripts work correctly

## Definition of Done
- [x] Scripts created and executable
- [x] Documentation updated (dev-workflow skill)
- [x] Scripts tested manually