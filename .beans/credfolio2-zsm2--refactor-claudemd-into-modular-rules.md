---
# credfolio2-zsm2
title: Refactor CLAUDE.md into modular rules
status: in-progress
type: task
priority: normal
created_at: 2026-02-07T15:59:13Z
updated_at: 2026-02-07T19:42:01Z
parent: credfolio2-ynmd
---

Break up the 327-line root CLAUDE.md into focused, modular files under `.claude/rules/`.

## Why

The entire CLAUDE.md loads into context at session start. Much of its content is situational — database debugging commands aren't needed when working on frontend components, and visual verification instructions aren't needed when writing migrations. Modular rules with path-specific scoping load content only when relevant, keeping the base context lean.

## What

Audit the current CLAUDE.md and extract concerns into `.claude/rules/`:

### Keep in CLAUDE.md (always-loaded essentials)
- Project overview & directory structure
- Key technical decisions (package manager, build system, stacks)
- Common commands (build, dev, clean)
- Git details & workflow rules
- The "STOP — Before Marking Any Work Complete" checklist
- Decision documentation reference

### Move to `.claude/rules/` (situational content)

| File | Content | Path scope |
|------|---------|------------|
| `database.md` | PostgreSQL connection details, psql commands, migration commands, database debugging | `src/backend/**` |
| `dev-servers.md` | Starting dev servers, killing stale processes, port cleanup, common issues | None (unconditional but rarely needed) |
| `visual-verification.md` | Agent-browser workflow, fixture resume, upload instructions | `src/frontend/**` |
| `backend.md` | Backend-specific patterns, Go conventions, file locations | `src/backend/**` |
| `frontend.md` | Frontend-specific patterns, Next.js/React/Tailwind conventions, file locations | `src/frontend/**` |
| `devcontainer.md` | Devcontainer notes, tool management, rebuild info | `.devcontainer/**` |

### Path-specific scoping example
```yaml
---
paths:
  - "src/backend/**"
---

# Database Debugging with psql
...
```

## Approach
- Read the current CLAUDE.md carefully
- Identify each section and whether it's always-needed or situational
- Create the rules files with appropriate `paths` frontmatter
- Slim down CLAUDE.md to the essentials
- Verify nothing is lost by comparing before/after

---

## Implementation Plan

### Approach

Extract four rule files from CLAUDE.md (database, dev-servers, visual-verification, devcontainer). Skip creating separate `backend.md` and `frontend.md` rule files because CLAUDE.md contains almost no backend/frontend-specific content beyond what `golang.md` already covers -- the "File Locations to Remember" section (12 lines) is a cross-cutting reference that belongs in the always-loaded base. If richer backend/frontend content is added later (e.g., via bean credfolio2-pqsx), it can be extracted at that time.

The existing `.claude/rules/golang.md` already handles Go conventions with `src/backend/**` scoping. It will not be duplicated.

### Pre-existing state

- `.claude/rules/golang.md` exists (3 lines of Go guidance, scoped to `src/backend/**` and `**/*.go`)
- No other rule files exist yet
- CLAUDE.md is 319 lines

### Section-by-line mapping of current CLAUDE.md

| Lines | Section | Disposition |
|-------|---------|-------------|
| 1 | Title | KEEP |
| 3-20 | STOP -- Before Marking Any Work Complete | KEEP |
| 22-40 | Directory Structure | KEEP |
| 42-78 | Key Technical Decisions (5 subsections) | KEEP |
| 80-96 | Common Commands (build/dev/clean only) | KEEP (trim DB/migration lines out) |
| 97-115 | Common Commands (DB/migration portion) | EXTRACT to `database.md` |
| 117-161 | Database Debugging with psql | EXTRACT to `database.md` |
| 163-194 | Important Context (Permissions, Build Req, Network, File Locations) | KEEP |
| 196-206 | Development Workflow | KEEP |
| 208-213 | Git Details | KEEP |
| 215-246 | Starting Dev Servers | EXTRACT to `dev-servers.md` |
| 248-282 | Visual Verification with Fixture Resume | EXTRACT to `visual-verification.md` |
| 284-291 | Devcontainer Notes | EXTRACT to `devcontainer.md` |
| 293-319 | Decision Documentation | KEEP |

