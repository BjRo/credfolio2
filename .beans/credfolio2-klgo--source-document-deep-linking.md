---
# credfolio2-klgo
title: Source Document Deep Linking
status: in-progress
type: epic
priority: normal
created_at: 2026-02-07T09:28:48Z
updated_at: 2026-02-07T09:43:47Z
parent: credfolio2-dwid
---

Enable deep linking from validated skills and testimonials to the respective quote in the source document. Clicking a source link opens a custom PDF viewer in a new tab with the validation quote highlighted.

## Context & Decisions

- **Entry points**: Both skill validation popovers AND testimonials section
- **Viewer**: Custom PDF.js viewer page at `/viewer` route
- **Quote location**: Text search at view time using PDF.js (no extraction pipeline changes)
- **Fallback UX**: Show PDF + info banner when quote text can't be found
- **Shareability**: URL-based (`/viewer?letterId=X&highlight=Y`), bookmarkable, auth-checked

## Architecture Overview

### New Route: `/viewer`
A standalone page that renders a PDF with quote highlighting.

**URL format**: `/viewer?letterId={uuid}&highlight={urlEncodedQuoteText}`

**Flow**:
1. Parse `letterId` and `highlight` from query params
2. Fetch reference letter file URL via GraphQL
3. Load PDF using PDF.js / react-pdf
4. Extract text layer from each page
5. Search for highlight text across pages (with fuzzy tolerance)
6. Scroll to the matching page and visually highlight the text
7. If no match: show subtle info banner, display full document

### Entry Point: Skill Validation Popover
The existing hover popover on skill badges shows who validated the skill and a quote snippet. Add a clickable "View in source" link that opens `/viewer?letterId={id}&highlight={quoteSnippet}` in a new tab.

**Data needed**: `reference_letter_id` and `quote_snippet` from `skill_validations` — already available via GraphQL.

### Entry Point: Testimonials Section  
The existing "View source document" kebab menu item currently opens the raw PDF. Modify it to use the new viewer URL, passing the testimonial quote as the highlight parameter.

## Implementation Plan

### Phase 1: PDF Viewer Infrastructure
1. Install `react-pdf` (wraps pdfjs-dist) as a frontend dependency
2. Configure PDF.js worker for Next.js (worker file needs special bundling setup)
3. Create a reusable `<PDFViewer>` component with:
   - Page-by-page rendering with virtualization (for performance on large PDFs)
   - Text layer enabled (required for search/selection)
   - Zoom controls (fit width, fit page, manual zoom)
   - Page navigation (prev/next, page number input)
   - Loading skeleton per page

### Phase 2: Text Search & Highlighting
1. After PDF loads, build a text index from each page's text content
2. Implement search logic:
   - Exact match first (normalized whitespace)
   - If no exact match, try substring matching (quote may be a fragment)
   - Optionally: fuzzy match with configurable tolerance for minor LLM paraphrasing
3. Highlight rendering:
   - Use PDF.js text layer spans to identify matching text ranges
   - Apply highlight CSS (background color overlay) to matched spans
   - Scroll the matched element into view with smooth scrolling
4. Fallback banner:
   - If no match found after all strategies, show an info banner at the top
   - Banner text: "Could not locate exact quote — showing full document"
   - Banner is dismissible

### Phase 3: Viewer Route (`/viewer`)
1. Create Next.js page at `src/frontend/src/app/viewer/page.tsx`
2. Layout:
   - Top toolbar: back button, document title, page N of M, zoom controls
   - Main area: PDF viewer component
   - Info banner area (conditionally shown)
3. GraphQL query: fetch reference letter by ID, including `file.url`
   - May need a new query or extend existing one
   - Must handle: letter not found, file missing, URL expired
4. Error states:
   - Invalid/missing letterId → "Document not found" page
   - File URL fetch failure → retry with error message
   - PDF load failure → error state with "Try again" option
5. Auth: ensure only authorized users can view the document (existing auth should cover this via presigned URL expiry)

### Phase 4: Wire Up Skill Validation Popover
1. Identify the validation popover component (in SkillsSection.tsx)
2. Ensure GraphQL query for skill validations includes `referenceLetterID` and `quoteSnippet`
3. Add a "View in source document" link/icon button to the popover
4. Link target: `/viewer?letterId={referenceLetterID}&highlight={encodeURIComponent(quoteSnippet)}`
5. Open in new tab (`target="_blank"`)

### Phase 5: Wire Up Testimonials Section
1. Modify the "View source document" menu item in TestimonialsSection.tsx
2. Change href from raw file URL to `/viewer?letterId={referenceLetterID}&highlight={encodeURIComponent(testimonialQuote)}`
3. Ensure the testimonial quote is available in the component context (it should be — it's the displayed quote)

## Technical Considerations

- **PDF.js worker**: Must be configured to load from the correct path in Next.js. Common approach: copy worker to `public/` or use CDN. For devcontainer (no external network), use local worker.
- **Large PDFs**: Use page-level virtualization (only render visible pages + buffer). Reference letters are typically 1-3 pages, so this may be premature optimization.
- **Text matching accuracy**: PDF text extraction can have quirks (ligatures, special characters, hyphenation). Normalize whitespace and punctuation before matching.
- **URL length**: Very long quotes could exceed URL length limits. Consider truncating the highlight param to first ~200 chars if needed, relying on substring search.
- **Presigned URL expiry**: The viewer fetches a presigned URL on load. If the user leaves the tab open for >1 hour, the URL expires. Consider a refresh mechanism or show a "reload" button.

## Out of Scope
- Modifying the LLM extraction pipeline (no page number extraction)
- Annotation/commenting on PDFs
- Editing PDFs
- Supporting non-PDF documents (DOCX, TXT) in the viewer — these could be a follow-up