---
# credfolio2-ntq2
title: Auto-trust mise.toml in devcontainer setup
status: todo
type: bug
priority: normal
created_at: 2026-02-05T13:42:37Z
updated_at: 2026-02-05T14:07:01Z
parent: credfolio2-2ex3
---

## Problem

When running `pnpm dev` for the first time in the devcontainer, mise doesn't trust the root `mise.toml` file. This requires manually running trust commands, wasting tokens and developer time.

## Solution

Configure mise to trust the project's `mise.toml` during the devcontainer build or setup phase, so tools are available immediately without manual intervention.

## Checklist
- [ ] Identify the right approach (e.g., `mise trust` in post-create, or mise config to auto-trust the workspace)
- [ ] Add the trust step to the devcontainer setup
- [ ] Verify `pnpm dev` works without mise trust prompts on fresh container

## Definition of Done
- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review