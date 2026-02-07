# Credfolio2 - Project Context for Claude Code

## STOP — Before Marking Any Work Complete

**You MUST complete ALL of these steps before marking a bean as completed or telling the user you're done:**

1. **Feature branch**: You MUST be on a feature branch, NOT main
2. **Run lint**: `pnpm lint` — fix all errors
3. **Run tests**: `pnpm test` — all tests must pass
4. **Visual verification** (for ANY UI changes): Use the `@qa` subagent (via Task tool) to verify the feature works in the browser
5. **Bean checklist**: Ensure ALL checklist items in the bean are checked off
6. **Push and create PR**: Push your branch and create a PR for human review
7. **Run automated code reviews**: Launch `@review-backend` and/or `@review-frontend` as subagents via Task tool — address any critical findings

**DO NOT skip these steps.**
**DO NOT commit directly to main.**
**DO NOT say "you can run tests to verify" — run them yourself.**
**DO NOT mark a bean complete if it has unchecked checklist items.**
**DO NOT merge your own PR — wait for human review.**

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

## Key Files (auto-imported)

@src/frontend/src/app/layout.tsx
@src/frontend/next.config.ts
@src/frontend/codegen.ts
@src/frontend/src/lib/urql/client.ts
@src/backend/gqlgen.yml
@src/backend/internal/domain/repository.go

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

**Backend (Go)**
- Entry point: `src/backend/cmd/server/main.go` -- HTTP server setup with chi router, middleware (logger, recoverer, CORS), route registration (REST + GraphQL), LLM provider initialization, River job queue setup
- Config: `src/backend/internal/config/config.go` -- Loads all env vars (database, MinIO, LLM providers, queue); see struct fields for available env var names
- Domain entities: `src/backend/internal/domain/entities.go` -- Bun ORM structs (User, File, ReferenceLetter, Author, Testimonial, SkillValidation, ExperienceValidation)
- Domain profile: `src/backend/internal/domain/profile.go` -- Profile, ProfileExperience, ProfileEducation, ProfileSkill entities
- Repository interfaces: `src/backend/internal/domain/repository.go` -- All repository interfaces (UserRepository, FileRepository, etc.)
- Repository implementations: `src/backend/internal/repository/postgres/` -- One file per repository (e.g., `resume_repository.go`, `profile_skill_repository.go`)
- GraphQL schema: `src/backend/internal/graphql/schema/schema.graphqls` -- All type definitions, queries, mutations (1591 lines)
- GraphQL resolvers: `src/backend/internal/graphql/resolver/schema.resolvers.go` -- All query/mutation implementations (auto-generated method stubs, manually implemented bodies)
- GraphQL converter: `src/backend/internal/graphql/resolver/converter.go` -- Domain-to-GraphQL model mapping functions
- GraphQL handler: `src/backend/internal/graphql/handler.go` -- gqlgen server setup and dependency injection
- GraphQL config: `src/backend/gqlgen.yml` -- Code generation config (schema paths, output locations, model mappings)
- Background jobs: `src/backend/internal/job/` -- River queue workers (document_detection, document_processing, reference_letter_processing, resume_processing)
- LLM infrastructure: `src/backend/internal/infrastructure/llm/` -- Anthropic/OpenAI providers, document extraction, resilience wrappers, prompts
- Materialization service: `src/backend/internal/service/materialization.go` -- Converts extracted resume/letter data into profile entities
- Migrations: `src/backend/migrations/` -- Timestamped up/down SQL files (golang-migrate format)
- Makefile: `src/backend/Makefile` -- Migration commands (make migrate-up, make migration name=X)

**Frontend (Next.js/React)**
- Layout: `src/frontend/src/app/layout.tsx` -- Root layout with ThemeProvider, UrqlProvider, SiteHeader
- Pages: `src/frontend/src/app/` -- App Router pages: `/` (home), `/upload` (document upload flow), `/upload-resume` (resume upload), `/profile/[id]` (profile view), `/viewer` (PDF viewer)
- Components by feature: `src/frontend/src/components/profile/` -- Profile display/edit components (ProfileHeader, EducationSection, SkillsSection, TestimonialsSection, etc.)
- Upload flow: `src/frontend/src/components/upload/` -- Multi-step upload wizard (DocumentUpload, DetectionProgress, ExtractionProgress, ExtractionReview)
- UI primitives: `src/frontend/src/components/ui/` -- shadcn/ui components (button, dialog, input, badge, etc.)
- GraphQL queries: `src/frontend/src/graphql/queries.graphql` -- All frontend query definitions
- GraphQL mutations: `src/frontend/src/graphql/mutations.graphql` -- All frontend mutation definitions
- Generated types: `src/frontend/src/graphql/generated/` -- graphql-codegen output (types + typed document nodes)
- urql client: `src/frontend/src/lib/urql/client.ts` -- GraphQL client setup, endpoint configuration
- urql provider: `src/frontend/src/lib/urql/provider.tsx` -- React context provider for urql
- Codegen config: `src/frontend/codegen.ts` -- graphql-codegen configuration (schema source, output, scalar mappings)
- Test setup: `src/frontend/vitest.config.ts` -- Vitest config with happy-dom, path aliases, @urql/next mock
- Test mocks: `src/frontend/src/test/mocks/` -- Test mock for @urql/next

