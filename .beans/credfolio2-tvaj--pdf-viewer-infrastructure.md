---
# credfolio2-tvaj
title: PDF Viewer Infrastructure
status: completed
type: feature
priority: normal
created_at: 2026-02-07T09:29:03Z
updated_at: 2026-02-07T10:01:32Z
parent: credfolio2-klgo
blocking:
    - credfolio2-5mkd
---

Set up the foundational PDF viewing capability in the frontend.

## Checklist

- [x] Install `react-pdf` (and its peer dep `pdfjs-dist`) in `src/frontend/`
- [x] Configure PDF.js worker for Next.js (copy worker to `public/` or configure bundler path)
- [x] Create `<PDFViewer>` component at `src/frontend/src/components/viewer/PDFViewer.tsx`
  - Props: `fileUrl: string`, `highlightText?: string`, `onHighlightResult?: (found: boolean) => void`
  - Renders PDF pages with text layer enabled
  - Handles loading state (skeleton per page)
  - Handles error state (corrupted/missing PDF)
- [x] Add zoom controls (fit width, fit page, zoom in/out)
- [x] Add page navigation (prev/next buttons, "Page N of M" indicator)
- [x] Verify worker loads correctly in dev mode (`pnpm dev`)
- [x] Verify text layer renders correctly (text is selectable in rendered PDF)

## Technical Notes

- `react-pdf` v9+ supports React 19
- Worker file path is typically: `pdfjs-dist/build/pdf.worker.min.mjs`
- For Next.js, may need to set `pdfjs.GlobalWorkerOptions.workerSrc` in a client component
- Text layer CSS from `react-pdf/dist/Page/TextLayer.css` must be imported
- Consider using `<Document>` and `<Page>` components from react-pdf

## Definition of Done
- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] Visual verification with agent-browser (for UI changes)
- [x] All other checklist items above are completed
- [x] Branch pushed and PR created for human review
- [x] Automated code review passed (`@review-backend` and/or `@review-frontend`)