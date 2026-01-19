# Add Decision Documentation Workflow

**Date**: 2026-01-19
**Bean**: workspace-lkqt

## Context

As the project grows, technical decisions accumulate but their rationale gets lost. Both human developers and AI agents need quick access to understand why certain choices were made, what patterns are in use, and what approaches have been deprecated.

Without documented decisions:
- New contributors spend time rediscovering context
- Agents may suggest approaches that were previously rejected
- The "why" behind technical choices fades over time

## Decision

Implemented an Architecture Decision Record (ADR) system with:

1. **Directory structure**: `/decisions/` at repository root
2. **Naming convention**: `YYYYMMDDHHMMSS-kebab-case-title.md` (Rails migration style)
3. **Standard template** including: Context, Decision, Reasoning, Consequences, and Bean reference
4. **Claude Code skill** (`/decision`) for easy creation
5. **CLAUDE.md integration** with workflow guidance

## Reasoning

- **Timestamped files**: Provides natural chronological ordering and avoids naming conflicts
- **Rails-style naming**: Familiar convention, sorts correctly in file explorers
- **Skill-based approach**: Makes it easy for agents to create decisions consistently
- **Bean references**: Links decisions back to work items for traceability
- **Minimal template**: Captures essential information without being burdensome

Alternatives considered:
- Wiki-based documentation: Rejected because it lives outside the repo and can drift
- Single DECISIONS.md file: Rejected because it becomes unwieldy and hard to navigate
- Numbered ADRs (ADR-001): Rejected in favor of timestamps which are more flexible

## Consequences

- All significant technical decisions should now be documented using `/decision`
- Decision files should be committed alongside the code changes they describe
- The `/decisions/` directory will grow over time as a project history
- Agents should check existing decisions before proposing changes to avoid revisiting settled questions
- CLAUDE.md now references this workflow, making it visible to future sessions
