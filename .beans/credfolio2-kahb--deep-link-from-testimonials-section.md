---
# credfolio2-kahb
title: Deep Link from Testimonials Section
status: in-progress
type: feature
priority: normal
created_at: 2026-02-07T09:29:40Z
updated_at: 2026-02-07T11:52:38Z
parent: credfolio2-klgo
---

Modify the existing "View source document" link in the testimonials section to use the new PDF viewer with quote highlighting instead of opening the raw PDF.

## Key Findings from Codebase Exploration

### 1. `buildViewerUrl` utility already exists and is fully tested
- **Location:** `src/frontend/src/lib/viewer.ts` (lines 1-18)
- **Tests:** `src/frontend/src/lib/viewer.test.ts` (lines 1-42)
- Handles: `letterId` param, `highlight` text param, URL encoding via `URLSearchParams`, truncation of highlight text > 500 chars
- JSDoc says "Used by TestimonialsSection and ValidationPopover" — designed for this purpose but **not yet imported by any component**

### 2. Current "View source document" link
- **File:** `src/frontend/src/components/profile/TestimonialsSection.tsx`
- **Line 328:** `getSourceUrl` function returns raw presigned file URL:
  ```typescript
  const getSourceUrl = (testimonial: Testimonial) => testimonial.referenceLetter?.file?.url;
  ```
- **Lines 237-243:** `<a>` tag in `QuoteItem` renders the link with `target="_blank" rel="noopener noreferrer"` (must be preserved)

### 3. GraphQL query already includes `referenceLetter.id`
- **File:** `src/frontend/src/graphql/queries.graphql` lines 288-294
- No GraphQL changes needed

### 4. Viewer page expects these URL params
- **File:** `src/frontend/src/app/viewer/page.tsx` lines 77-78
- Reads `letterId` and `highlight` from `searchParams` — matches `buildViewerUrl` output exactly

### 5. `testimonial.quote` is readily accessible
- Already rendered on line 269 of TestimonialsSection.tsx

### 6. Test fixtures have all needed data
- `mockTestimonialsWithSourceBadge` has `referenceLetter.id: "ref-1"` (line 116) and `quote` (line 99)

---

## Implementation Plan

### Step 1: Update Tests First (TDD)

**File: `src/frontend/src/components/profile/TestimonialsSection.test.tsx`**

**a) Update test "'View source document' links to the PDF file URL" (lines 311-324):**
- Change assertion from expecting `https://example.com/reference-letter.pdf` to expecting `/viewer?letterId=ref-1&highlight=Great+team+player+with+excellent+leadership+skills.`
- Consider renaming test to "'View source document' links to the viewer with quote highlight"

**b) Verify "'View source document' opens in a new tab" test (lines 326-337):**
- Checks `target="_blank"` and `rel="noopener noreferrer"` — should pass without modification

**c) Add new test: "truncates long quotes in viewer URL highlight param":**
- Create testimonial with 600-char quote
- Verify the highlight param is truncated to 500 chars
- Verify letterId is still correct

### Step 2: Modify `getSourceUrl` Function

**File: `src/frontend/src/components/profile/TestimonialsSection.tsx`**

**a) Add import (around line 28-30):**
```typescript
import { buildViewerUrl } from "@/lib/viewer";
```

**b) Replace `getSourceUrl` (line 328):**

Current:
```typescript
const getSourceUrl = (testimonial: Testimonial) => testimonial.referenceLetter?.file?.url;
```

New:
```typescript
const getSourceUrl = (testimonial: Testimonial) => {
  const letterId = testimonial.referenceLetter?.id;
  const hasFile = !!testimonial.referenceLetter?.file?.url;
  if (!letterId || !hasFile) return undefined;
  return buildViewerUrl(letterId, testimonial.quote);
};
```

**No other code changes needed.** The `<a>` tag in `QuoteItem` works with any URL string.

### Step 3: Run Tests & Lint
```bash
pnpm test
pnpm lint
```

### Step 4: Visual Verification
Use `agent-browser` to verify the link opens the viewer page with highlighted quote.

---

## Edge Cases

| Case | Behavior |
|------|----------|
| No file URL but referenceLetter.id exists | `!hasFile` guard returns undefined → menu not shown |
| Empty quote | `buildViewerUrl` omits `highlight` param → viewer shows full doc |
| Very long quotes (>500 chars) | `buildViewerUrl` truncates to 500 chars → viewer substring search still matches |
| Special characters in quotes | `URLSearchParams.set()` handles URL encoding correctly |

## Summary of Files to Change

| File | Change |
|------|--------|
| `src/frontend/src/components/profile/TestimonialsSection.tsx` | Add `buildViewerUrl` import + modify `getSourceUrl` (line 328) |
| `src/frontend/src/components/profile/TestimonialsSection.test.tsx` | Update href assertion (lines 320-323), add long quote truncation test |

## Checklist

- [ ] In TestimonialsSection.tsx, locate the "View source document" dropdown menu item
- [ ] Change the `href` from the raw file URL to the viewer URL:
  - Old: `{testimonial.referenceLetter.file.url}` (raw presigned URL)
  - New: `/viewer?letterId={testimonial.referenceLetter.id}&highlight={encodeURIComponent(testimonial.quote)}`
- [ ] Ensure `referenceLetter.id` is available in the GraphQL query (it should be — verify) ✅ Confirmed available
- [ ] Handle long quotes: if `testimonial.quote` exceeds ~500 chars, truncate for the URL param (the viewer's substring search will still match) — handled by existing `buildViewerUrl`
- [ ] Keep `target="_blank" rel="noopener noreferrer"` behavior
- [ ] Update existing tests in TestimonialsSection.test.tsx to reflect the new URL format
- [ ] Add test for long quote truncation
- [ ] Verify the link works end-to-end with the viewer page

## Definition of Done
- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Visual verification with agent-browser (for UI changes)
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review
- [ ] Automated code review passed (`@review-backend` and/or `@review-frontend`)