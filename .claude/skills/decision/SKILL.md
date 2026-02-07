---
name: decision
description: Document a technical or process decision. Use after making significant changes to the tech stack, introducing new patterns, or deprecating existing approaches.
---

# Decision Documentation

You are documenting a technical or process decision for this project.

## When This Skill Applies

Use this skill when:
- Adding or removing dependencies, frameworks, or tools
- Introducing new architectural patterns or concepts
- Deprecating existing approaches
- Making significant technical choices that affect the codebase

If none of the above apply to the work you just completed, skip this skill and mark the ADR checklist item in the Definition of Done as N/A.

## When This Skill Does NOT Apply

Do NOT create a decision record for:
- Bug fixes that use existing patterns and dependencies
- New features that follow established conventions without introducing new patterns
- Routine refactoring that does not change design or architecture
- Test additions or improvements
- Documentation updates (unless documenting a new documentation strategy)
- Configuration changes that follow existing patterns (e.g., adding a new route using the existing router)

## Process

### 1. Gather Information

If the user didn't provide details, ask about:
- **What was done?** (the decision/change)
- **Why was it done?** (the reasoning)
- **What bean introduced this?** (if applicable)

If you just completed work on a bean, you can infer most of this from context.

### 2. Generate the Decision File

Create the file in `/decisions/` with this naming convention:

```
YYYYMMDDHHMMSS-kebab-case-title.md
```

Use the current timestamp. Convert the title to kebab-case (lowercase, hyphens instead of spaces).

### 3. Use This Template

```markdown
# [Title: What was done]

**Date**: YYYY-MM-DD
**Bean**: [bean-id or "N/A" if not applicable]

## Context

[What situation led to this decision? What problem were we solving?]

## Decision

[What was decided and implemented? Be specific about what changed.]

## Reasoning

[Why was this approach chosen? What alternatives were considered?]

## Consequences

[What are the implications? What changes for the codebase going forward?
What should future developers/agents know about this decision?]
```

### 4. Update the Decision Index

**MANDATORY** — You MUST update `/decisions/CLAUDE.md` immediately after creating the decision file. Add a new row to the index table:

```markdown
| [filename.md](filename.md) | Title | YYYY-MM-DD | Brief one-line summary |
```

**Do NOT consider the decision complete until the index is updated.** The decision file and the index update must be committed together.

### 5. Commit Both Files

Commit the decision file AND the updated index together:

```bash
git add decisions/<new-decision-file>.md decisions/CLAUDE.md
git commit --no-gpg-sign -m "docs: Add ADR for <title>

Co-Authored-By: Claude <noreply@anthropic.com>"
```

If committing alongside related code changes, include both decision files in that commit.

## Example

If the user says `/decision` after adding Redis caching:

1. Generate filename: `20260119143022-add-redis-caching-layer.md`
2. Fill in template based on the work done
3. Create file in `/decisions/`
4. Remind to commit with related changes

## Tips

- Keep decisions focused on one topic
- Be specific about what changed and why
- Include enough context for someone unfamiliar with the project
- Reference the bean ID so decisions can be traced back to work items
