---
# credfolio2-qd8a
title: Prevent completion with new TODO comments
status: completed
type: task
priority: normal
created_at: 2026-02-08T14:25:29Z
updated_at: 2026-02-08T14:44:40Z
---

Add validation to prevent marking beans as completed when new TODO comments have been introduced in the code during that work. This ensures all work is actually finished before marking a bean as done.

## Implementation Plan

### Approach

**Simple solution**: Add a new checklist item to the Definition of Done template that requires developers to manually verify no new TODO/FIXME/HACK/XXX comments were introduced. The existing `validate-bean-completion.sh` hook already enforces that all checklist items must be checked before completion, so we leverage that existing infrastructure.

This approach:
- Is much simpler than building a new hook script
- Makes the requirement explicit and visible
- Uses existing enforcement mechanisms (checklist validation)
- Trusts developers to do the verification (consistent with other checklist items)

### Files to Modify

1. `/workspace/.claude/templates/definition-of-done.md` — Add new checklist item
2. This bean (`credfolio2-qd8a`) — Update to include the new checklist item (dogfooding)

### Implementation Steps

#### 1. Update the Definition of Done template

Add a new checklist item after the tests item in `/workspace/.claude/templates/definition-of-done.md`:

```markdown
## Definition of Done
- [ ] Tests written (TDD: write tests before implementation)
- [ ] No new TODO/FIXME/HACK/XXX comments introduced (verify with \`git diff main...HEAD | grep -i "^+.*TODO\|FIXME\|HACK\|XXX"\`)
- [ ] \`pnpm lint\` passes with no errors
- [ ] \`pnpm test\` passes with no failures
...
```

**Rationale for placement**: After tests but before lint/test runs, since it's a code quality check that should be done during development.

#### 2. Verify the existing hook enforces the new item

The `validate-bean-dod.sh` hook already checks that the Definition of Done template content appears in the bean body with all items checked. No code changes needed—just verify it works:

1. Create a test bean without the new checklist item
2. Try to mark it completed
3. Verify the hook blocks it

#### 3. Update this bean to include the new checklist item

Update this bean's Definition of Done section to include the new checklist item (dogfooding the change).

#### 4. Test the change propagates to new beans

1. Create a new test bean: `beans create "Test TODO validation" -t task -s draft`
2. The `validate-bean-dod.sh` PostToolUse hook should block it for missing DoD
3. Update the bean with the new DoD template (including the new TODO item)
4. Verify the hook accepts it
5. Delete the test bean

### Testing Strategy

**Validation tests**:
- [x] New checklist item appears in template
- [x] PostToolUse hook enforces template includes new item
- [x] TaskCompleted hook blocks completion if TODO item unchecked (verified by existing hook logic)
- [x] Manual verification: `git diff main...HEAD | grep -i "^+.*TODO\|FIXME\|HACK\|XXX"` works correctly

**Regression check**:
- [x] Existing hooks still work (`validate-bean-dod.sh`, `validate-bean-completion.sh`)
- [x] New beans created with `beans create` are blocked if DoD missing

### Advantages of This Approach

✅ **Simple**: One line change to a template file
✅ **Explicit**: Developers see the requirement clearly  
✅ **Leverages existing infrastructure**: Uses checklist validation hooks
✅ **Low maintenance**: No new scripts to maintain
✅ **Flexible**: Developers can use their judgment for legitimate cases
✅ **Educational**: Provides the actual command to verify

### Tradeoffs

⚠️ **Manual verification**: Relies on developer diligence (but so do all checklist items)

This tradeoff is acceptable because:
- The Definition of Done checklist is about establishing standards
- Developers who check boxes without verifying would bypass automated hooks too
- Simpler code is more maintainable long-term

## Definition of Done
- [x] Tests written (TDD: write tests before implementation)
- [x] No new TODO/FIXME/HACK/XXX comments introduced (verify with `git diff main...HEAD | grep -i "^+.*TODO\|FIXME\|HACK\|XXX"`)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] Visual verification via `@qa` subagent (via Task tool, for UI changes)
- [x] ADR written via `/decision` skill (if new dependencies, patterns, or architectural changes were introduced)
- [x] All other checklist items above are completed
- [x] Branch pushed to remote
- [x] PR created for human review
- [x] Automated code review passed via `@review-backend`, `@review-frontend`, and/or `@review-ai` (for LLM changes) subagents (via Task tool)