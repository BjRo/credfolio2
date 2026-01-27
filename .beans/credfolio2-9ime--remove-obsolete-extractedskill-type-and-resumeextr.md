---
# credfolio2-9ime
title: Remove obsolete ExtractedSkill type and ResumeExtractedData.skills
status: completed
type: task
priority: normal
created_at: 2026-01-27T18:41:32Z
updated_at: 2026-01-27T18:45:37Z
---

## Context

The skills unification refactoring (dee3216) materialized extracted skills into ProfileSkill records.
This makes `ExtractedSkill` (reference letter extraction type) and `ResumeExtractedData.skills` obsolete
from the GraphQL API perspective, matching the pattern already applied to education/experience.

## Checklist

- [x] Remove `ExtractedSkill` GraphQL type from schema
- [x] Change `ExtractedLetterData.skills` from `[ExtractedSkill!]!` to `[String!]!`
- [x] Remove `skills` field from `ResumeExtractedData` in GraphQL schema
- [x] Remove `ExtractedSkill` model mapping from gqlgen.yml
- [x] Remove `model.ExtractedSkill` Go struct from extraction.go
- [x] Update `model.ExtractedLetterData.Skills` to `[]string`
- [x] Remove `model.ResumeExtractedData.Skills` field
- [x] Remove `toGraphQLExtractedSkills` converter
- [x] Update `toGraphQLExtractedData` to convert skills to string names
- [x] Update `toGraphQLResumeExtractedData` to drop skills
- [x] Remove `extractedSkillResolver` and `ExtractedSkill()` resolver method
- [x] Regenerate backend generated code (go generate)
- [x] Update frontend queries (GetReferenceLetter, GetResume)
- [x] Regenerate frontend GraphQL types

## Definition of Done
- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] All other checklist items above are completed