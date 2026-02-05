---
# credfolio2-fap7
title: Explore splitting resume extraction into focused LLM steps
status: draft
type: feature
priority: high
created_at: 2026-02-05T17:40:18Z
updated_at: 2026-02-05T17:56:45Z
parent: credfolio2-dwid
---

## Context

Observation from testing: gpt-5-mini performs well on skill extraction but poorly on work experience extraction and summarization. This suggests the current monolithic single-call approach (one model does everything) isn't optimal — different extraction tasks have different complexity requirements.

## Current State

Resume extraction is a single LLM call that asks one model to extract everything at once: skills, work experiences, education, personal info, and summarize accomplishments.

## Problem

- Skill extraction is structured pattern matching — smaller/faster models handle it well
- Experience summarization requires reasoning about context and distilling accomplishments — needs a more capable model
- A single model choice forces a compromise: either overpay for simple tasks or underperform on complex ones

## Possible Approaches

### 1. Multi-step pipeline (sequential)
Run skill extraction with a fast/cheap model (e.g., gpt-5-mini), then run experience/education extraction with a more capable model (e.g., gpt-5-nano or larger). Each step gets a focused prompt.

**Pros:** Simple to implement, focused prompts improve quality  
**Cons:** Higher total latency (sequential), more LLM calls

### 2. Parallel fan-out extraction
Fan out multiple focused LLM calls simultaneously (skills, experiences, education, personal info), potentially with different models per task.

**Pros:** Faster wall-clock time, model-per-task flexibility  
**Cons:** More complex orchestration, need to merge results

### 3. Tiered approach
Start with a fast model for everything, then selectively re-process sections that need deeper reasoning with a more capable model.

**Pros:** Fast for easy resumes, quality boost only where needed  
**Cons:** Most complex logic, harder to reason about

## Relationship to Other Work

- Related to credfolio2-ujt0 (separate River queue for LLM extraction) — parallel fan-out would benefit from independent worker pools
- Related to credfolio2-yiqg (reference extraction model) — same theme of right-sizing models to tasks

## Checklist

- [ ] Profile current extraction to understand which parts are slow vs fast
- [ ] Experiment with splitting the extraction prompt into focused sub-prompts
- [ ] Benchmark quality of different models on each sub-task independently
- [ ] Choose an approach based on findings
- [ ] Implement the chosen approach
- [ ] Update tests