---
# credfolio2-c026
title: start-work.sh fails to query bean metadata
status: todo
type: bug
priority: normal
created_at: 2026-02-07T18:43:13Z
updated_at: 2026-02-07T18:43:13Z
---

## Problem

When running `.claude/scripts/start-work.sh credfolio2-o624`, the script fails at step 2/5 with:

```
[0;31mError: Bean 'credfolio2-o624' not found[0m
```

The bean exists and can be queried directly via `beans query`, so the issue is in how the script queries bean metadata.

## Likely Cause

The script may be using an incorrect GraphQL query path (e.g., `.data.bean.type` instead of `.bean.type`) or passing the bean ID in an unexpected format to the `beans query` command.

The bean body for credfolio2-o624 mentions: "The `beans query --json` output uses `.bean.type` (not `.data.bean.type`) — the existing `start-work.sh` script has a bug here; don't replicate it"

## Steps to Reproduce

```bash
.claude/scripts/start-work.sh credfolio2-o624
```

## Expected

Script should query the bean successfully and create a feature branch.

## Actual

Script errors with "Bean 'credfolio2-o624' not found" at step 2/5.

## Definition of Done
- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Visual verification via `@qa` subagent (via Task tool, for UI changes)
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review
- [ ] Automated code review passed via `@review-backend` and/or `@review-frontend` subagents (via Task tool)