### Files to Create

1. **`.claude/rules/database.md`** -- PostgreSQL connection details, psql debugging, migration commands
2. **`.claude/rules/dev-servers.md`** -- Dev server startup/shutdown procedure, port cleanup, common issues
3. **`.claude/rules/visual-verification.md`** -- Agent-browser workflow, fixture resume, upload instructions
4. **`.claude/rules/devcontainer.md`** -- Devcontainer tooling, mise, rebuild info

### Files to Modify

1. **`CLAUDE.md`** -- Remove extracted sections, keep everything else

### Files NOT Created (deviation from original bean spec)

- **`backend.md`** -- Not needed. `golang.md` already exists for Go conventions. The backend file locations (4 lines) are part of a cross-cutting "File Locations" section that should stay in CLAUDE.md. No other backend-specific prose exists in CLAUDE.md to extract.
- **`frontend.md`** -- Not needed. Frontend file locations (3 lines) are part of the same cross-cutting section. The network restrictions note (3 lines) is generic advice. No frontend-specific patterns or conventions are documented in CLAUDE.md. The `vercel-react-best-practices` skill already provides frontend guidance when activated.

### Steps

#### Step 1: Create `.claude/rules/database.md`

Create the file at `/workspace/.claude/rules/database.md` with the following content:

**Frontmatter:**
```yaml
---
paths:
  - "src/backend/**"
---
```

**Body -- composed from two CLAUDE.md sections:**

1. The database/migration commands from "Common Commands" (current lines 97-114):
   - Starting with `# Database (run from host, not devcontainer)`
   - Through `# Password: credfolio_dev`
   - Wrap these in a `## Common Database Commands` heading and a code block

2. The entire "Database Debugging with psql" section (current lines 117-161):
   - Keep all three subsections: Connection Details (table), Quick Commands, Environment Variable for Password
   - Keep verbatim

#### Step 2: Create `.claude/rules/dev-servers.md`

Create the file at `/workspace/.claude/rules/dev-servers.md` with the following content:

**Frontmatter -- no path scoping (unconditional, loaded when no path matches):**
```yaml
---
description: "How to start, stop, and troubleshoot dev servers"
---
```

Note: Files without `paths` frontmatter are loaded unconditionally by Claude Code. Since the bean says "None (unconditional but rarely needed)", we use `description` only. This is acceptable -- it is still only ~30 lines and the content is operational knowledge needed across contexts.

**Body -- from CLAUDE.md lines 215-246:**
- The entire "Starting Dev Servers" section verbatim
- Includes the bash code block (pkill/fuser/lsof sequence), "Why this approach" list, and "Common issues" list

#### Step 3: Create `.claude/rules/visual-verification.md`

Create the file at `/workspace/.claude/rules/visual-verification.md` with the following content:

**Frontmatter:**
```yaml
---
paths:
  - "src/frontend/**"
---
```

**Body -- from CLAUDE.md lines 248-282:**
- The entire "Visual Verification with Fixture Resume" section verbatim
- Includes the note about @qa subagent, the agent-browser code block, and the "Key points" list

#### Step 4: Create `.claude/rules/devcontainer.md`

Create the file at `/workspace/.claude/rules/devcontainer.md` with the following content:

**Frontmatter:**
```yaml
---
paths:
  - ".devcontainer/**"
---
```

**Body -- from CLAUDE.md lines 284-291:**
- The entire "Devcontainer Notes" section verbatim (6 bullet points)

#### Step 5: Slim down CLAUDE.md

Edit `/workspace/CLAUDE.md` to remove the extracted sections. The resulting file should contain these sections in order:

1. **Title** (line 1) -- keep
2. **STOP -- Before Marking Any Work Complete** (lines 3-20) -- keep verbatim
3. **Directory Structure** (lines 22-40) -- keep verbatim
4. **Key Technical Decisions** (lines 42-78) -- keep verbatim
5. **Common Commands** (lines 80-96) -- keep but REMOVE lines 97-114 (the database/migration commands). The code block should end after `rm -rf .turbo  # Clear Turborepo cache` with its closing triple-backtick
6. **Important Context** (lines 163-194) -- keep verbatim (Permissions, Build Requirements, Network Restrictions, File Locations)
7. **Development Workflow** (lines 196-206) -- keep verbatim
8. **Git Details** (lines 208-213) -- keep verbatim
9. **Decision Documentation** (lines 293-319) -- keep verbatim

