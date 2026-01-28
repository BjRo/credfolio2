---
# credfolio2-7a0g
title: Investigate broken LLM extraction for WorkExperience, Education, Skills
status: completed
type: bug
created_at: 2026-01-27T19:17:36Z
updated_at: 2026-01-28T06:07:00Z
---

Resume upload extracts data from LLM but results don't appear in profile. Suspect data is silently dropped due to model validation changes during skills management UI work.

## Investigation Findings

### Root Cause: Silent Materialization Failure

The issue is in [resume_processing.go:170-176](src/backend/internal/job/resume_processing.go#L170-L176) - **materialization errors are logged but NOT propagated as job failures**:

```go
if matErr := w.materializeExtractedData(ctx, args.ResumeID, resume.UserID, extractedData); matErr != nil {
    w.log.Error("Failed to materialize extracted data into profile", ...)
    // Log but don't fail — extraction data is still saved in JSONB
}
```

This means ANY error during materialization causes ALL profile data (experiences, education, skills) to be missing from the profile tables, while the resume status shows as `COMPLETED`.

### Why Materialization Fails: Unique Constraint Violations

The `profile_skills` table has a unique constraint on `(profile_id, normalized_name)`:

```sql
-- From migration 20260127120001
CREATE UNIQUE INDEX IF NOT EXISTS idx_profile_skills_unique_name
    ON profile_skills(profile_id, normalized_name);
```

The `materializeSkills()` function at [resume_processing.go:387-413](src/backend/internal/job/resume_processing.go#L387-L413) will fail when:

1. **LLM returns duplicate skills** - e.g., `["Python", "PYTHON", "python"]` all normalize to "python"
2. **Skills already exist from another source** - manual skills or skills from a previously uploaded resume conflict with new extraction

When the first duplicate skill INSERT fails with a unique constraint violation, the entire `materializeExtractedData()` call returns an error, and NO data (including experiences and education) is saved to profile tables.

### Additional Issues Found

1. **No skill deduplication before INSERT** - [resume_processing.go:393-411](src/backend/internal/job/resume_processing.go#L393-L411) iterates through all skills without deduplicating by normalized name first

2. **Unit tests don't catch this** - Mock `mockProfileSkillRepository.Create()` at [resume_processing_test.go:193-199](src/backend/internal/job/resume_processing_test.go#L193-L199) doesn't enforce unique constraint

3. **Delete only removes same-source skills** - `DeleteBySourceResumeID()` at [resume_processing.go:297-299](src/backend/internal/job/resume_processing.go#L297-L299) only deletes skills from the SAME resume, not skills that would conflict from other sources

### Data Flow Summary

```
Resume Upload → LLM Extraction → Save to JSONB ✓ → Materialize to Profile Tables ✗
                                                         ↓
                                              Unique constraint violation
                                                         ↓
                                              Error logged, job continues
                                                         ↓
                                              Resume status = COMPLETED
                                              Profile tables = EMPTY
```

## Fix Applied

### Phase 1: Handle Duplicates Gracefully

#### 1.1 Add skill deduplication in `materializeSkills()`
- [x] Create helper function `deduplicateSkills(skills []string) []string` that:
  - Iterates through skills, normalizing each to lowercase
  - Tracks seen normalized names in a map
  - Returns slice with only first occurrence of each normalized name
- [x] Call deduplication at start of `materializeSkills()` before the loop
- [x] Add unit test `TestDeduplicateSkills` covering:
  - Mixed case duplicates: `["Python", "PYTHON", "python"]` → `["Python"]`
  - No duplicates: `["Go", "Rust"]` → `["Go", "Rust"]`
  - Empty/whitespace handling

#### 1.2 Use `ON CONFLICT DO NOTHING` for skill INSERTs
- [x] Add new method `CreateIgnoreDuplicate(ctx, skill)` to `ProfileSkillRepository` interface
- [x] Implement in postgres repository using `ON CONFLICT (profile_id, normalized_name) DO NOTHING`
- [x] Update `materializeSkills()` to use the new method instead of `Create()`
- [x] Add unit test for repository method (mock should return nil on duplicate)

#### 1.3 Add tests for duplicate skill scenarios
- [x] `TestMaterializeSkillsWithDuplicatesInExtraction` - same skill different cases
- [x] `TestMaterializeSkillsWithExistingManualSkill` - skill already exists from manual entry
- [x] Update mock `mockProfileSkillRepository` to track normalized names and return error on duplicate (simulating real DB behavior)

### Phase 2: Improve Resilience

#### 2.1 Make materialization more granular
- [x] Change `materializeExtractedData()` to collect errors instead of returning on first error
- [x] Process experiences, education, and skills independently (don't abort all if one fails)
- [x] Log individual failures with context (which item failed, why)
- [x] Return aggregated error only if any category failed

#### 2.2 Add tests for partial materialization success
- [x] `TestMaterializePartialSuccess_SkillsFail` - experiences/education succeed, skills fail
- [x] `TestMaterializePartialSuccess_ExperiencesFail` - education/skills succeed, experiences fail
- [x] Verify that successful categories are persisted even when others fail

### Follow-up Work

Optional observability improvements tracked in credfolio2-sfna.

## Definition of Done
- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] Visual verification with agent-browser (for UI changes) - N/A (backend-only changes)
- [x] All checklist items completed
