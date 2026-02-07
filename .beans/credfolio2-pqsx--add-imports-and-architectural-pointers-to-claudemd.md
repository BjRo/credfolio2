---
# credfolio2-pqsx
title: Add @imports and architectural pointers to CLAUDE.md
status: completed
type: task
priority: normal
created_at: 2026-02-07T16:17:24Z
updated_at: 2026-02-07T20:27:02Z
parent: credfolio2-ynmd
---

Enrich CLAUDE.md with @imports of key files and detailed "where to find X" pointers so Claude can jump directly to the right place instead of exploring.

## Why

Every session, Claude spends early turns grepping and globbing to locate key files — config, schema, resolvers, API client, routes. This is wasted context and time. By importing key files and adding rich architectural pointers, Claude starts with a mental map of the codebase.

## What

### 1. Add @imports for frequently-needed files

Identify files Claude reaches for in almost every session and import them directly. These get loaded at session start:

```markdown
@src/backend/graph/schema.graphqls
@src/backend/internal/config/config.go
@src/frontend/src/app/layout.tsx
@src/frontend/next.config.ts
```

Be selective — each import adds to context. Only import files that are:
- Referenced in most sessions (not niche)
- Relatively small (not 500+ lines)
- Structurally important (schema, config, routing)

### 2. Enrich "File Locations to Remember" section

The current section lists paths but not what's inside them. Add brief descriptions of what each file/directory contains and how it relates to the architecture:

**Before:**
```markdown
- Backend entry point: `src/backend/cmd/server/main.go`
```

**After:**
```markdown
- Backend entry point: `src/backend/cmd/server/main.go` — HTTP server setup, middleware chain, route registration
- GraphQL schema: `src/backend/graph/schema.graphqls` — all type definitions, queries, mutations
- Resolvers: `src/backend/graph/resolver/` — one file per domain (resume.go, education.go, etc.)
- Database models: `src/backend/internal/models/` — Bun ORM structs mapping to PostgreSQL tables
- Frontend API layer: `src/frontend/src/lib/graphql/` — urql client setup + query/mutation definitions
- Frontend pages: `src/frontend/src/app/` — Next.js App Router pages and layouts
- Components: `src/frontend/src/components/` — shared React components
```

### 3. Add a "Key Patterns" section

Document the high-level patterns so Claude doesn't need to infer them:
- How GraphQL resolvers connect to the database layer
- How frontend pages fetch data (urql + server components vs client components)
- How migrations are structured
- How tests are organized (Go test files, frontend test setup)

## Approach

1. Audit recent session transcripts (if available) to see what Claude explores most
2. Identify the top 5-10 files/directories Claude gravitates to
3. Add @imports for the most critical ones
4. Enrich the file locations section with architectural context
5. Keep it concise — this is a map, not documentation

## Note

This task should be done alongside or after credfolio2-zsm2 (CLAUDE.md refactor into modular rules). The refactor slims down CLAUDE.md; this task enriches what remains with higher-signal content.

---

## Implementation Plan

### Approach

Enrich CLAUDE.md with @imports for small, structurally important files and add two new sections: an enriched "File Locations to Remember" section with accurate paths and descriptions, and a new "Key Patterns" section documenting architectural patterns. The bean's original path suggestions contain several inaccuracies (discovered during codebase exploration), which this plan corrects.

**Key corrections from the original bean spec:**
- `src/backend/graph/schema.graphqls` does not exist. Correct path: `src/backend/internal/graphql/schema/schema.graphqls` (1591 lines -- too large for @import)
- `src/backend/internal/models/` does not exist. Domain entities live in `src/backend/internal/domain/`
- `src/frontend/src/lib/graphql/` does not exist. GraphQL queries/mutations are at `src/frontend/src/graphql/`; the urql client is at `src/frontend/src/lib/urql/`
- `src/backend/graph/resolver/` does not exist. Resolvers are at `src/backend/internal/graphql/resolver/`

### @import Candidates (size analysis)

| File | Lines | Import? | Rationale |
|------|-------|---------|-----------|
| `src/backend/internal/graphql/schema/schema.graphqls` | 1591 | NO | Too large (500+ lines). Would consume excessive context. |
| `src/backend/internal/config/config.go` | 298 | NO | Borderline large; env var names and defaults are useful but 298 lines is significant context. Better referenced as a pointer. |
| `src/frontend/src/app/layout.tsx` | 27 | YES | Small, shows app-level providers (ThemeProvider, UrqlProvider, SiteHeader). |
| `src/frontend/next.config.ts` | 40 | YES | Small, shows proxy rewrites (critical for understanding API routing). |
| `src/backend/gqlgen.yml` | 105 | YES | Small, shows GraphQL code generation config (schema location, resolver layout, model mappings). |
| `src/backend/internal/domain/repository.go` | 159 | YES | Medium, defines all repository interfaces -- central to understanding data access patterns. |
| `src/frontend/codegen.ts` | 23 | YES | Very small, shows frontend GraphQL codegen config and schema source. |
| `src/frontend/src/lib/urql/client.ts` | 24 | YES | Very small, shows GraphQL client setup and endpoint configuration. |

