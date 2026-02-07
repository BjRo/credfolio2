---
# credfolio2-xflj
title: Update bean checklists and templates to reference new subagents
status: in-progress
type: task
priority: normal
created_at: 2026-02-07T16:23:26Z
updated_at: 2026-02-07T17:21:13Z
parent: credfolio2-ynmd
---

Update all checklist templates and existing beans to reference the new subagents (QA, review-backend, review-frontend) instead of the old skills.

## Why

Once the review skills become subagents and the QA subagent is created, the "Definition of Done" template and existing beans will reference stale skill names. Claude needs to know to invoke the correct subagents, and the TaskCompleted hook (credfolio2-tdjg) will enforce that checklist items are completed — so the items need to be actionable and accurate.

## What

### 1. Update CLAUDE.md Definition of Done template

**Before:**
```markdown
- [ ] Visual verification with agent-browser (for UI changes)
- [ ] Automated code review passed (`@review-backend` and/or `@review-frontend`)
```

**After (example):**
```markdown
- [ ] Visual verification via QA subagent (for UI changes)
- [ ] Automated code review via review-backend and/or review-frontend subagents
```

### 2. Update dev-workflow skill

The dev-workflow skill generates the Definition of Done checklist when creating beans. Update it to reference the new subagents and provide clear invocation instructions so Claude knows exactly how to satisfy each item.

### 3. Update existing beans

Find all in-progress and todo beans that contain the old references and update them:
- `agent-browser` → QA subagent
- `@review-backend` skill → review-backend subagent
- `@review-frontend` skill → review-frontend subagent

```bash
# Find affected beans
beans query '{ beans(filter: { excludeStatus: ["completed", "scrapped"] }) { id title body } }' --json | jq -r '.data.beans[] | select(.body | test("agent-browser|@review-backend|@review-frontend")) | .id + " " + .title'
```

### 4. Verify invocation clarity

Ensure the checklist items are specific enough that Claude can act on them. Each item should make it obvious how to satisfy it:
- Which subagent to invoke
- What to pass to it (e.g., PR number for reviews, URL for QA)
- What a "pass" looks like

## Dependencies

This task should be done **after** the subagent tasks are complete:
- credfolio2-hpu8 (review-backend subagent)
- credfolio2-fvlb (review-frontend subagent)
- credfolio2-ik4e (QA subagent)

## Implementation Plan

### Approach

This is a documentation-only task: update template text and bean bodies so they reference the `@qa`, `@review-backend`, and `@review-frontend` subagents (invoked via Task tool) instead of the old skill names or raw `agent-browser` commands. Four files and six beans need changes. No source code changes are required.

### Files to Modify

1. **`/workspace/CLAUDE.md`** — Two locations: the "STOP" checklist (already correct) and the "Mandatory Bean Checklist" Definition of Done template (lines 204-213). Also the "Visual Verification with Fixture Resume" section (lines 257-289) needs a note pointing to the `@qa` subagent. Also the "Before Marking Work Complete" subsection within dev-workflow SKILL.md (line 314) references `/skill agent-browser`.
2. **`/workspace/.claude/skills/dev-workflow/SKILL.md`** — Three locations: Step 5 (lines 72-114), the Quick Reference snippet (lines 255-260), the "Mandatory Definition of Done" template (lines 288-297), and the "Before Marking Work Complete" section (lines 301-320).
3. **Six bean files** in `/workspace/.beans/` that contain stale references.

### Checklist

#### A. Update `/workspace/CLAUDE.md`

- [x] **A1. Definition of Done template (lines 204-213):** Change:
  - `- [ ] Visual verification with \`@qa\` subagent (for UI changes)` to `- [ ] Visual verification via \`@qa\` subagent (via Task tool, for UI changes)`
  - `- [ ] Automated code review passed (\`@review-backend\` and/or \`@review-frontend\`)` to `- [ ] Automated code review passed via \`@review-backend\` and/or \`@review-frontend\` subagents (via Task tool)`

- [x] **A2. "Visual Verification with Fixture Resume" section (lines 257-289):** Keep the `agent-browser` reference docs as-is (they document the CLI tool itself), but add a note at the top of the section:
  ```
  > **Note:** For routine visual verification during development, use the `@qa` subagent (via Task tool) instead of running these commands manually. The QA subagent handles dev server management, browser automation, and error checking automatically. The commands below are reference documentation for the underlying `agent-browser` CLI.
  ```

#### B. Fully Rewrite Step 5 of `/workspace/.claude/skills/dev-workflow/SKILL.md`

- [x] **B1. Rewrite Step 5 "Smoke Test with Browser Automation" (lines 72-114):** Replace the entire section with a subagent-based approach, modeled on how Step 8 handles code reviews. The new Step 5 should:
  - Be titled "Visual Verification via QA Subagent" (or similar)
  - Explain that the `@qa` subagent should be launched via the Task tool to keep verbose browser output out of the main context
  - Provide the exact Task tool invocation parameters:
    ```
    Task tool call:
      subagent_type: "qa"
      description: "Visual verification of <feature>"
      prompt: "Verify <feature description>. Start dev servers if needed, navigate to <URL>, test <interactions>, check for errors, and report pass/fail."
    ```
  - Explain what to do after the QA subagent returns (read summary, address failures)
  - Include the "For backend-only changes" note about verifying the API
  - Keep it concise but clear

- [x] **B2. Update Quick Reference snippet (lines 255-260):** Replace the `agent-browser` commands:
  ```
  # After implementation - visual verification
  # Launch @qa subagent via Task tool to verify feature works
  ```

