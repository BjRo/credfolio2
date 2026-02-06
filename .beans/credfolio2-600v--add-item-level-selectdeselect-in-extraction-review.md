---
# credfolio2-600v
title: Add item-level select/deselect in extraction review step
status: in-progress
type: task
priority: normal
created_at: 2026-02-05T23:09:37Z
updated_at: 2026-02-06T14:23:37Z
parent: credfolio2-3ram
---

## Summary

The ExtractionReview step currently shows extracted data (experiences, education, skills, testimonials) in a read-only list with a single "Import to profile" button that imports everything. Users cannot select or deselect individual items before importing.

Users should be able to review extracted data and choose which items to import — e.g., uncheck an incorrectly extracted work experience or deselect skills they don't want on their profile.

The UX should be consistent with the dedicated reference letter preview page (`/profile/[id]/reference-letters/[referenceLetterID]/preview/`), which already implements a polished selection UX with checkbox cards, SelectionControls, and per-section bulk actions.

## Current Behavior

- `ExtractionReview.tsx` renders `CareerInfoSection` and `TestimonialSection` as read-only displays
- "Import to profile" sends `resumeId` and `referenceLetterID` to `importDocumentResults` mutation
- Backend materializes ALL extracted data — no filtering
- GraphQL `ResumeExtractedData` type only exposes header fields (name, email, etc.) — NOT experiences, education, or skills

## Desired Behavior

- Each extracted item (experience, education entry, skill, testimonial) has a checkbox
- All items start selected by default (discovered skills from reference letters NOT pre-selected, matching preview page convention)
- Users can deselect items they don't want imported
- Only selected items get materialized into profile tables
- Visual UX matches the reference letter preview page (checkbox cards, selection controls, footer counter)

## Checklist

### Step 1: Expose resume experiences/education/skills in GraphQL schema

- [x] Add `ExtractedWorkExperience` type to `schema.graphqls` (company, title, location, startDate, endDate, isCurrent, description)
- [x] Add `ExtractedEducation` type to `schema.graphqls` (institution, degree, field, startDate, endDate, gpa, achievements)
- [x] Add `experiences: [ExtractedWorkExperience!]!`, `educations: [ExtractedEducation!]!`, `skills: [String!]!` fields to existing `ResumeExtractedData` type
- [x] Run `go generate ./...` to regenerate models
- [x] Add resolvers for new fields in `converter.go` (mapping from `domain.ResumeExtractedData`)

**Files:** `src/backend/internal/graphql/schema/schema.graphqls`, `src/backend/internal/graphql/resolver/converter.go`

### Step 2: Add selection fields to ImportDocumentResultsInput

- [x] Add `selectedExperienceIndices: [Int!]` to `ImportDocumentResultsInput` (null = import all)
- [x] Add `selectedEducationIndices: [Int!]` to `ImportDocumentResultsInput` (null = import all)
- [x] Add `selectedSkills: [String!]` to `ImportDocumentResultsInput` (null = import all)
- [x] Add `selectedTestimonialIndices: [Int!]` to `ImportDocumentResultsInput` (null = import all)
- [x] Run `go generate ./...` to regenerate models

**Files:** `src/backend/internal/graphql/schema/schema.graphqls`

### Step 3: Implement backend filtering in resolver

- [x] Add generic `FilterByIndices[T any](items []T, indices []int) []T` helper function
- [x] Add `FilterSkillsByName(skills []string, selected []string) []string` helper function
- [x] In `ImportDocumentResults` resolver, filter extracted data arrays before passing to MaterializationService when selection fields are non-nil
- [x] Filter testimonials in reference letter data when `selectedTestimonialIndices` is non-nil
- [x] Write backend tests for `FilterByIndices` (valid/invalid/empty/out-of-range indices, structs, duplicates)
- [x] Write backend tests for `FilterSkillsByName` (exact, case-insensitive, empty, no matches, all, whitespace)

**Files:** `src/backend/internal/graphql/resolver/schema.resolvers.go`, `src/backend/internal/service/materialization.go`, `src/backend/internal/service/materialization_test.go`

### Step 4: Update frontend types and polling query

- [x] Add `ExtractedWorkExperience` interface to `types.ts` (company, title, location, startDate, endDate, isCurrent, description)
- [x] Add `ExtractedEducation` interface to `types.ts` (institution, degree, field, startDate, endDate)
- [x] Extend `ResumeExtractionData.extractedData` to include `experiences: ExtractedWorkExperience[]`, `educations: ExtractedEducation[]`, `skills: string[]`
- [x] Update `DOCUMENT_PROCESSING_STATUS_QUERY` in `ExtractionProgress.tsx` to fetch `experiences { ... }`, `educations { ... }`, `skills`

**Files:** `src/frontend/src/components/upload/types.ts`, `src/frontend/src/components/upload/ExtractionProgress.tsx`

### Step 5: Move SelectionControls to shared location

- [x] Create `SelectionControls` at `src/frontend/src/components/ui/selection-controls.tsx`
- [x] Update import in reference letter preview page (CorroborationsSection, TestimonialsSection, DiscoveredSkillsSection)

**Files:** `src/frontend/src/components/ui/selection-controls.tsx`, reference letter preview components

### Step 6: Rewrite ExtractionReview with selection UX

- [x] Add selection state per category: `Set<number>` for experiences/education/testimonials, `Set<string>` for skills
- [x] Initialize all items as selected by default (discovered skills NOT pre-selected)
- [x] Add checkbox card UI for each experience (bg-success/5 border-success/30 when selected)
- [x] Add checkbox card UI for each education entry (bg-success/5 border-success/30 when selected)
- [x] Add checkbox card UI for each skill (bg-success/5 border-success/30 when selected)
- [x] Add checkbox card UI for each testimonial (bg-primary/5 border-primary/30 when selected, matching preview page)
- [x] Add checkbox card UI for discovered skills (bg-warning/5 border-warning/50 border-2 border-dashed)
- [x] Add SelectionControls (select all / deselect all) per section
- [x] Add footer with "X item(s) selected" counter + "Import Selected" button + "Back" button
- [x] All cards: `role="checkbox"`, keyboard accessible (Space/Enter to toggle), `Checkbox` component
- [x] Update import mutation call to pass selection indices/names

**File:** `src/frontend/src/components/upload/ExtractionReview.tsx`

### Step 7: Update tests

- [x] Update `ExtractionReview.test.tsx` fixtures to include experiences, education, skills arrays
- [x] Test all items render with checkboxes
- [x] Test items start selected by default
- [x] Test toggle behavior (click to deselect/reselect)
- [x] Test mutation receives correct selection indices/names
- [x] Test import disabled when nothing selected
- [x] Test "X items selected" counter updates correctly

**File:** `src/frontend/src/components/upload/ExtractionReview.test.tsx`

### Definition of Done
- [x] Tests written
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [ ] Visual verification with agent-browser (for UI changes)
- [x] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review
- [ ] Automated code review passed (`@review-backend` and/or `@review-frontend`)