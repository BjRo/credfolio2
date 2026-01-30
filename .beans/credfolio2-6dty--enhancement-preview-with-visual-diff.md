---
# credfolio2-6dty
title: Validation preview with granular selection
status: in-progress
type: feature
priority: normal
created_at: 2026-01-23T16:28:37Z
updated_at: 2026-01-30T10:18:19Z
parent: credfolio2-1kt0
blocking:
    - credfolio2-1twf
---

Preview screen showing what a reference letter will validate, with granular selection controls.

## Preview Layout

After reference letter extraction completes, show a preview with three sections:

### 1. Corroborations Section
Existing skills and experiences that this reference letter validates.

- Header: "Skills & Experiences Your Reference Validates"
- Each item shows:
  - The skill/experience name
  - Quote snippet from the reference letter
  - Checkbox to include/exclude
- Visual: Green checkmark or validation icon

### 2. Testimonials Section
Full quotes suitable for display on profile.

- Header: "Testimonials to Add"
- Each testimonial shows:
  - Full quote text
  - Author attribution (name, title, relationship)
  - Checkbox to include/exclude
  - "Skills validated" tags below quote

### 3. Discovered Skills Section
Skills the reference mentioned that aren't in the profile yet.

- Header: "Skills Your Reference Noticed"
- Styled differently (suggestion/prompt style, not default selected)
- Each shows:
  - Skill name
  - Quote context
  - Checkbox (unchecked by default)
- Encourage user to add if relevant

## Selection Controls

- "Select All" / "Deselect All" buttons per section
- Individual checkboxes per item
- Running count: "X of Y selected"
- "Apply Selected" primary action button
- "Cancel" secondary button

## Data Flow

1. Receive `referenceLetterID` from upload flow
2. Query extraction results: `referenceLetter(id) { extractedData }`
3. Query current profile to compute matches
4. Display preview with pre-computed matches
5. On "Apply Selected" -> call applyValidations mutation
6. On success -> navigate to profile page

## Checklist

- [x] Create ValidationPreview page/component
- [x] Create CorroborationsSection component
- [x] Create TestimonialsSection component
- [x] Create DiscoveredSkillsSection component
- [x] Implement selection state management
- [x] Create SelectionControls component (select all, count)
- [x] Query extraction results and current profile
- [x] Compute skill/experience matches
- [x] Style discovered skills as suggestions
- [x] Wire up "Apply Selected" to mutation
- [x] Handle loading and error states
- [x] Navigate to profile on success

## Definition of Done

- [ ] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] Visual verification with agent-browser (UI changes)
- [ ] All checklist items above are completed
- [ ] Branch pushed and PR created for human review
