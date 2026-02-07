---
# credfolio2-0gh5
title: Deep Link from Skill Validation Popover
status: draft
type: feature
priority: normal
created_at: 2026-02-07T09:29:33Z
updated_at: 2026-02-07T09:29:33Z
parent: credfolio2-klgo
---

Add a "View in source document" link to the skill validation popover that opens the PDF viewer with the validating quote highlighted.

## Checklist

- [ ] Identify the validation popover component in SkillsSection.tsx
- [ ] Verify the GraphQL query for skill validations includes:
  - `referenceLetterID` (or `referenceLetter.id`)
  - `quoteSnippet` (the text used for validation)
  - If not available, extend the query/schema
- [ ] Add a "View in source" link/button to the validation popover:
  - Icon: `ExternalLink` or `FileText` from lucide-react
  - Text: "View in source" or just an icon with tooltip
  - Opens in new tab: `target="_blank"`
  - URL: `/viewer?letterId={referenceLetterID}&highlight={encodeURIComponent(quoteSnippet)}`
- [ ] Only show the link when a reference letter is available (some validations may not have a source document)
- [ ] Style consistently with existing popover content

## Technical Notes

- The skill validation popover currently shows: author name, author title/company, quote snippet
- The popover data comes from the `GetProfileSkills` query which includes validation details
- Check if `reference_letter_id` is exposed in the GraphQL response for skill validations
- If not, extend `SkillValidation` GraphQL type to include `referenceLetter { id }`

## Definition of Done
- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Visual verification with agent-browser (for UI changes)
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review
- [ ] Automated code review passed (`@review-backend` and/or `@review-frontend`)