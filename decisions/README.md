# Decisions

This directory contains Architecture Decision Records (ADRs) documenting significant technical and process decisions for this project.

## Purpose

- **For agents**: Quick context on key decisions and their rationale
- **For humans**: Understand project evolution and the "why" behind choices

## When to Document

Create a decision record when:

1. Modifying the tech stack (adding/removing dependencies, frameworks, tools)
2. Introducing new concepts or patterns into the codebase
3. Deprecating existing patterns or approaches
4. Making significant architectural decisions

## File Naming

Files use timestamped names following Rails migration conventions:

```
YYYYMMDDHHMMSS-kebab-case-title.md
```

Example: `20260119120000-add-turborepo-for-monorepo-builds.md`

## Template

Each decision file should contain:

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

## Creating Decisions

Use the `/decision` skill in Claude Code to create new decision records:

```
/decision
```

This will guide you through creating a properly formatted decision file.
