---
# credfolio2-o624
title: Add PostToolUse hook to validate Definition of Done on bean creation
status: todo
type: task
created_at: 2026-02-07T16:25:29Z
updated_at: 2026-02-07T16:25:29Z
parent: credfolio2-ynmd
---

Add a PostToolUse hook that checks newly created beans for the required Definition of Done checklist and prompts Claude to add it if missing.

## Why

The Definition of Done checklist is mandatory for all actionable beans (tasks, bugs, features). Currently this is only enforced by instructions in CLAUDE.md and the dev-workflow skill. If Claude creates a bean without the DoD, nothing catches it until much later (or never). A PostToolUse hook gives immediate feedback.

## What

### 1. Definition of Done template file

Create `.claude/templates/definition-of-done.md` as the single source of truth:
```markdown
## Definition of Done
- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Visual verification via QA subagent (for UI changes)
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review
- [ ] Automated code review via review-backend and/or review-frontend subagents
```

This file is read by the hook and can also be referenced by the dev-workflow skill, keeping the template in one place.

### 2. Hook script

Create `.claude/hooks/validate-bean-dod.sh` that:
- Reads the PostToolUse JSON input from stdin
- Checks if the Bash command was a `beans create` command
- Extracts the bean ID from the tool response output
- Queries the bean for its type and body via `beans query`
- **Skips validation** for epics and milestones (they don't have DoD)
- Checks the bean body for the required DoD items (reads from the template file)
- If DoD is missing or incomplete: returns `decision: "block"` with a reason listing what's missing
- If DoD is present: exits 0 silently

### 3. Hook registration

In `.claude/settings.json`:
```json
{
  "hooks": {
    "PostToolUse": [
      {
        "matcher": "Bash",
        "hooks": [
          {
            "type": "command",
            "command": "\"$CLAUDE_PROJECT_DIR\"/.claude/hooks/validate-bean-dod.sh"
          }
        ]
      }
    ]
  }
}
```

### 4. Integration

- Update the dev-workflow skill to reference the template file
- Update CLAUDE.md to reference the template file as the canonical DoD
- This ensures the hook, skill, and documentation all use the same checklist

## Example flow

```
# Claude creates a bean without DoD:
beans create "Add login page" -t feature -s todo -d "Create a login page with email/password"

# Hook detects missing DoD, returns:
{
  "decision": "block",
  "reason": "Bean credfolio2-abc1 is missing the Definition of Done checklist. Read .claude/templates/definition-of-done.md and append it to the bean body using beans update."
}

# Claude reads the template and updates the bean
```

## Notes

- The hook should be fast — only triggers on `beans create`, exits 0 immediately for all other Bash commands
- The template file approach means credfolio2-xflj (updating references to new subagents) only needs to update one file
- Consider also validating on `beans update` if the DoD gets accidentally removed (optional, might be noisy)

## Definition of Done
- [ ] Template file created at `.claude/templates/definition-of-done.md`
- [ ] Hook script created at `.claude/hooks/validate-bean-dod.sh`
- [ ] Hook registered in `.claude/settings.json`
- [ ] Tested: bean creation without DoD → Claude prompted to add it
- [ ] Tested: bean creation with DoD → passes silently
- [ ] Tested: epic/milestone creation → skipped (no validation)
- [ ] Tested: non-beans-create Bash commands → ignored
- [ ] dev-workflow skill and CLAUDE.md updated to reference template file
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Branch pushed and PR created for human review