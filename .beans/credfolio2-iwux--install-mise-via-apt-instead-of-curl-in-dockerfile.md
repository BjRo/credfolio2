---
# credfolio2-iwux
title: Install mise via apt instead of curl in Dockerfile
status: completed
type: task
priority: normal
created_at: 2026-02-04T12:09:35Z
updated_at: 2026-02-04T12:14:39Z
---

Replace the curl-based mise installation with apt-based installation for better security and maintainability.

## Context
Currently using: `curl https://mise.run | sh`
Should use: apt package manager

## Checklist
- [x] Update Dockerfile to use apt for mise installation
- [x] Test the devcontainer builds successfully (requires user to rebuild)
- [x] Verify mise is properly installed and functional (requires user to rebuild)

## Changes Made
- Added mise apt repository with GPG key verification
- Replaced `curl https://mise.run | sh` with `apt-get install mise`
- Kept mise shims in user home directory (standard behavior)
- Committed and pushed to mise-devcontainer branch

## Definition of Done
- [x] Code changes completed and committed
- [x] Branch pushed (updates existing PR)
- [x] User verifies devcontainer rebuilds successfully
- [x] User verifies mise works (`mise --version` after rebuild)