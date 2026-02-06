---
# credfolio2-pnu4
title: Fix React Strict Mode double-firing mutations in ExtractionProgress
status: todo
type: bug
priority: normal
created_at: 2026-02-05T23:08:28Z
updated_at: 2026-02-06T07:29:27Z
parent: credfolio2-3ram
---

## Summary

In development, React 19 Strict Mode double-invokes effects (mount → unmount → remount). The `ExtractionProgress` component fires the `processDocument` mutation inside a `useEffect`, causing it to execute **twice** — creating duplicate resume/reference letter records and enqueuing duplicate background jobs.

Observed in logs: two "Letter extraction completed" entries with different `reference_letter_id` values for the same `file_id`.

## Root Cause

`ExtractionProgress.tsx` lines 184-191:

```tsx
useEffect(() => {
    mountedRef.current = true;
    startProcessing();  // fires processDocument mutation
    return () => { mountedRef.current = false; };
}, [startProcessing]);
```

The `mountedRef` guard only prevents state updates after unmount — it does NOT prevent the second mount from calling `startProcessing()` again. Both fetches fire, both mutations succeed, both create DB records.

**Also affects**: `DetectionProgress.tsx` likely has the same pattern (polling `setInterval` in effect). Audit both components.

## Fix

Add an `isStartedRef` to prevent re-invocation:

```tsx
const isStartedRef = useRef(false);

useEffect(() => {
    if (isStartedRef.current) return;
    isStartedRef.current = true;
    startProcessing();
    ...
```

## Scope

- Dev-only issue (Strict Mode double-fire doesn't happen in production builds)
- But causes confusing duplicate logs and wasted LLM API calls during development

## Checklist

- [ ] Add `isStartedRef` guard to `ExtractionProgress.tsx` `startProcessing` effect
- [ ] Audit `DetectionProgress.tsx` for same pattern
- [ ] Add test verifying mutation is only called once

### Definition of Done
- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review