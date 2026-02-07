---
# credfolio2-pqsx
title: Add @imports and architectural pointers to CLAUDE.md
status: todo
type: task
created_at: 2026-02-07T16:17:24Z
updated_at: 2026-02-07T16:17:24Z
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

## Definition of Done
- [ ] Key files identified and @imports added to CLAUDE.md
- [ ] "File Locations to Remember" section enriched with descriptions
- [ ] "Key Patterns" section added with architectural overview
- [ ] Verified imports resolve correctly (no broken paths)
- [ ] Total CLAUDE.md size remains reasonable after enrichment
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Branch pushed and PR created for human review