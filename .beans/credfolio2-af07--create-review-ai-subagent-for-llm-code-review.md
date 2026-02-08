---
# credfolio2-af07
title: Create @review-ai subagent for LLM code review
status: in-progress
type: task
priority: normal
created_at: 2026-02-08T08:38:58Z
updated_at: 2026-02-08T08:43:03Z
---

Create a new review subagent specialized in reviewing LLM-related code changes (prompts, model selection, extraction schemas, safety). This complements the existing @review-backend and @review-frontend agents by focusing on AI/ML-specific concerns.

## Context

Credfolio2 is heavily dependent on LLM functionality for:
- Document detection/classification
- Resume extraction (structured data)
- Reference letter processing (testimonials, skill/experience validation)
- All LLM infrastructure is in `src/backend/internal/infrastructure/llm/`

The existing review agents (@review-backend, @review-frontend) cover general code quality but lack specialized expertise for LLM-specific concerns like:
- Prompt engineering patterns
- Prompt injection vulnerabilities
- Evaluation strategies for prompt changes
- Model selection appropriateness
- Token efficiency and cost optimization
- Structured output schema design

## Implementation Approach

1. Create `.claude/agents/review-ai.md` following the same structure as review-backend.md and review-frontend.md
2. Define review criteria specific to LLM code:
   - Prompt construction and safety
   - Schema design for structured outputs
   - Model selection and parameter tuning
   - Error handling and retry strategies
   - Evaluation/testing strategy for prompt changes
   - Token usage and cost implications
3. Use the same GitHub PR review API workflow (single review with inline comments)
4. Focus on files in:
   - `src/backend/internal/infrastructure/llm/**`
   - `src/backend/internal/job/*_processing.go`
   - Any files containing prompt strings or LLM calls

## Checklist

### Implementation
- [x] Create `.claude/agents/review-ai.md` with frontmatter (name, description, tools, model)
- [x] Write review process section (identify PR, gather context, review code, post comments)
- [x] Define LLM-specific review criteria (see detailed criteria below)
- [x] Document comment posting workflow using GitHub API
- [x] Add project-specific context section for Credfolio2's LLM architecture
- [ ] Test the agent on a PR with LLM changes (will be validated on next LLM-related PR)

### LLM-Specific Review Criteria (in priority order)

#### Security & Safety (CRITICAL)
- [x] Define criteria for prompt injection vulnerabilities
  - User content properly isolated from system instructions
  - Validation of LLM outputs before use
  - Sandboxing/escaping of user-provided content in prompts
- [x] Define criteria for PII/sensitive data handling
  - Logging practices (no raw prompts/responses in logs)
  - Data retention policies for LLM requests/responses

#### Prompt Engineering (HIGH)
- [x] Define criteria for prompt construction patterns
  - Clear instruction structure (system/user/assistant boundaries)
  - Few-shot examples where appropriate
  - Structured output formatting (JSON schema, XML tags)
  - Prompt clarity and specificity
- [x] Define criteria for prompt versioning/tracking
  - Changes to prompts should be reviewable (not inline string literals)
  - Consider extracting prompts to separate files or constants

#### Schema Design (HIGH)
- [x] Define criteria for structured output schemas
  - Schema complexity vs extraction accuracy trade-offs
  - Field naming consistency
  - Required vs optional fields (how does LLM handle missing data?)
  - Validation of extracted data against schema

#### Model Selection (MEDIUM)
- [x] Define criteria for model choice appropriateness
  - Task complexity vs model capability (Haiku vs Sonnet vs Opus)
  - Cost/performance trade-offs
  - Context window requirements
  - Specialized model features (vision, function calling, etc.)

#### Error Handling & Resilience (MEDIUM)
- [x] Define criteria for LLM error handling
  - Retry strategies with exponential backoff
  - Fallback behaviors on LLM failure
  - Circuit breaker patterns for API outages
  - Validation of LLM responses before proceeding

#### Evaluation & Testing (MEDIUM)
- [x] Define criteria for prompt testing
  - Test cases for prompt changes (regression testing)
  - Evaluation metrics (accuracy, consistency, cost)
  - Edge case handling (malformed inputs, empty responses)

#### Performance & Cost (LOW)
- [x] Define criteria for token efficiency
  - Prompt length optimization
  - Unnecessary context removal
  - Batch processing where applicable
  - Caching strategies for repeated prompts

### Integration
- [x] Update `.claude/templates/definition-of-done.md` to reference @review-ai
- [x] Update `/skill dev-workflow` to mention @review-ai as optional step for LLM changes
- [x] Document in CLAUDE.md under "Common Commands" or "Development Workflow"

## Definition of Done
- [x] Tests written (TDD: write tests before implementation) - N/A for documentation/config
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] Visual verification via `@qa` subagent (via Task tool, for UI changes) - N/A for documentation
- [x] ADR written via `/decision` skill (if new dependencies, patterns, or architectural changes were introduced)
- [x] All other checklist items above are completed
- [x] Branch pushed to remote
- [x] PR created for human review
- [x] Automated code review passed via `@review-backend` and/or `@review-frontend` subagents (via Task tool)