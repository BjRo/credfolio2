---
# credfolio2-ntq2
title: Auto-trust mise.toml in devcontainer setup
status: completed
type: bug
priority: normal
created_at: 2026-02-05T13:42:37Z
updated_at: 2026-02-05T14:48:36Z
parent: credfolio2-2ex3
---

## Problem

When running `pnpm dev` for the first time in the devcontainer, mise doesn't trust the root `mise.toml` file. This requires manually running trust commands, wasting tokens and developer time.

## Solution

Configure mise to trust the project's `mise.toml` during the devcontainer build or setup phase, so tools are available immediately without manual intervention.

## Checklist
- [x] Identify the right approach (e.g., `mise trust` in post-create, or mise config to auto-trust the workspace)
- [x] Add the trust step to the devcontainer setup
- [x] Verify `pnpm dev` works without mise trust prompts on fresh container

## Definition of Done
- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review