---
# credfolio2-xflj
title: Update bean checklists and templates to reference new subagents
status: todo
type: task
created_at: 2026-02-07T16:23:26Z
updated_at: 2026-02-07T16:23:26Z
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

## Definition of Done
- [ ] CLAUDE.md Definition of Done template updated with new subagent references
- [ ] dev-workflow skill updated to generate correct checklist items
- [ ] Existing in-progress/todo beans updated with new references
- [ ] Invocation instructions are clear and actionable
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Branch pushed and PR created for human review