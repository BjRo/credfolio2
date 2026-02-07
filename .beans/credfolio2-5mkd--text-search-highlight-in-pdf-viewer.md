---
# credfolio2-5mkd
title: Text Search & Highlight in PDF Viewer
status: in-progress
type: feature
priority: normal
created_at: 2026-02-07T09:29:14Z
updated_at: 2026-02-07T10:21:40Z
parent: credfolio2-klgo
blocking:
    - credfolio2-owyl
---

Implement the ability to search for a text string in the loaded PDF and visually highlight the match.

## Design Decisions

- **Architecture:** Custom hook (`useTextHighlight`) for testability, with pure utility functions in `textSearch.ts`
- **Scope:** Single-page matching only (cross-page highlighting is a follow-up)
- **Highlight color:** Yellow with 40% opacity (`rgba(255, 255, 0, 0.4)`)
- **Approach:** `customTextRenderer` with two-pass rendering (see Technical Notes)

## Architecture: Three-Layer Design

```
textSearch.ts          (pure functions, zero React deps)
    ↓
useTextHighlight.ts    (custom hook, bridges react-pdf to utilities)
    ↓
PDFViewer.tsx          (wires hook outputs to <Page> props)
```

## Checklist

- [x] Create `src/frontend/src/lib/textSearch.test.ts` — write tests for `normalizeText`, `findMatchInPage`, `renderHighlightedText`, `escapeHtml`
- [x] Create `src/frontend/src/lib/textSearch.ts` — implement pure utility functions (types: `HighlightRange`, `TextItemSlice`, `PageMatchResult`)
- [x] Create `src/frontend/src/hooks/useTextHighlight.test.ts` — write hook tests with `renderHook`
- [x] Create `src/frontend/src/hooks/useTextHighlight.ts` — implement custom hook
- [x] Update `src/frontend/src/components/viewer/PDFViewer.test.tsx` — extend react-pdf mock, add highlight integration tests
- [x] Update `src/frontend/src/components/viewer/PDFViewer.tsx` — integrate hook, wire props to `<Page>`
- [x] Update `src/frontend/src/app/globals.css` — add `.pdf-highlight` CSS class

## Technical Notes

### Two-Pass Rendering (Critical Constraint)

In react-pdf v10, `customTextRenderer` runs in a `useLayoutEffect` BEFORE `onGetTextSuccess` fires in a `useEffect`. So match data computed in `onGetTextSuccess` is NOT available during the first `customTextRenderer` call.

**Solution:** Two-pass with `renderKey` state counter:
1. First pass: `customTextRenderer` has no match data → returns plain text
2. `onGetTextSuccess` fires → computes match → stores in ref → increments `renderKey`
3. `renderKey` change gives `customTextRenderer` a new identity → react-pdf re-renders text layer
4. Second pass: `customTextRenderer` reads match data from ref → returns highlighted HTML

The gap is ~one frame, imperceptible.

### textSearch.ts Utilities

**`normalizeText(text)`** — collapse whitespace, trim, replace ligatures (`\uFB00-\uFB04`), smart quotes (`\u2018-\u201D`), dashes (`\u2013-\u2014`), NBSP (`\u00A0`).

**`findMatchInPage(items, searchText)`** algorithm:
1. For each item, `normalizeWithMapping(item.str)` → `{ normalized, mapping }` where `mapping[normalizedIdx] = originalCharIdx`
2. Concatenate normalized strings → `pageText`, build `globalMapping[globalIdx] = { itemIndex, originalCharIdx }`
3. Case-insensitive `indexOf` on `pageText` for normalized `searchText`
4. Map match range back to per-item `HighlightRange`s using `globalMapping`
5. Group by `itemIndex`: `startOffset = min(originalCharIdx)`, `endOffset = max(originalCharIdx) + 1`

**`renderHighlightedText(str, ranges)`** — produce HTML string for `customTextRenderer`. Wraps matched portions in `<mark class="pdf-highlight">`. All text portions HTML-escaped (XSS prevention via `innerHTML`).

### useTextHighlight Hook

**Options:** `{ highlightText?, numPages, onHighlightResult?, pageRefs }`

**Internal state:**
- `matchDataRef: Map<number, PageMatchResult>` (ref — read synchronously by customTextRenderer)
- `renderKey: number` (state — triggers new customTextRenderer identity)
- `firstMatchPageRef`, `hasScrolledRef`, `hasReportedRef`, `searchedPagesRef` (all refs)

**Returns:**
- `getOnGetTextSuccess(pageNumber)` — factory for per-page `onGetTextSuccess` callbacks
- `customTextRenderer` — pass to `<Page>`, memoized with `[highlightText, renderKey]`, `undefined` when no search
- `getOnRenderTextLayerSuccess(pageNumber)` — factory for auto-scroll trigger

**`itemIndex` correctness:** react-pdf passes `itemIndex` as the index in the full `textContent.items` array (including `TextMarkedContent` items). The hook must use the same indexing: iterate `textContent.items`, skip non-TextItems, but use the array index as the key.

**Auto-scroll:** In `getOnRenderTextLayerSuccess`, query for `.pdf-highlight` element in the page DOM, call `scrollIntoView({ behavior: "smooth", block: "center" })`.

**Reset:** `useEffect([highlightText])` clears all refs and resets `renderKey`.

### PDFViewer.tsx Integration

1. Remove `_` prefix from `highlightText` and `onHighlightResult`
2. Call `useTextHighlight` hook
3. Create memoized per-page callbacks via `useMemo([numPages, getOnGetTextSuccess, getOnRenderTextLayerSuccess])`
4. Pass `onGetTextSuccess`, `customTextRenderer`, `onRenderTextLayerSuccess` to each `<Page>`

### CSS

```css
.pdf-highlight {
  background-color: rgba(255, 255, 0, 0.4);
  border-radius: 2px;
}
```

Yellow on white PDF pages (always white even in dark mode per existing convention).

## Definition of Done
- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] Visual verification with agent-browser (N/A — component not yet mounted in any route; that's credfolio2-owyl)
- [x] All other checklist items above are completed
- [x] Branch pushed and PR created for human review
- [ ] Automated code review passed (`@review-backend` and/or `@review-frontend`)