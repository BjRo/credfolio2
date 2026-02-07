---
# credfolio2-owyl
title: Viewer Route (/viewer page)
status: draft
type: feature
priority: normal
created_at: 2026-02-07T09:29:28Z
updated_at: 2026-02-07T09:29:45Z
parent: credfolio2-klgo
blocking:
    - credfolio2-0gh5
    - credfolio2-kahb
---

Create the /viewer Next.js route that ties together the PDF viewer, text highlighting, and data fetching.

## Checklist

- [ ] Create page at `src/frontend/src/app/viewer/page.tsx` (client component)
- [ ] Parse query params: `letterId` (required UUID), `highlight` (optional string)
- [ ] GraphQL query to fetch reference letter by ID:
  - Need `file.url` (presigned URL for the PDF)
  - Need `title` or `authorName` for the toolbar display
  - Check if existing queries suffice or if a new one is needed
- [ ] Page layout:
  - Top toolbar: back/close button, document title/author, page indicator, zoom controls
  - Main area: `<PDFViewer>` component
  - Info banner area (conditionally shown when highlight not found)
- [ ] Handle states:
  - Loading: show skeleton/spinner while fetching letter data
  - Error (letter not found / unauthorized): show "Document not found" with link back
  - Error (PDF load failure): show error with retry button
  - Success (highlight found): PDF shown, quote highlighted
  - Success (highlight not found): PDF shown + info banner "Could not locate exact quote — showing full document"
- [ ] Info banner:
  - Subtle, non-intrusive (e.g., top of page, light yellow/blue background)
  - Dismissible (X button)
  - Text: "Could not locate exact quote — showing full document"
- [ ] Ensure page works when opened in a new tab (no shared client state needed)
- [ ] Test with fixture PDF from `fixtures/CV_TEMPLATE_0004.pdf` or a reference letter

## Technical Notes

- This is a `"use client"` page since it uses PDF.js (browser-only APIs)
- Query params are accessed via `useSearchParams()` from `next/navigation`
- The viewer should work standalone — all data comes from the URL params + GraphQL
- Presigned URLs expire after 1 hour. If the user returns to a stale tab, the PDF will fail to load. Show a "Document link expired — reload" message with a reload button.

## Definition of Done
- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Visual verification with agent-browser (for UI changes)
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review
- [ ] Automated code review passed (`@review-backend` and/or `@review-frontend`)