**Project-level**
- Turborepo config: `turbo.json`
- Workspace definition: `pnpm-workspace.yaml`
- Docker services: `docker-compose.yml`
- Next.js config: `src/frontend/next.config.ts` -- Proxy rewrites for GraphQL and MinIO storage

## Key Patterns

### Data Flow: GraphQL Request

1. Frontend sends query/mutation via urql client to `/api/graphql` (proxied by Next.js to `localhost:8080/graphql`)
2. gqlgen routes to resolver method in `schema.resolvers.go`
3. Resolver calls domain repository interfaces (defined in `domain/repository.go`)
4. Repository implementations in `internal/repository/postgres/` execute SQL via Bun ORM
5. Resolver uses converter functions (`converter.go`) to map domain entities to GraphQL models
6. Response flows back through urql to the React component

### Data Flow: Document Upload and Extraction

1. Frontend uploads file via GraphQL `uploadDocument` mutation
2. Backend stores file in MinIO (object storage), creates File record, enqueues River job
3. Detection worker (`job/document_detection.go`) classifies document type using LLM
4. Processing worker (`job/document_processing.go` or `job/reference_letter_processing.go`) extracts structured data using LLM
5. Materialization service (`service/materialization.go`) converts extracted data into Profile entities (experiences, education, skills, testimonials)
6. Frontend polls for status updates via GraphQL queries

### GraphQL Code Generation

- **Backend**: gqlgen (Go) -- Schema in `internal/graphql/schema/schema.graphqls`, config in `gqlgen.yml`. Run `go generate ./internal/graphql/` to regenerate. Produces `generated/generated.go` (executable schema) and `model/models_gen.go` (Go types). Resolver stubs in `resolver/schema.resolvers.go`.
- **Frontend**: graphql-codegen (TypeScript) -- Config in `codegen.ts`, reads schema from backend path. Run `pnpm codegen` in frontend. Produces typed document nodes in `src/graphql/generated/`.

### Domain Layer Structure

- `internal/domain/` contains pure business types and interfaces (no infrastructure dependencies)
- `internal/domain/entities.go` -- Core entities with Bun ORM tags (User, File, ReferenceLetter, etc.)
- `internal/domain/profile.go` -- Profile aggregate (Profile, ProfileExperience, ProfileEducation, ProfileSkill)
- `internal/domain/repository.go` -- Repository interfaces (one per entity)
- `internal/domain/llm.go` -- LLM provider and extractor interfaces
- `internal/domain/storage.go` -- Object storage interface
- `internal/domain/job.go` -- Job enqueuer interface

### Test Organization

- **Go tests**: Co-located with source (`*_test.go` files next to implementation). Repository tests use real PostgreSQL (`credfolio_test` database). Job/service tests use mocks.
- **Frontend tests**: Co-located with source (`*.test.tsx` / `*.test.ts`). Use vitest with happy-dom. urql mocked via `src/test/mocks/urql-next.tsx` (aliased in vitest.config.ts). Run with `pnpm test` from frontend or root.

### Migration Conventions

- Timestamp-prefixed pairs: `YYYYMMDDHHMMSS_name.up.sql` / `YYYYMMDDHHMMSS_name.down.sql`
- Created via `make migration name=X` in `src/backend/`
- Applied via `make migrate-up` (dev) or `CREDFOLIO_ENV=test make migrate-up` (test)
- Tool: golang-migrate CLI

## Development Workflow

This project follows a strict workflow. See `/skill dev-workflow` for the full process.

### Mandatory Bean Checklist

**Every bean MUST include a "Definition of Done" section.** A PostToolUse hook automatically validates this on `beans create` commands.

@.claude/templates/definition-of-done.md

You cannot mark a bean as completed while it has unchecked items.

## Git Details

- Main branch: `main`
- Remote: `github.com:BjRo/credfolio2.git`
- Commits use `--no-gpg-sign` flag
- Co-authored by: `Claude <noreply@anthropic.com>`

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
