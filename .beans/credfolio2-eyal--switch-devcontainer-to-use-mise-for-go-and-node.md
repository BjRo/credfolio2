---
# credfolio2-eyal
title: Switch devcontainer to use mise for Go and Node
status: completed
type: task
priority: normal
created_at: 2026-02-04T11:42:53Z
updated_at: 2026-02-04T12:01:43Z
---

Replace manual Go/Node installation in Dockerfile with mise for version consistency between local dev and devcontainer.

## Approach
- Switch base image from node:20 to debian:bookworm-slim
- Install mise in Dockerfile
- Use mise to install Go 1.24.1 and Node 20 (from mise.toml)
- Keep manual installation for all other tools (golangci-lint, migrate, beans, etc.)
- Verify devcontainer rebuilds successfully

## Implementation Details
- Changed Dockerfile base from node:20 to debian:bookworm-slim
- Added mise installation as node user
- Copy mise.toml and run `mise install` during build
- Updated npm/corepack/npx commands to use `mise exec --` prefix
- Set MISE_GLOBAL_CONFIG_FILE env var for consistent config location
- Added `mise trust` command to trust the config file
- Changed devcontainer.json build context to `..` (needed for COPY mise.toml)
- Updated COPY paths for firewall scripts to use `.devcontainer/` prefix

## Checklist
- [x] Modify Dockerfile to install mise
- [x] Switch from manual Go installation to mise-managed
- [x] Switch from node:20 base to debian + mise-managed Node
- [x] Test devcontainer rebuild
- [x] Verify go and node versions are correct (Go 1.24.1, Node v20.20.0)
- [x] Verify pnpm build works
- [x] Verify pnpm dev works
- [x] Run tests (246 frontend tests, 15 backend packages - all passing)
- [x] Update CLAUDE.md with mise details

## Definition of Done
- [x] Tests written (TDD: write tests before implementation) - N/A for infrastructure change
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures (246 frontend tests, 15 backend packages)
- [x] Visual verification with agent-browser (for UI changes) - N/A for infrastructure change
- [x] All other checklist items above are completed
- [x] Branch pushed and PR created for human review