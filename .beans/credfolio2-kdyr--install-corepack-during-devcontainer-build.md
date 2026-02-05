---
# credfolio2-kdyr
title: Install corepack during devcontainer build
status: todo
type: bug
priority: normal
created_at: 2026-02-05T13:42:31Z
updated_at: 2026-02-05T14:06:37Z
parent: credfolio2-2ex3
---

## Problem

When running `pnpm dev` for the first time in the devcontainer, corepack prompts to download itself. This wastes tokens and time during development.

## Solution

Install/enable corepack during the container build phase (in the Dockerfile or a post-create script) so it's ready to use immediately.

## Checklist
- [ ] Identify where corepack should be installed (Dockerfile or devcontainer lifecycle script)
- [ ] Add corepack enable/install step
- [ ] Verify `pnpm dev` runs without corepack prompt on fresh container

## Definition of Done
- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review