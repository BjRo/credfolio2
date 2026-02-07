---
# credfolio2-tvaj
title: PDF Viewer Infrastructure
status: draft
type: feature
priority: normal
created_at: 2026-02-07T09:29:03Z
updated_at: 2026-02-07T09:29:45Z
parent: credfolio2-klgo
blocking:
    - credfolio2-5mkd
---

Set up the foundational PDF viewing capability in the frontend.

## Checklist

- [ ] Install `react-pdf` (and its peer dep `pdfjs-dist`) in `src/frontend/`
- [ ] Configure PDF.js worker for Next.js (copy worker to `public/` or configure bundler path)
- [ ] Create `<PDFViewer>` component at `src/frontend/src/components/viewer/PDFViewer.tsx`
  - Props: `fileUrl: string`, `highlightText?: string`, `onHighlightResult?: (found: boolean) => void`
  - Renders PDF pages with text layer enabled
  - Handles loading state (skeleton per page)
  - Handles error state (corrupted/missing PDF)
- [ ] Add zoom controls (fit width, fit page, zoom in/out)
- [ ] Add page navigation (prev/next buttons, "Page N of M" indicator)
- [ ] Verify worker loads correctly in dev mode (`pnpm dev`)
- [ ] Verify text layer renders correctly (text is selectable in rendered PDF)

## Technical Notes

- `react-pdf` v9+ supports React 19
- Worker file path is typically: `pdfjs-dist/build/pdf.worker.min.mjs`
- For Next.js, may need to set `pdfjs.GlobalWorkerOptions.workerSrc` in a client component
- Text layer CSS from `react-pdf/dist/Page/TextLayer.css` must be imported
- Consider using `<Document>` and `<Page>` components from react-pdf

## Definition of Done
- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Visual verification with agent-browser (for UI changes)
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review
- [ ] Automated code review passed (`@review-backend` and/or `@review-frontend`)