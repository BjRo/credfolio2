---
# credfolio2-600v
title: Add item-level select/deselect in extraction review step
status: completed
type: task
priority: normal
created_at: 2026-02-05T23:09:37Z
updated_at: 2026-02-06T20:11:34Z
parent: credfolio2-3ram
---

## Summary

The ExtractionReview step currently shows extracted data (experiences, education, skills, testimonials) in a read-only list with a single "Import to profile" button that imports everything. Users cannot select or deselect individual items before importing.

Users should be able to review extracted data and choose which items to import — e.g., uncheck an incorrectly extracted work experience or deselect skills they don't want on their profile.

The UX should be consistent with the dedicated reference letter preview page (`/profile/[id]/reference-letters/[referenceLetterID]/preview/`), which already implements a polished selection UX with checkbox cards, SelectionControls, and per-section bulk actions.

## Checklist

### Step 1: Expose resume experiences/education/skills in GraphQL schema
- [x] Add `ExtractedWorkExperience` and `ExtractedEducation` types to `schema.graphqls`
- [x] Add `experiences`, `educations`, `skills` fields to `ResumeExtractedData`
- [x] Run `go generate ./...` and add resolvers in `converter.go`

### Step 2: Add selection fields to ImportDocumentResultsInput
- [x] Add `selectedExperienceIndices`, `selectedEducationIndices`, `selectedSkills`, `selectedTestimonialIndices`, `selectedDiscoveredSkills` (all nullable)
- [x] Run `go generate ./...`

### Step 3: Implement backend filtering in resolver
- [x] Add `FilterByIndices` (deduplicating), `FilterSkillsByName`, `FilterDiscoveredSkillsByName` helpers
- [x] Apply filters in resolver before materialization
- [x] Write comprehensive backend tests

### Step 4: Update frontend types and polling query
- [x] Add TypeScript interfaces for extracted work/education
- [x] Update polling query to fetch new fields

### Step 5: Move SelectionControls to shared location
- [x] Create shared component at `@/components/ui/selection-controls.tsx`
- [x] Update imports in reference letter preview components

### Step 6: Rewrite ExtractionReview with selection UX
- [x] Full rewrite with checkbox cards, per-section SelectionControls, footer counter
- [x] All items pre-selected by default (discovered skills NOT pre-selected)
- [x] Keyboard accessible, role="checkbox", consistent styling with preview page

### Step 7: Update tests
- [x] 32 frontend tests covering selection, toggling, mutation payloads, disabled states, counters

### Definition of Done
- [x] Tests written
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures (360 tests)
- [x] Visual verification with agent-browser
- [x] All other checklist items above are completed
- [x] Branch pushed and PR created for human review (PR #95)
- [x] Automated code review passed — addressed critical finding (discovered skills not transmitted)