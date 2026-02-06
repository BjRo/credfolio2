---
# credfolio2-fjk0
title: 'Fix discovered skills: persistence bug, LLM categorization, and grouped review UI'
status: completed
type: feature
priority: high
created_at: 2026-02-06T20:51:38Z
updated_at: 2026-02-06T22:06:43Z
---

## Problem

Discovered skills from reference letters have multiple issues across both user flows (unified review + separate reference letter extraction):

1. **Unified review flow: discovered skills are silently dropped (bug)**
   - In `ImportDocumentResults` (schema.resolvers.go:980-982), selected discovered skills are filtered but `MaterializeReferenceLetterData` (materialization.go:248-305) only creates testimonials — it never creates `ProfileSkill` records for the selected discovered skills.
   - The `CrossReferenceValidations` step also can't help because those skills don't exist on the profile yet.
   - **Result:** Users select discovered skills, click Import, and nothing happens.

2. **Separate reference letter flow: wrong source + hardcoded category**
   - In `ApplyReferenceLetterValidations` (~schema.resolvers.go:2520), discovered skills are created with `Source: ExperienceSourceManual` instead of `ExperienceSourceLetterDiscovered`.
   - Category is hardcoded to `SkillCategory.Soft` in the frontend (preview/page.tsx:298) regardless of actual skill type.

3. **LLM extraction lacks category assignment**
   - The `DiscoveredSkill` struct only has `skill`, `quote`, and `context` (free-text).
   - No structured category (TECHNICAL/SOFT/DOMAIN) is assigned, so the UI can't group or display skills meaningfully.
   - The prompt loosely says "Focus on soft skills..." but doesn't enforce categorization.

4. **UI shows flat list with no categorization**
   - Both `DiscoveredSkillsSection` components (ExtractionReview.tsx and preview/DiscoveredSkillsSection.tsx) show a flat list of skill name + quote.
   - The `context` field is available but never displayed.
   - Users have no way to understand what kind of skill they're adding.

5. **No category editing before import**
   - Users cannot override the skill category before importing.

## Solution

### UX Decisions (confirmed with user)
- **Skill types:** Extract all types (TECHNICAL, SOFT, DOMAIN) but with proper LLM-assigned categories
- **Review UI:** Group discovered skills by category with sub-headers
- **Category editing:** Users can change the LLM-assigned category via dropdown before importing
- **Attribution on profile:** No special treatment — rely on existing validation/credibility dots
- **Skill-level quotes:** Each discovered skill's quote should be persisted as a SkillValidation for the ValidationPopover

### Implementation scope: Both flows fixed together

## Checklist

### Backend — Domain & Extraction
- [x] Add `Category` field (string: TECHNICAL/SOFT/DOMAIN) to `DiscoveredSkill` struct in `internal/domain/extraction.go`
- [x] Update LLM system prompt (`reference_letter_extraction_system.txt`) to require `category` field on each discovered skill with clear examples
- [x] Update LLM structured output schema / JSON parsing in `internal/infrastructure/llm/extraction.go` to map the new `category` field
- [x] Add `category` field to `DiscoveredSkill` GraphQL type in `schema.graphqls`
- [x] Regenerate GraphQL code (`go generate`)

### Backend — Persistence Fix (Unified Flow)
- [x] In `MaterializeReferenceLetterData` (materialization.go) OR in the `ImportDocumentResults` resolver, add logic to create `ProfileSkill` records for selected discovered skills
  - Source: `ExperienceSourceLetterDiscovered`
  - Category: from LLM-assigned (or user-overridden) category
  - SourceReferenceLetterID: set to the reference letter ID
- [x] Create `SkillValidation` record for each discovered skill, linking to the reference letter with `quoteSnippet` from the extracted quote
- [x] Update `ImportDocumentResultsInput.selectedDiscoveredSkills` from `[String!]` to a structured input (e.g., `[SelectedDiscoveredSkillInput!]`) carrying both `name` and `category`

### Backend — Source Fix (Separate Reference Letter Flow)
- [x] In `ApplyReferenceLetterValidations` resolver (~schema.resolvers.go:2520), change `ExperienceSourceManual` to `ExperienceSourceLetterDiscovered`
- [x] Use the LLM-provided category (from `extractedData`) instead of relying on frontend hardcoding
- [x] Ensure `SourceReferenceLetterID` is set on created ProfileSkill records

### Frontend — Unified Review UI
- [x] Refactor `DiscoveredSkillsSection` in `ExtractionReview.tsx` to group skills by category with sub-headers ("Soft Skills", "Domain Knowledge", "Technical")
- [x] Add category dropdown per skill (editable before import)
- [x] Update selection state from `Set<string>` to a map of skill name → { selected: boolean, category: SkillCategory } to track overrides
- [x] Update `handleImport` to send the new structured input with category per skill

### Frontend — Reference Letter Preview UI
- [x] Update `DiscoveredSkillsSection.tsx` (preview page) with same grouped layout and category editing
- [x] Remove hardcoded `SkillCategory.Soft` in `preview/page.tsx:294-301` — use LLM-provided or user-overridden category
- [x] Update the `ApplyValidationsInput.newSkills` mapping to use per-skill category

### Frontend — Types & GraphQL
- [x] Update GraphQL queries/fragments to include new `category` field on `DiscoveredSkill`
- [x] Regenerate frontend GraphQL types (`pnpm graphql-codegen` or equivalent)
- [x] Update `types.ts` for the extraction review component

### Testing
- [x] Backend: Test that `MaterializeReferenceLetterData` (or resolver) creates ProfileSkill + SkillValidation for discovered skills
- [x] Backend: Test category field flows through extraction → GraphQL → materialization
- [x] Backend: Test `ApplyReferenceLetterValidations` uses `ExperienceSourceLetterDiscovered`
- [x] Frontend: Test grouped rendering and category override behavior

## Definition of Done
- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [ ] Visual verification with agent-browser (for UI changes)
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review
- [ ] Automated code review passed (`@review-backend` and/or `@review-frontend`)

## Key Files

**Backend:**
- `src/backend/internal/domain/extraction.go` — DiscoveredSkill struct
- `src/backend/internal/infrastructure/llm/prompts/reference_letter_extraction_system.txt` — LLM prompt
- `src/backend/internal/infrastructure/llm/prompts/reference_letter_extraction_user.txt` — User prompt template
- `src/backend/internal/infrastructure/llm/extraction.go` — JSON parsing (~line 627)
- `src/backend/internal/service/materialization.go` — MaterializeReferenceLetterData + CrossReferenceValidations
- `src/backend/internal/graphql/schema/schema.graphqls` — GraphQL schema
- `src/backend/internal/graphql/resolver/schema.resolvers.go` — ImportDocumentResults (~line 826), ApplyReferenceLetterValidations (~line 2219)

**Frontend:**
- `src/frontend/src/components/upload/ExtractionReview.tsx` — Unified review DiscoveredSkillsSection (~line 653)
- `src/frontend/src/components/upload/types.ts` — Extraction types
- `src/frontend/src/app/profile/[id]/reference-letters/[referenceLetterID]/preview/DiscoveredSkillsSection.tsx` — Separate flow component
- `src/frontend/src/app/profile/[id]/reference-letters/[referenceLetterID]/preview/page.tsx` — Separate flow page (hardcoded SOFT category at line 298)