- [x] **B3. Update "Mandatory Definition of Done" template (lines 288-297):** Change:
  - `- [ ] Visual verification with agent-browser (for UI changes)` to `- [ ] Visual verification via \`@qa\` subagent (via Task tool, for UI changes)`
  - `- [ ] Automated code review passed (\`@review-backend\` and/or \`@review-frontend\`)` to `- [ ] Automated code review passed via \`@review-backend\` and/or \`@review-frontend\` subagents (via Task tool)`

- [x] **B4. Update "Before Marking Work Complete" section (lines 301-320):** Change:
  - `# Then use /skill agent-browser to verify` to `# Then launch @qa subagent via Task tool to verify`
  - Update the numbered list to say "Visual verification via `@qa` subagent (for UI changes)" instead of referencing agent-browser or skills

#### C. Update Stale Bean Files

Six beans have stale references. For each, edit the bean's `.md` file directly:

- [x] **C1. `credfolio2-xgdk` — Server-side redirect for root page when profile exists**
  File: `/workspace/.beans/credfolio2-xgdk--server-side-redirect-for-root-page-when-profile-ex.md`
  - Line 32: `- [ ] Visual verification with agent-browser (for UI changes)` → `- [ ] Visual verification via \`@qa\` subagent (via Task tool, for UI changes)`
  - The Definition of Done is also missing the review and "Branch pushed" items — add:
    - `- [ ] Branch pushed and PR created for human review`
    - `- [ ] Automated code review passed via \`@review-backend\` and/or \`@review-frontend\` subagents (via Task tool)`

- [x] **C2. `credfolio2-mzg8` — Use E164 for phone storage and validation**
  File: `/workspace/.beans/credfolio2-mzg8--use-e164-for-phone-storage-and-validation.md`
  - Line 34: `- [ ] Visual verification with agent-browser (for UI changes)` → `- [ ] Visual verification via \`@qa\` subagent (via Task tool, for UI changes)`
  - Add missing: `- [ ] Automated code review passed via \`@review-backend\` and/or \`@review-frontend\` subagents (via Task tool)`

- [x] **C3. `credfolio2-dfpn` — Prevent prompt injection in LLM extraction pipeline**
  File: `/workspace/.beans/credfolio2-dfpn--prevent-prompt-injection-in-llm-extraction-pipeline.md`
  - Definition of Done line: `- [ ] Automated code review passed (\`@review-backend\` and/or \`@review-frontend\`)` → `- [ ] Automated code review passed via \`@review-backend\` and/or \`@review-frontend\` subagents (via Task tool)`

- [x] **C4. `credfolio2-67n7` — Review current backend types, flows, and API design**
  File: `/workspace/.beans/credfolio2-67n7--review-current-backend-types-flows-and-api-design.md`
  - Line 15: `Use the \`@review-backend\` skill to perform the review. This skill provides...` → `Use the \`@review-backend\` subagent (via Task tool) to perform the review. This subagent provides...`
  - Line 35: `- [ ] Run \`@review-backend\` against the backend codebase` → `- [ ] Run \`@review-backend\` subagent (via Task tool) against the backend codebase`

- [x] **C5. `credfolio2-fbrc` — Audit & restructure large skills to use supporting files**
  File: `/workspace/.beans/credfolio2-fbrc--audit-restructure-large-skills-to-use-supporting-f.md`
  - The "Skills to audit" table (lines 22-26) references `review-backend` and `review-frontend` as skills. These are now subagents defined in `.claude/agents/`, not skills. Update the table to remove them and add a note:
    > **Note:** `review-backend` and `review-frontend` have been converted from skills to subagents (`.claude/agents/`). They are no longer skill files and do not need to be audited here.
  - Line 48: Update the "Note" to reflect that the conversion has already happened:
    > `review-backend`, `review-frontend`, and `agent-browser` have already been restructured as part of the subagent migration. `review-backend` and `review-frontend` are now subagent definitions in `.claude/agents/`. The `agent-browser` skill is preloaded by the QA subagent. Focus restructuring efforts on the remaining skills.

- [x] **C6. `credfolio2-tdjg` — Add TaskCompleted hook to enforce bean checklist completion**
  File: `/workspace/.beans/credfolio2-tdjg--add-taskcompleted-hook-to-enforce-bean-checklist-c.md`
  - Line 68 in the "Example flow" section: `- [ ] Visual verification with agent-browser` → `- [ ] Visual verification via \`@qa\` subagent (via Task tool)`

#### D. Verify Consistency

- [x] **D1.** After all edits, run a final grep to confirm no remaining stale references exist in non-completed beans:
  ```bash
  grep -rn "Visual verification with agent-browser\|/skill agent-browser\|/skill review-backend\|/skill review-frontend\|@review-backend.*skill\|@review-frontend.*skill" /workspace/.beans/ /workspace/CLAUDE.md /workspace/.claude/skills/dev-workflow/SKILL.md
  ```
  (References to `agent-browser` in the CLAUDE.md "Visual Verification" reference section and in the `agent-browser` SKILL.md itself are expected and correct — those document the CLI tool, not the workflow.)

- [x] **D2.** Verify that every Definition of Done template across CLAUDE.md and dev-workflow SKILL.md uses identical wording for the visual verification and code review items.

### Testing Strategy

- No automated tests are needed (documentation-only changes)
- Run `pnpm lint` to confirm no build issues
- Run `pnpm test` to confirm nothing is broken
- Manually review each changed file to verify wording consistency

### Open Questions

None — all clarifying questions have been answered by the user.

## Definition of Done
- [x] CLAUDE.md Definition of Done template updated with new subagent references
- [x] dev-workflow skill updated to generate correct checklist items
- [x] Existing in-progress/todo beans updated with new references
- [x] Invocation instructions are clear and actionable
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [ ] Branch pushed and PR created for human review