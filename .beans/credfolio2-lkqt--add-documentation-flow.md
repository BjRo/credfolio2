---
# credfolio2-lkqt
title: add-documentation-flow
status: completed
type: task
priority: normal
created_at: 2026-01-19T16:19:33Z
updated_at: 2026-01-19T16:40:34Z
---

## Goal

Continuously maintain documentation for technical and process decisions. This documentation serves:
- **Agents**: Quick access to important decisions for context
- **Humans**: Fast onboarding and understanding of project evolution

## Triggers for Documentation

Document when:
1. Modifying the tech stack (adding/removing dependencies, frameworks, tools)
2. Introducing new concepts into the codebase
3. Deprecating existing patterns or approaches
4. Making significant architectural decisions

## Implementation

### 1. Directory Structure

Create `/decisions/` at repository root with timestamped markdown files:

```
/decisions/
├── 20260119120000-add-turborepo-for-monorepo-builds.md
├── 20260119130000-remove-google-fonts-for-network-isolation.md
└── ...
```

**Naming convention**: `YYYYMMDDHHMMSS-kebab-case-title.md`

### 2. Decision File Template

Each decision file must contain:

```markdown
# [Title: What was done]

**Date**: YYYY-MM-DD
**Bean**: [bean-id]

## Context

[What situation led to this decision?]

## Decision

[What was decided and implemented?]

## Reasoning

[Why was this approach chosen over alternatives?]

## Consequences

[What are the implications? What changes for the codebase going forward?]
```

### 3. Claude Code Skill

Create a `/decision` skill that:
- Prompts for decision details (or infers from recent work)
- Generates the timestamped filename
- Creates the decision file with the template filled in
- Reminds to commit the decision file with related changes

Location: `.claude/skills/decision/SKILL.md`

### 4. CLAUDE.md Integration

Add a section to CLAUDE.md instructing Claude to:
- Consider documenting decisions after significant changes
- Reference the `/decision` skill
- Include decision files in commits alongside code changes

## Checklist

- [x] Create `/decisions/` directory with a README explaining the format
- [x] Create the decision skill at `.claude/skills/decision/SKILL.md`
- [x] Update `CLAUDE.md` with documentation workflow guidance
- [x] Create an initial decision documenting this documentation system itself
- [x] Test the skill by running `/decision`

