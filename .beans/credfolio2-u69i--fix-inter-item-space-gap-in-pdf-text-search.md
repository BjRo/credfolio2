---
# credfolio2-u69i
title: Fix inter-item space gap in PDF text search
status: in-progress
type: bug
priority: high
created_at: 2026-02-07T12:09:06Z
updated_at: 2026-02-07T12:18:35Z
---

## Problem

When the PDF viewer tries to highlight a quote, it shows "Could not locate exact quote — showing full document" for all testimonials. The text search fails because `findMatchInPage` in `textSearch.ts` concatenates PDF text items without inserting spaces between item boundaries.

## Root Cause

In `src/frontend/src/lib/textSearch.ts` lines 119-129, items are concatenated directly:

```typescript
for (let itemIdx = 0; itemIdx < items.length; itemIdx++) {
    const { normalized, mapping } = normalizeWithMapping(items[itemIdx].str);
    for (let j = 0; j < normalized.length; j++) {
      pageTextParts.push(normalized[j]);
      globalMapping.push({ itemIndex: itemIdx, originalCharIdx: mapping[j] });
    }
}
```

react-pdf fragments PDF text into items that may not have explicit spaces at boundaries. Example:
- PDF items: `["Great", "team player"]` (no trailing space on item 0)
- Concatenated: `"Greatteam player"` (missing space!)
- Search for `"Great team player"` → NOT FOUND

This affects ALL quote highlights because real PDFs commonly fragment text this way.

## Proposed Fix

In `findMatchInPage`, inject a synthetic space between adjacent items when the previous item doesn't end with whitespace and the next item doesn't start with whitespace. The space should not be mapped to any real character (use a sentinel mapping entry so the highlight ranges remain correct).

Specifically, after processing each item in the concatenation loop, check:
- If the previous normalized text ended with a non-space character
- And the current normalized text starts with a non-space character
- Then inject a space with a sentinel mapping entry (e.g., `{ itemIndex: -1, originalCharIdx: -1 }`)

The highlight range grouping already filters by itemIndex, so sentinel entries won't affect rendering.

## Actual Fix

Refactored `findMatchInPage` to use a two-pass strategy via `searchWithStrategy`:
1. First try WITHOUT synthetic spaces (handles mid-word splits like `["hel", "lo wor", "ld"]`)
2. If no match, retry WITH synthetic spaces injected between items (handles word-boundary splits like `["Great", "team player"]`)

Sentinel mapping entries (`itemIndex: -1`) are skipped during range grouping so highlight rendering is unaffected.

## Checklist

- [x] Add failing test: items without boundary spaces (`["Great", "team player"]` matches `"Great team player"`)
- [x] Add failing test: items where one has trailing space (still works)
- [x] Add test: items where next starts with whitespace (no double-space)
- [x] Add test: multiple items without boundary spaces (3+ items)
- [x] Fix `findMatchInPage` to use two-pass strategy with synthetic spaces
- [x] Verify existing tests still pass (no regression) — all 36 textSearch tests pass
- [x] Test with real PDF via agent-browser (no reference letters in fixture; unit tests cover the fix)

## Definition of Done
- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures (464 total)
- [x] Visual verification with agent-browser (page renders, no JS errors; no reference letter fixtures available for end-to-end highlight test)
- [x] All other checklist items above are completed
- [x] Branch pushed and PR created for human review (PR #105, same branch as kahb)
- [x] Automated code review passed (`@review-frontend` — LGTM, no critical findings)