---
# credfolio2-o624
title: Add PostToolUse hook to validate Definition of Done on bean creation
status: in-progress
type: task
priority: normal
created_at: 2026-02-07T16:25:29Z
updated_at: 2026-02-07T18:42:12Z
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
- [ ] Visual verification via `@qa` subagent (via Task tool, for UI changes)
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review
- [ ] Automated code review passed via `@review-backend` and/or `@review-frontend` subagents (via Task tool)
```

This file is read by the hook and can also be referenced by the dev-workflow skill, keeping the template in one place.

### 2. Hook script

Create `.claude/hooks/validate-bean-dod.sh` that:

**Early exit (fast path):**
1. Read PostToolUse JSON from stdin (see protocol details below)
2. Extract command from `.tool_input.command`
3. If command does not match `beans create` → exit 0 immediately (no output)

**Validation path:**
4. Extract bean ID from `.tool_response` (the `beans create` output contains `Created credfolio2-xxxx ...`)
5. Query the bean for its type and body via `beans query '{ bean(id: "...") { type body } }' --json`
6. **Skip validation** for epics and milestones (they don't need DoD)
7. **Validate all other types** including draft beans (tasks, bugs, features in any status)
8. Check the bean body for required DoD items using **key phrase matching** (see matching strategy below)
9. If DoD is missing or incomplete: exit 0 with JSON `{"decision": "block", "reason": "..."}`
10. If DoD is present: exit 0 silently (no output)

**Key phrase matching strategy:**
- Read the template file at `.claude/templates/definition-of-done.md`
- Extract all `- [ ]` lines, strip the checkbox prefix to get the key phrases
- For each key phrase, check if a substring appears in the bean body (case-insensitive)
- Required substrings to check: `Tests written`, `pnpm lint`, `pnpm test`, `All other checklist items`, `Branch pushed and PR created`, `code review`
- The `Visual verification` item is NOT required (only applies to UI changes) — but still checked as part of the template
- This approach is resilient to minor rewording while still catching missing items

### 3. PostToolUse hook protocol

**Stdin JSON schema** (this is the first PostToolUse hook in the codebase):
```json
{
  "session_id": "abc123",
  "transcript_path": "/path/to/transcript.jsonl",
  "cwd": "/workspace",
  "permission_mode": "default",
  "hook_event_name": "PostToolUse",
  "tool_name": "Bash",
  "tool_input": { "command": "beans create ..." },
  "tool_response": "Created credfolio2-xxxx credfolio2-xxxx--slug.md",
  "tool_use_id": "toolu_01ABC123..."
}
```

Key differences from PreToolUse hooks:
- Has `tool_response` field (the tool's output — the tool has already executed)
- `decision: "block"` does NOT prevent execution (bean is already created). It injects feedback into Claude's context telling it to fix the issue.

**Exit codes:**
| Code | Behavior |
|------|----------|
| 0 | Success. Stdout is parsed as JSON if present. |
| 2 | Shows stderr to Claude as feedback. JSON on stdout is ignored. |
| Other | Non-blocking error, shown in verbose mode only. |

**Output format** (exit 0 + JSON on stdout):
```json
{
  "decision": "block",
  "reason": "Bean credfolio2-abc1 is missing the Definition of Done checklist. Read .claude/templates/definition-of-done.md and append it to the bean body using beans update."
}
```

### 4. Hook registration

Add a `PostToolUse` entry to `.claude/settings.json` (alongside existing `PreToolUse`, `SessionStart`, `PreCompact`):
```json
{
  "hooks": {
    "PostToolUse": [
      {
        "matcher": "Bash",
        "hooks": [
          {
            "type": "command",
            "command": "$CLAUDE_PROJECT_DIR/.claude/hooks/validate-bean-dod.sh"
          }
        ]
      }
    ]
  }
}
```

Note: existing hooks use `$CLAUDE_PROJECT_DIR` without quoting (not `"$CLAUDE_PROJECT_DIR"`). Follow the same convention for consistency.

### 5. Integration

- Update the dev-workflow skill to reference the template file instead of inline DoD
- Update CLAUDE.md to reference the template file as the canonical DoD
- This ensures the hook, skill, and documentation all use the same checklist

## Example flow

```
# Claude creates a bean without DoD:
beans create "Add login page" -t feature -s todo -d "Create a login page with email/password"

# PostToolUse hook fires, reads stdin JSON with tool_response containing the bean ID
# Hook queries the new bean, finds no DoD items, returns:
{"decision": "block", "reason": "Bean credfolio2-abc1 is missing the Definition of Done checklist. Read .claude/templates/definition-of-done.md and append it to the bean body using beans update."}

# Claude reads the template and updates the bean with the DoD appended
```

## Design decisions

- **Draft beans ARE validated.** The DoD should be present from creation — it's part of the bean's definition, not something added later.
- **Only `beans create` is caught, not GraphQL `createBean` mutations.** The `beans create` CLI is the primary path. Catching mutations would add complexity for minimal benefit.
- **Key phrase matching over exact string matching.** Resilient to minor formatting differences while still catching missing items. Required phrases are derived from the template file, not hardcoded, so updating the template automatically updates validation.
- **The template file is the single source of truth.** The hook reads it at runtime. CLAUDE.md and the dev-workflow skill reference it. This prevents drift.
- **`tool_response` field for bean ID extraction.** The `beans create` output format is `Created credfolio2-xxxx ...` — extract the ID with a regex like `credfolio2-[a-zA-Z0-9]+`.

## Notes

- The hook should be fast — only triggers on `beans create`, exits 0 immediately for all other Bash commands
- The template file approach means credfolio2-xflj (updating references to new subagents) only needs to update one file
- The `beans query --json` output uses `.bean.type` (not `.data.bean.type`) — the existing `start-work.sh` script has a bug here; don't replicate it
- Consider also validating on `beans update` if the DoD gets accidentally removed (optional future enhancement, might be noisy)

## Checklist
- [x] Template file created at `.claude/templates/definition-of-done.md` (content must match CLAUDE.md exactly)
- [x] Hook script created at `.claude/hooks/validate-bean-dod.sh`
- [x] Hook registered in `.claude/settings.json` under `PostToolUse`
- [x] Tested: bean creation without DoD → Claude prompted to add it
- [x] Tested: bean creation with DoD → passes silently
- [x] Tested: epic/milestone creation → skipped (no validation)
- [x] Tested: draft bean creation without DoD → Claude prompted to add it
- [x] Tested: non-beans-create Bash commands → ignored (fast exit)
- [x] dev-workflow skill updated to reference template file
- [x] CLAUDE.md updated to reference template file as canonical DoD

## Definition of Done
- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [ ] Visual verification via `@qa` subagent (via Task tool, for UI changes)
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review
- [ ] Automated code review passed via `@review-backend` and/or `@review-frontend` subagents (via Task tool)