# Add @review-ai Subagent for LLM Code Review

**Date**: 2026-02-08
**Bean**: credfolio2-af07

## Context

Credfolio2 is heavily dependent on LLM functionality for core features:
- Document detection and classification
- Resume extraction (structured data from unstructured text)
- Reference letter processing (testimonials, skill/experience validation)

All LLM infrastructure is centralized in `src/backend/internal/infrastructure/llm/` and `src/backend/internal/job/*_processing.go`. While the existing `@review-backend` and `@review-frontend` subagents provide excellent general code review, they lack specialized expertise in LLM-specific concerns that are critical to the application's correctness, safety, and cost-effectiveness.

Key gaps in current review coverage:
- **Security**: Prompt injection vulnerabilities when user content is incorporated into prompts
- **Prompt Engineering**: Unclear instructions, missing examples, suboptimal output format specifications
- **Schema Design**: Extraction schemas that are too complex (low accuracy) or too simple (missing critical data)
- **Model Selection**: Using expensive models (Opus) when cheaper ones (Haiku/Sonnet) would suffice, or vice versa
- **Evaluation**: No systematic way to test prompt changes for regression
- **Cost**: Inefficient prompts, unnecessary context, missed caching opportunities

These concerns are distinct from general backend code quality and warrant specialized review.

## Decision

Introduced a new specialized code review subagent: `@review-ai` (AI/ML Engineer reviewer).

**Implementation**:
1. Created `.claude/agents/review-ai.md` following the same structure as `@review-backend` and `@review-frontend`
2. Defined LLM-specific review criteria across 6 priority levels:
   - **Security & Safety** (CRITICAL): Prompt injection, PII handling, output validation
   - **Prompt Engineering** (HIGH): Construction patterns, versioning, structured output
   - **Schema Design** (HIGH): Complexity trade-offs, field design, validation
   - **Model Selection** (MEDIUM): Task appropriateness, cost/performance balance
   - **Error Handling & Resilience** (MEDIUM): Retry logic, fallbacks, circuit breakers
   - **Evaluation & Testing** (MEDIUM): Regression testing, metrics, edge cases
   - **Performance & Cost** (LOW): Token efficiency, caching, observability
3. Integrated with existing workflow:
   - Updated `.claude/templates/definition-of-done.md` to reference `@review-ai`
   - Updated `.claude/skills/dev-workflow/SKILL.md` with usage instructions
   - Updated `CLAUDE.md` pre-completion checklist

**Usage Pattern**:
The agent is invoked conditionally when LLM-related files change (`src/backend/internal/infrastructure/llm/**`, `src/backend/internal/job/*_processing.go`). It posts findings as inline PR comments via the GitHub API, just like the other review agents.

## Reasoning

**Why a separate agent vs extending @review-backend?**
- LLM review requires fundamentally different expertise (AI/ML vs Go/backend engineering)
- Keeps agents focused and maintainable (single responsibility principle)
- Allows different review criteria priorities (e.g., prompt injection is CRITICAL for LLM code, irrelevant for general backend)
- Enables selective invocation (only run when LLM code changes, avoiding noise on non-LLM PRs)

**Why prioritize these specific criteria?**
- **Security first**: Prompt injection can lead to data exfiltration, privilege escalation, or corrupted extraction
- **Prompt quality second**: Directly impacts extraction accuracy, which is the core value proposition
- **Cost last**: Important but rarely critical (can be optimized post-launch)

**Alternatives considered**:
1. **Manual LLM review by human expert**: Too slow, doesn't scale, blocks PRs
2. **Expand @review-backend with LLM checklist**: Dilutes focus, makes agent too large, hard to maintain
3. **Post-hoc LLM audits**: Catches issues too late, expensive to fix after merge

## Consequences

**Positive**:
- **Early detection of LLM-specific issues**: Prompt injection, schema problems, model misuse caught in PR review
- **Knowledge sharing**: Review comments educate team on LLM best practices
- **Cost optimization**: Identifies expensive model usage before it reaches production
- **Systematic evaluation**: Encourages regression testing for prompt changes

**Neutral**:
- **Additional review step**: PRs with LLM changes now require running `@review-ai` (adds ~30-60s per PR)
- **Agent maintenance**: New agent file to keep updated as LLM patterns evolve

**Future implications**:
- When adding new LLM features, developers should expect `@review-ai` to scrutinize prompt design
- Prompt changes should include rationale and evaluation strategy (agent will flag if missing)
- Consider extracting prompts to separate files/constants as codebase grows (agent recommends this)
- May want to add LLM-specific testing patterns (golden datasets for prompt validation) as the team matures

**Migration path**:
- Existing LLM code is not automatically reviewed; agent only applies to new PRs going forward
- Consider backlog bean to audit existing LLM code against `@review-ai` criteria
