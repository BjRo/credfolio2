---
# credfolio2-5mkd
title: Text Search & Highlight in PDF Viewer
status: draft
type: feature
priority: normal
created_at: 2026-02-07T09:29:14Z
updated_at: 2026-02-07T09:29:45Z
parent: credfolio2-klgo
blocking:
    - credfolio2-owyl
---

Implement the ability to search for a text string in the loaded PDF and visually highlight the match.

## Checklist

- [ ] After PDF loads, extract text content from each page via PDF.js text layer API
- [ ] Implement text search logic:
  - Normalize whitespace (collapse multiple spaces, trim)
  - Normalize Unicode (ligatures, smart quotes → standard equivalents)
  - Try exact match first
  - Fall back to substring match (the highlight text may be a fragment of a larger sentence)
- [ ] Highlight rendering:
  - Identify the DOM spans in the text layer that contain the matched text
  - Apply a highlight overlay (CSS class with background color, e.g., yellow with some opacity)
  - Ensure highlight is visible across page boundaries if quote spans pages
- [ ] Auto-scroll to the first highlighted match (smooth scroll, centered in viewport)
- [ ] Expose callback: `onHighlightResult(found: boolean)` so the parent can show a fallback banner

## Technical Notes

- PDF.js text layer renders each text run as a `<span>` in the `.react-pdf__Page__textContent` container
- Matching may need to span multiple spans (a sentence can be split across spans)
- Consider using a mark/highlight approach: wrap matched text in `<mark>` elements or use CSS `::highlight()` if browser support allows
- Normalize curly quotes ('' "" → '' "") and em-dashes (— → --) before matching
- URL-encoded highlight text should be decoded before searching
- If highlight text is very long (>500 chars), truncate to first ~200 chars for search to avoid URL length issues

## Definition of Done
- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Visual verification with agent-browser (for UI changes)
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review
- [ ] Automated code review passed (`@review-backend` and/or `@review-frontend`)