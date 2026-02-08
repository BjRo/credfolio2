---
# credfolio2-8do3
title: Improve skill validation matching algorithm
status: todo
type: task
priority: normal
created_at: 2026-02-08T11:14:51Z
updated_at: 2026-02-08T11:14:51Z
parent: credfolio2-nihn
---

Replace aggressive substring matching in skill validation with more sophisticated algorithm to reduce false positives.

## Problem (from @review-backend)

Current implementation uses simple substring matching which causes false positives:
- "Java" matches "JavaScript" 
- "C" matches "CSS", "React", "Docker", etc.
- "Go" matches "Django", "PostgreSQL", etc.

This leads to incorrect skill validations and poor user experience.

## Current Implementation

```go
// src/backend/internal/service/materialization.go
strings.Contains(strings.ToLower(context), strings.ToLower(skillName))
```

## Impact

- **User Experience**: Users see skills validated that weren't actually mentioned
- **Data Quality**: Profile skill validation counts are inflated and misleading
- **Trust**: Users may lose confidence in AI extraction accuracy

## Proposed Solutions

1. **Word Boundary Matching** - Use regex with word boundaries (`\bJava\b`)
2. **Fuzzy Matching** - Allow slight variations but require close match (Levenshtein distance)
3. **NLP-based Matching** - Use embeddings or entity recognition
4. **Hybrid Approach** - Exact match for short skills (≤3 chars), word boundary for longer

## Files Affected

- `src/backend/internal/service/materialization.go` (validateSkillMention function)
- Tests for skill validation logic

## Acceptance Criteria

- [ ] "Java" does not match "JavaScript"
- [ ] "C" does not match "CSS" or "React"  
- [ ] "Go" does not match "Django"
- [ ] Legitimate variations still match ("React.js" matches "React", "Node" matches "Node.js")
- [ ] Tests cover common false positive cases
- [ ] Performance remains acceptable (no N+1 queries added)

## Reference

See: /documentation/reviews/2026-02-08-comprehensive-codebase-review.md#warnings-4

## Definition of Done
- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Visual verification via `@qa` subagent (via Task tool, for UI changes)
- [ ] ADR written via `/decision` skill (if new dependencies, patterns, or architectural changes were introduced)
- [ ] All other checklist items above are completed
- [ ] Branch pushed to remote
- [ ] PR created for human review
- [ ] Automated code review passed via `@review-backend`, `@review-frontend`, and/or `@review-ai` (for LLM changes) subagents (via Task tool)