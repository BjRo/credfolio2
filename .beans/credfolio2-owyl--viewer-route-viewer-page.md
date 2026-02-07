---
# credfolio2-owyl
title: Viewer Route (/viewer page)
status: completed
type: feature
priority: normal
created_at: 2026-02-07T09:29:28Z
updated_at: 2026-02-07T11:07:35Z
parent: credfolio2-klgo
blocking:
    - credfolio2-0gh5
    - credfolio2-kahb
---

Create the /viewer Next.js route that ties together the PDF viewer, text highlighting, and data fetching.

## Checklist

- [x] Create page at `src/frontend/src/app/viewer/page.tsx` (client component)
- [x] Parse query params: `letterId` (required UUID), `highlight` (optional string)
- [x] GraphQL query to fetch reference letter by ID:
  - Need `file.url` (presigned URL for the PDF)
  - Need `title` or `authorName` for the toolbar display
  - Created new lightweight `GetReferenceLetterForViewer` query
- [x] Page layout:
  - Top toolbar: back/close button, document title/author
  - PDFViewer component (has its own toolbar with page nav + zoom controls)
  - Info banner area (conditionally shown when highlight not found)
- [x] Handle states:
  - Loading: show spinner while fetching letter data
  - Error (letter not found / unauthorized): show "Document not found" with back button
  - Error (PDF load failure): PDFViewer's built-in error display
  - Success (highlight found): PDF shown, quote highlighted
  - Success (highlight not found): PDF shown + info banner "Could not locate exact quote — showing full document"
- [x] Info banner:
  - Subtle, non-intrusive (amber/yellow background)
  - Dismissible (X button)
  - Text: "Could not locate exact quote — showing full document"
- [x] Ensure page works when opened in a new tab (no shared client state needed)
- [ ] Test with fixture PDF from `fixtures/CV_TEMPLATE_0004.pdf` or a reference letter

## Technical Notes

- This is a `"use client"` page since it uses PDF.js (browser-only APIs)
- PDFViewer loaded via `next/dynamic` with `ssr: false` to avoid DOMMatrix SSR error
- Page wrapped in `<Suspense>` boundary for `useSearchParams()` compatibility
- Query params are accessed via `useSearchParams()` from `next/navigation`
- The viewer should work standalone — all data comes from the URL params + GraphQL
- `buildViewerUrl()` utility created at `src/frontend/src/lib/viewer.ts` for downstream beans
- Presigned URLs expire after 1 hour. PDFViewer's built-in error display handles load failures.

## Definition of Done
- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [ ] Visual verification with agent-browser (for UI changes)
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review
- [ ] Automated code review passed (`@review-backend` and/or `@review-frontend`)
