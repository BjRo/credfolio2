---
# credfolio2-ccly
title: Skills extracted from reference letters not linked to source document
status: in-progress
type: bug
priority: normal
created_at: 2026-02-02T14:18:26Z
updated_at: 2026-02-02T15:07:36Z
parent: credfolio2-2ex3
---

## Problem

When skills are extracted from reference letters, they are not properly:
1. Marked as credible/validated
2. Linked back to the **testimonials** that mention them

This makes it unclear which skills have third-party validation vs self-reported skills.

## Root Cause Analysis

Three related issues:

### Issue 1: Prompt Structure

The extraction prompt (`reference_letter_extraction_system.txt`) extracts skills in two different formats:

```
SkillMentions:     [{skill, quote, context}]  ← Rich data with attribution
DiscoveredSkills:  ["skill1", "skill2"]       ← Just strings, no attribution
```

Problems:
- The LLM doesn't know what skills exist on the profile, so it can't meaningfully distinguish "mentions of existing skills" vs "discovered new skills"
- `DiscoveredSkills` lacks quotes/context, making it impossible to link them back to the source text

### Issue 2: Data Model Gap

Current state:
- `Testimonial.SkillsMentioned` is `[]string` - just skill names, not links to `ProfileSkill` records
- `SkillValidation` links `ProfileSkillID` → `ReferenceLetterID`, but should link to `TestimonialID` for granular attribution

The link should be: **ProfileSkill ↔ Testimonial** (not just ReferenceLetter)

This enables: "Go is validated by testimonial from John Doe: 'Their Go expertise was exceptional...'"

### Issue 3: Processing Gap

Even with rich extraction data, `reference_letter_processing.go` only stores the JSON blob. It doesn't:
1. Create `SkillValidation` records linking profile skills to testimonials
2. Create new `ProfileSkill` records for discovered skills with proper source tracking

## Solution

### Part A: Prompt Enhancement

1. **Pass existing profile skills in the user prompt** so the LLM can distinguish:
   - `SkillMentions` = skills that match existing profile skills (for validation)
   - `DiscoveredSkills` = new skills not currently on the profile

2. **Align the structures** - both should include:
   - `skill`: The skill name
   - `quote`: Supporting quote from the letter
   - `context`: Category (technical, leadership, etc.)

### Part B: Data Model Enhancement

Update `SkillValidation` to link to testimonials:
- Add `TestimonialID` field (the specific testimonial that validates this skill)
- Keep `ReferenceLetterID` for backward compat / denormalization
- The `QuoteSnippet` becomes the testimonial quote itself

### Part C: Processing Enhancement

After extraction, `reference_letter_processing.go` should:
1. Create `Testimonial` records from extracted testimonials
2. Match skills mentioned in each testimonial against `ProfileSkill` records
3. Create `SkillValidation` records linking `ProfileSkill` → `Testimonial`
4. For discovered skills: create new `ProfileSkill` records with source tracking

## Files to Modify

- `src/backend/internal/infrastructure/llm/prompts/reference_letter_extraction_system.txt` - update prompt instructions
- `src/backend/internal/infrastructure/llm/prompts/reference_letter_extraction_user.txt` - add profile skills context
- `src/backend/internal/domain/extraction.go` - align `DiscoveredSkills` structure
- `src/backend/internal/domain/entities.go` - add `TestimonialID` to `SkillValidation`
- `src/backend/internal/domain/profile.go` - add `SourceReferenceLetterID` to `ProfileSkill`
- `src/backend/internal/job/reference_letter_processing.go` - create testimonials and validation links
- Migration for schema changes

## Checklist

### Part A: Prompt Changes
- [x] Update `ExtractedSkillMention` struct (or create shared struct) for consistent skill format
- [x] Update `DiscoveredSkills` from `[]string` to `[]DiscoveredSkill` (with quote/context)
- [x] Update system prompt to clarify distinction with profile context available
- [x] Update user prompt template to include existing profile skills
- [x] Update `ExtractLetterData` to accept profile skills as context
- [x] Update extraction tests

### Part B: Data Model Changes
- [x] Add `TestimonialID` field to `SkillValidation` entity
- [x] Add `SourceReferenceLetterID` field to `ProfileSkill` entity
- [x] Create migration for new columns
- [x] Update GraphQL schema to expose skill → testimonial links

### Part C: Processing Changes
- [x] Update `reference_letter_processing.go` to:
  - [x] Fetch existing profile skills before extraction
  - [x] Pass skills to extraction call
  - [x] Create `Testimonial` records from extracted data
  - [x] Match `Testimonial.SkillsMentioned` against `ProfileSkill` records
  - [x] Create `SkillValidation` records linking skills to testimonials
  - [x] Create `ProfileSkill` records for discovered skills
- [x] Update processing tests

## Definition of Done
- [x] All checklist items above completed
- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] Visual verification with agent-browser (for UI changes) - N/A: backend-only changes
- [ ] Branch pushed and PR created for human review