**Selected imports (6 files, ~378 lines total):**
1. `src/frontend/src/app/layout.tsx` (27 lines)
2. `src/frontend/next.config.ts` (40 lines)
3. `src/frontend/codegen.ts` (23 lines)
4. `src/frontend/src/lib/urql/client.ts` (24 lines)
5. `src/backend/gqlgen.yml` (105 lines)
6. `src/backend/internal/domain/repository.go` (159 lines)

### Files to Modify

- `CLAUDE.md` -- Add @imports section, enrich "File Locations to Remember", add "Key Patterns" section

### Steps

#### Step 1: Add @imports to CLAUDE.md

Add an `@import` block near the top of CLAUDE.md (after the "Directory Structure" section but before "Key Technical Decisions"). This keeps structural context early where Claude sees it first.

```markdown
## Key Files (auto-imported)

@src/frontend/src/app/layout.tsx
@src/frontend/next.config.ts
@src/frontend/codegen.ts
@src/frontend/src/lib/urql/client.ts
@src/backend/gqlgen.yml
@src/backend/internal/domain/repository.go
```

Each import must use the exact relative path from the workspace root. No leading slash.

#### Step 2: Enrich "File Locations to Remember" section

Replace the current "File Locations to Remember" section (lines 119-130 of CLAUDE.md) with an enriched version using **correct paths** and brief descriptions. Organize by backend and frontend with sub-groupings:

```markdown
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
- Pages: `src/frontend/src/app/` -- App Router pages: `/` (home), `/upload` (document upload flow), `/profile/[id]` (profile view), `/viewer` (PDF viewer)
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
```

#### Step 3: Add "Key Patterns" section

Add a new "Key Patterns" section after "File Locations to Remember". This documents the architectural patterns so Claude does not need to infer them:

```markdown
## Key Patterns

### Data Flow: GraphQL Request

1. Frontend sends query/mutation via urql client to `/api/graphql` (proxied by Next.js to `localhost:8080/graphql`)
2. gqlgen routes to resolver method in `schema.resolvers.go`
3. Resolver calls domain repository interfaces (defined in `domain/repository.go`)
4. Repository implementations in `internal/repository/postgres/` execute SQL via Bun ORM
5. Resolver uses converter functions (`converter.go`) to map domain entities to GraphQL models
6. Response flows back through urql to the React component

### Data Flow: Document Upload & Extraction

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
```

#### Step 4: Verify @import paths resolve correctly

After editing, verify each @import path exists:
```bash
ls -la src/frontend/src/app/layout.tsx
ls -la src/frontend/next.config.ts
ls -la src/frontend/codegen.ts
ls -la src/frontend/src/lib/urql/client.ts
ls -la src/backend/gqlgen.yml
ls -la src/backend/internal/domain/repository.go
```

All paths must resolve. If any fail, fix the @import line.

#### Step 5: Check total CLAUDE.md size

After editing, verify the file remains reasonable:
```bash
wc -l CLAUDE.md
wc -c CLAUDE.md
```

The current file is 178 lines / 5833 bytes. The enriched version will be larger (estimated ~280-320 lines) but this is acceptable because:
- The bean spec says to enrich with higher-signal content
- The predecessor bean (credfolio2-zsm2) already slimmed CLAUDE.md by ~50% by extracting situational content
- The new content is high-value architectural pointers that prevent redundant exploration
- @imports add file contents at load time, not to CLAUDE.md itself

#### Step 6: Run lint and tests

```bash
pnpm lint
pnpm test
```

Both must pass. Since this is a documentation-only change, failures are unlikely.

#### Step 7: Create branch, commit, push, and create PR

```bash
git checkout -b pqsx-add-imports-and-architectural-pointers
git add CLAUDE.md
git commit --no-gpg-sign -m "docs: Add @imports and architectural pointers to CLAUDE.md

Add 6 @imports for frequently-needed small files (layout, config, codegen,
urql client, gqlgen config, repository interfaces). Enrich File Locations
section with accurate paths and descriptions. Add Key Patterns section
documenting data flows, code generation, domain structure, and test
organization.

Co-Authored-By: Claude <noreply@anthropic.com>"
git push -u origin pqsx-add-imports-and-architectural-pointers
```

Create PR for human review.

### Testing Strategy

- **No automated tests to write** -- this is a documentation-only change
- **Verification**: `pnpm lint` and `pnpm test` must pass
- **Path verification**: `ls -la` each @import path to confirm it resolves
- **Size check**: Verify CLAUDE.md stays under ~350 lines after enrichment
- **Manual review**: Read through the enriched sections to confirm accuracy against actual codebase structure

### Open Questions

None -- all path corrections and import decisions have been resolved through codebase exploration.

---

## Definition of Done
- [x] Key files identified and @imports added to CLAUDE.md
- [x] "File Locations to Remember" section enriched with descriptions
- [x] "Key Patterns" section added with architectural overview
- [x] Verified imports resolve correctly (no broken paths)
- [x] Total CLAUDE.md size remains reasonable after enrichment (261 lines)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] Branch pushed and PR created for human review
