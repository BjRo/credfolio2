# Credfolio2 - Project Context for Claude Code

## Directory Structure

```
/workspace/
├── .devcontainer/          # Dev container (Node 20 + Go 1.24.1)
├── .claude/                # Claude Code settings (Write/Edit permissions enabled)
├── src/
│   ├── frontend/           # Next.js 16 app (TypeScript, Tailwind CSS 4, React 19)
│   │   ├── src/app/        # Next.js app directory structure
│   │   └── package.json    # Has "backend": "workspace:*" dependency
│   └── backend/            # Go 1.24 backend
│       ├── cmd/server/main.go  # HTTP server on :8080
│       └── go.mod
├── turbo.json              # Turborepo pipeline config
├── pnpm-workspace.yaml     # Defines: src/frontend, src/backend
└── package.json            # Root with Turborepo scripts
```

## Key Technical Decisions

### Package Manager

- **pnpm 10.28.1** (not npm/yarn)
- Configured via `packageManager` field in package.json
- Uses pnpm workspaces for monorepo

### Build System

- **Turborepo 2.7.5** orchestrates builds
- Build order: backend FIRST, then frontend (enforced via workspace dependency)
- Command: `pnpm build` (builds everything in correct order)
- Caches in `.turbo/` (gitignored)

### Frontend Stack

- Next.js 16 with App Router
- TypeScript
- Tailwind CSS 4
- React 19
- SWC compiler (Go-based, built into Next.js)
- **NO Google Fonts** (removed due to network restrictions in devcontainer)

### Backend Stack

- Go 1.24.1
- Standard library HTTP server
- Runs on port 8080
- Routes: `/` (hello), `/health` (health check)

## Common Commands

```bash
# Build everything (backend → frontend)
pnpm build

# Dev mode (both services)
pnpm dev

# Individual package commands
cd src/frontend && pnpm dev    # Next.js on :3000
cd src/backend && pnpm dev     # Go server on :8080
cd src/backend && pnpm build   # Compiles to bin/server

# Cleanup
pnpm clean                     # Clean all packages
rm -rf .turbo                  # Clear Turborepo cache
```

## Important Context

### Permissions

- Claude Code has `Write(*)` and `Edit(*)` permissions enabled in `.claude/settings.local.json`
- Running in devcontainer provides isolation from host machine
- Can freely modify files in `/workspace/`

### Build Requirements

- Frontend build depends on backend build completing first
- This is enforced via `"backend": "workspace:*"` in frontend's devDependencies
- Turborepo's `^build` notation in turbo.json respects this dependency

### Network Restrictions

- Devcontainer may have limited external network access
- Google Fonts were removed from layout.tsx for this reason
- Be aware when adding external resource dependencies

### File Locations to Remember

- Backend entry point: `src/backend/cmd/server/main.go`
- Frontend layout: `src/frontend/src/app/layout.tsx`
- Frontend homepage: `src/frontend/src/app/page.tsx`
- Next.js config: `src/frontend/next.config.ts`
- Turborepo config: `turbo.json`
- Workspace definition: `pnpm-workspace.yaml`

## Git Workflow

- Main branch: `main`
- Remote: `github.com:BjRo/credfolio2.git`
- Commits use `--no-gpg-sign` flag
- Co-authored by: `Claude Sonnet 4.5 <noreply@anthropic.com>`

## Devcontainer Notes

- Based on `node:20` image
- Go 1.24.1 installed via wget during build
- Includes: git, gh, fzf, zsh, and development tools
- Rebuild required when Dockerfile changes
