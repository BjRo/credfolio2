# Credfolio2 - Project Context for Claude Code

## STOP — Before Marking Any Work Complete

**You MUST complete ALL of these steps before marking a bean as completed or telling the user you're done:**

1. **Run lint**: `pnpm lint` — fix all errors
2. **Run tests**: `pnpm test` — all tests must pass
3. **Visual verification** (for ANY UI changes): Use `/skill agent-browser` to verify the feature works in the browser
4. **Bean checklist**: Ensure ALL checklist items in the bean are checked off

**DO NOT skip these steps.**
**DO NOT say "you can run tests to verify" — run them yourself.**
**DO NOT mark a bean complete if it has unchecked checklist items.**

---

## Directory Structure

```
/workspace/
├── .devcontainer/          # Dev container (Node 20 + Go 1.24.1)
├── .claude/                # Claude Code settings (Write/Edit permissions enabled)
├── decisions/              # Architecture Decision Records (ADRs)
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

### Database

- PostgreSQL 16 (via docker-compose)
- Database names: `credfolio_dev` (default), `credfolio_test`
- Environment selection via `CREDFOLIO_ENV` (defaults to `dev`)
- Migrations: golang-migrate with timestamp versioning

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

# Database (run from host, not devcontainer)
docker-compose up -d           # Start PostgreSQL and MinIO
docker-compose down            # Stop services

# Migrations (run from src/backend/)
cd src/backend
make help                      # Show all migration commands
make migration name=create_users  # Create new migration
make migrate-up                # Run pending migrations (on credfolio_dev)
make migrate-down              # Rollback last migration
CREDFOLIO_ENV=test make migrate-up   # Run migrations on test database
make migrate-status            # Show migration status for all environments
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
- Backend config: `src/backend/internal/config/config.go`
- Backend Makefile: `src/backend/Makefile`
- Migrations: `src/backend/migrations/`
- Frontend layout: `src/frontend/src/app/layout.tsx`
- Frontend homepage: `src/frontend/src/app/page.tsx`
- Next.js config: `src/frontend/next.config.ts`
- Turborepo config: `turbo.json`
- Workspace definition: `pnpm-workspace.yaml`
- Docker services: `docker-compose.yml`

## Development Workflow

This project follows a strict workflow for all feature development. See `.claude/skills/dev-workflow/SKILL.md` for full details.

### Key Principles

1. **Feature Branches**: Never commit directly to main
2. **TDD**: Write tests before implementation (see `/skill tdd`)
3. **PR Review**: All changes require human review before merge
4. **Beans Tracking**: Use beans to track all work (see `/skill issue-tracking-with-beans`)

### Mandatory Bean Checklist Items

**Every bean you create MUST include a "Definition of Done" checklist section.** This section should be at the end of the bean body and include these mandatory items:

```markdown
## Definition of Done
- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Visual verification with agent-browser (for UI changes)
- [ ] All other checklist items above are completed
```

You cannot mark a bean as completed while it has unchecked items. This structurally enforces the workflow compliance.

### Quick Reference

```bash
# Start work on a bean
git checkout main && git pull origin main
git checkout -b feat/<bean-id>-<description>
beans update <bean-id> --status in-progress

# Work using TDD, update bean checklist, commit with bean file

# When done
git push -u origin <branch>
gh pr create --title "..." --body "Closes beans-<id>"
# WAIT for human review - do NOT merge yourself

# After human merges
git checkout main && git pull origin main
beans update <bean-id> --status completed
```

## Git Details

- Main branch: `main`
- Remote: `github.com:BjRo/credfolio2.git`
- Commits use `--no-gpg-sign` flag
- Co-authored by: `Claude <noreply@anthropic.com>`

## Devcontainer Notes

- Based on `node:20` image
- Go 1.24.1 installed via wget during build
- Includes: git, gh, fzf, zsh, and development tools
- Rebuild required when Dockerfile changes

## Decision Documentation

This project maintains Architecture Decision Records (ADRs) in `/decisions/`.

### When to Document

After completing work that involves:
- Adding or removing dependencies, frameworks, or tools
- Introducing new architectural patterns or concepts
- Deprecating existing approaches
- Making significant technical choices

### How to Document

Use the `/decision` skill to create a new decision record:

```
/decision
```

This generates a timestamped file in `/decisions/` with the standard template.

### Important

- Include decision files in commits alongside related code changes
- Reference the bean ID that introduced the decision
- See `/decisions/README.md` for the full template and guidelines
