---
# credfolio2-kdyr
title: Install corepack during devcontainer build
status: in-progress
type: bug
priority: normal
created_at: 2026-02-05T13:42:31Z
updated_at: 2026-02-05T14:50:01Z
parent: credfolio2-2ex3
---

## Problem

When running `pnpm dev` for the first time in the devcontainer, corepack prompts to download itself. This wastes tokens and time during development.

## Solution

Install/enable corepack during the container build phase (in the Dockerfile or a post-create script) so it's ready to use immediately.

## Checklist
- [x] Identify where corepack should be installed (Dockerfile or devcontainer lifecycle script)
- [x] Add corepack enable/install step
- [ ] Verify `pnpm dev` runs without corepack prompt on fresh container

## Definition of Done
- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review