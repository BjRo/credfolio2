---
# credfolio2-iwux
title: Install mise via apt instead of curl in Dockerfile
status: in-progress
type: task
created_at: 2026-02-04T12:09:35Z
updated_at: 2026-02-04T12:09:35Z
---

Replace the curl-based mise installation with apt-based installation for better security and maintainability.

## Context
Currently using: `curl https://mise.run | sh`
Should use: apt package manager

## Checklist
- [x] Update Dockerfile to use apt for mise installation
- [ ] Test the devcontainer builds successfully
- [ ] Verify mise is properly installed and functional

## Definition of Done
- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Visual verification with agent-browser (for UI changes)
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review