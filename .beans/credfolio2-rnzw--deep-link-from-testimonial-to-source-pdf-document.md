---
# credfolio2-rnzw
title: Deep link from testimonial to source PDF document
status: scrapped
type: feature
priority: normal
created_at: 2026-02-02T14:20:08Z
updated_at: 2026-02-02T16:57:28Z
parent: credfolio2-2ex3
---

## Overview

Allow users to navigate directly from a testimonial quote to the specific location in the source PDF reference letter where the quote appears.

## User Story

As a profile viewer, I want to click on a testimonial and be taken to the exact page/location in the PDF where that quote appears, so I can verify the context and read the full reference letter.

## Investigation Areas

- [ ] Research PDF deep linking capabilities (page anchors, text fragments)
- [ ] Evaluate PDF viewer options that support highlighting/navigation
- [ ] Determine if we store page numbers during testimonial extraction
- [ ] Consider browser-native PDF viewing vs embedded viewer (pdf.js)
- [ ] Explore text fragment URLs (Chrome's `#:~:text=` syntax)

## Technical Considerations

- PDF.js for in-app viewing with programmatic navigation
- Store extraction metadata (page number, bounding box) during processing
- Fallback: open PDF at correct page if exact text highlighting not possible
- Mobile considerations for PDF viewing experience

## Out of Scope (for now)

- Highlighting the exact text within the PDF (stretch goal)
- Editing or annotating the PDF

## Definition of Done
- [ ] Research completed and approach documented
- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Visual verification with agent-browser (for UI changes)
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review