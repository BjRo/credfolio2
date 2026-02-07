---
# credfolio2-fbrc
title: Audit & restructure large skills to use supporting files
status: completed
type: task
priority: normal
created_at: 2026-02-07T16:01:38Z
updated_at: 2026-02-07T20:35:18Z
parent: credfolio2-ynmd
---

Audit all skills for size and restructure large ones to use supporting files, keeping SKILL.md focused and under 500 lines.

## Why

Large skills consume significant context when loaded. Claude Code's documentation recommends keeping SKILL.md under 500 lines and moving detailed reference material to separate files that load on-demand. This keeps the skill's core instructions lean while preserving all the detailed content for when it's actually needed.

## Skills to audit

| Skill | Concern |
|-------|---------|
| `vercel-react-best-practices` | 45 rules across 8 categories — likely well over 500 lines |
| `agent-browser` | 60+ commands with detailed usage — very large |
| `web-design-guidelines` | Fetches external content, review size |
| `dev-workflow` | Automated workflow with scripts, review size |
| Others | Quick check, likely fine |

> **Note:** `review-backend` and `review-frontend` have been converted from skills to subagents (`.claude/agents/`). They no longer need to be audited as skills. The `agent-browser` skill is preloaded by the `@qa` subagent but still exists as a standalone skill and should be audited for size.

## What

For each oversized skill:
1. Measure current SKILL.md line count
2. Identify what must stay in SKILL.md (overview, when to use, key instructions)
3. Extract detailed content into supporting files:
   - e.g., `vercel-react-best-practices/rules.md` for the full 45 rules
   - e.g., `agent-browser/commands.md` for the command reference
   - e.g., `agent-browser/examples.md` for usage examples
4. Reference supporting files from SKILL.md so Claude knows when to load them:
   ```markdown
   ## Additional resources
   - For the complete rule set, see [rules.md](rules.md)
   - For command reference, see [commands.md](commands.md)
   ```
5. Verify skills still work correctly after restructuring

## Note

`review-backend`, `review-frontend`, and `agent-browser` have already been restructured as part of the subagent migration. `review-backend` and `review-frontend` are now subagent definitions in `.claude/agents/`. The `agent-browser` skill is preloaded by the `@qa` subagent. Focus restructuring efforts on the remaining skills listed above.

## Implementation Plan

### Audit Results

Complete line count audit of all SKILL.md files:

| Skill | SKILL.md Lines | Over 500? | Has Supporting Files? |
|-------|---------------|-----------|----------------------|
| `docker-expert` | 408 | No | No |
| `tdd` | 371 | No | Yes (`testing-anti-patterns.md`) |
| `dev-workflow` | 292 | No | No |
| `agent-browser` | 253 | No | No |
| `vercel-react-best-practices` | 126 | No | Yes (`rules/` dir, `AGENTS.md`, `README.md`) |
| `decision` | 114 | No | No |
| `issue-tracking-with-beans` | 104 | No | No |
| `implement` | 76 | No | No |
| `refine` | 51 | No | No |
| `web-design-guidelines` | 40 | No | No |

### Key Finding: No Skills Exceed 500 Lines

**None of the SKILL.md files exceed the 500-line threshold.** The original bean was created based on the assumption that several skills were "likely well over 500 lines," but the actual measurements show all are under 500 lines:

- `vercel-react-best-practices` (126 lines) -- Already well-structured. The SKILL.md contains only a quick-reference summary of the 45 rules with category tables. The detailed content lives in 53 individual rule files in `rules/` and a compiled `AGENTS.md` (2516 lines). This skill is already a model of the supporting-files pattern.
- `agent-browser` (253 lines) -- Contains a comprehensive command reference as a single SKILL.md. While substantial, it is under 500 lines. It serves as a standalone skill and is also preloaded by the `@qa` subagent.
- `docker-expert` (408 lines) -- The largest SKILL.md at 408 lines. Contains expertise areas, code examples, review checklists, and diagnostics. Under the threshold but the closest to it.
- `tdd` (371 lines) -- Already uses a supporting file (`testing-anti-patterns.md`, 300 lines) referenced from the main SKILL.md. Good existing pattern.
- `dev-workflow` (292 lines) -- Contains the full development workflow. Under threshold.
- `web-design-guidelines` (40 lines) -- Minimal. Fetches guidelines from an external URL at runtime.
- All others (`decision`, `issue-tracking-with-beans`, `implement`, `refine`) are well under 200 lines.

### Approach

Since no skills exceed 500 lines, this bean requires **audit only, not restructuring**. The work is:

1. Document the audit results (this plan captures them)
2. Verify skills that already use supporting files (`vercel-react-best-practices`, `tdd`) have proper references
3. No files need to be created, modified, or restructured

### Steps

1. **Record the audit** -- The audit is complete. All 10 skills have been measured. Results are documented above.

2. **Verify existing supporting file patterns** -- Two skills already use supporting files correctly:
   - `vercel-react-best-practices/SKILL.md` (line 110-116, 124-126): References `rules/*.md` files and `AGENTS.md` for detailed content. Pattern is correct.
   - `tdd/SKILL.md` (line 359): References `testing-anti-patterns.md` via `@testing-anti-patterns.md`. Pattern is correct.

3. **No restructuring needed** -- All SKILL.md files are under 500 lines. The bean's checklist items for "skills over 500 lines restructured" and "supporting files contain detailed content" are vacuously satisfied (no skills are over 500 lines, so no restructuring is needed).

4. **Run lint and tests** -- Verify `pnpm lint` and `pnpm test` pass (no changes to source code, so these should pass trivially).

5. **Create PR** -- Branch, commit the audit documentation update to this bean, push, and create PR.

### Testing Strategy

- No source code changes, so no new tests needed
- `pnpm lint` and `pnpm test` should pass without changes
- Verify skills are still invocable after any bean metadata updates (manual spot check)

### Open Questions

None -- the audit is straightforward and the results are clear. No skills need restructuring.

## Definition of Done
- [x] All skills audited for line count
- [x] Skills over 500 lines restructured with supporting files (N/A — none exceeded 500 lines)
- [x] SKILL.md files contain overview + references to supporting files (N/A — no restructuring needed)
- [x] Supporting files contain detailed content (rules, commands, examples) (N/A — no restructuring needed)
- [x] All restructured skills tested to verify they still work (N/A — no restructuring needed)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] Branch pushed and PR created for human review