**Sections REMOVED from CLAUDE.md:**
- Lines 97-115: Database/migration commands from Common Commands code block
- Lines 117-161: "Database Debugging with psql" entire section
- Lines 215-246: "Starting Dev Servers" entire section
- Lines 248-282: "Visual Verification with Fixture Resume" entire section
- Lines 284-291: "Devcontainer Notes" entire section

**Estimated slimmed CLAUDE.md size:** ~160 lines (down from 319, a ~50% reduction)

#### Step 6: Verify no content was lost

Run a verification pass:
1. Concatenate the slimmed CLAUDE.md + all rule files (database.md, dev-servers.md, visual-verification.md, devcontainer.md, golang.md)
2. Confirm every heading from the original CLAUDE.md appears in exactly one of these files
3. Spot-check key details: PostgreSQL connection table, pkill sequence, agent-browser commands, mise note

Specific checks:
- `psql -h credfolio2-postgres` appears in `database.md`
- `pkill -f "turbo run dev"` appears in `dev-servers.md`
- `agent-browser open` appears in `visual-verification.md`
- `debian:bookworm-slim` appears in `devcontainer.md`
- `make migrate-up` appears in `database.md`
- `pnpm build` appears in slimmed `CLAUDE.md`
- `@.claude/templates/definition-of-done.md` appears in slimmed `CLAUDE.md`

#### Step 7: Run lint and tests

```bash
pnpm lint
pnpm test
```

Both must pass. Since this change only modifies markdown files, failures are unlikely but should be verified.

#### Step 8: Create branch, commit, and push

```bash
git checkout -b zsm2-refactor-claude-md-modular-rules
git add CLAUDE.md .claude/rules/database.md .claude/rules/dev-servers.md .claude/rules/visual-verification.md .claude/rules/devcontainer.md
git commit --no-gpg-sign -m "refactor: Extract situational CLAUDE.md content into modular .claude/rules/ files

Move database debugging, dev server management, visual verification, and
devcontainer notes into path-scoped rule files. Reduces always-loaded
context by ~50% while preserving all content.

Co-Authored-By: Claude <noreply@anthropic.com>"
git push -u origin zsm2-refactor-claude-md-modular-rules
```

#### Step 9: Create PR

Create a pull request with a summary listing what was extracted and where.

### Testing Strategy

- **No automated tests to write** -- this is a documentation-only change
- **Verification**: Confirm `pnpm lint` and `pnpm test` pass (no code changes)
- **Manual check**: Compare original CLAUDE.md line count (319) vs slimmed version (~160) + rule files total to ensure no content is missing
- **Grep check**: Search for key distinctive strings from each extracted section to confirm they appear in the correct rule file

### Coordination Notes

- **Bean credfolio2-pqsx** ("Add @imports and architectural pointers to CLAUDE.md") should be done AFTER this bean. That bean enriches CLAUDE.md with higher-signal content; this bean slims it down first. The "File Locations to Remember" section stays in CLAUDE.md for now and can be enriched by that follow-up bean.
- **Bean credfolio2-nj3i** ("Create dev-server shutdown script") may eventually replace the manual shutdown instructions in `dev-servers.md`. That is fine -- the rule file documents current state.
- **Bean credfolio2-ubn6** ("Create database reset script") mentions updating CLAUDE.md Common Commands. After this refactor, database commands live in `.claude/rules/database.md` instead, so that bean's plan should be updated accordingly.

### Open Questions

None -- the scope is well-defined and all decisions are resolved.

---

## Definition of Done
- [x] `.claude/rules/` directory created with modular rule files
- [x] Path-specific scoping applied where appropriate
- [x] Root CLAUDE.md reduced to essential always-needed content
- [x] No information lost — all content preserved across files
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [ ] Branch pushed and PR created for human review
