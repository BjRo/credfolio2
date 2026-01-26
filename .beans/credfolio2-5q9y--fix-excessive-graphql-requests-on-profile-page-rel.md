---
# credfolio2-5q9y
title: Fix excessive GraphQL requests on profile page reload
status: completed
type: bug
created_at: 2026-01-26T14:47:54Z
updated_at: 2026-01-26T14:47:54Z
---

The profile page was issuing infinite GraphQL requests on reload due to URQL's suspense mode.

## Root Cause
The URQL client had `suspense: true` configured, which caused an infinite loop of requests. This happens when:
1. Suspense mode throws promises to trigger React Suspense boundaries
2. Without proper Suspense boundary handling (or with React 19 / Next.js 16 behaviors), this causes repeated re-renders
3. Each re-render triggers new GraphQL requests

## Checklist
- [x] Investigate and confirm root cause with code review
- [x] Disable suspense mode in URQL client
- [x] Memoize callback functions with useCallback
- [x] Test that requests are reduced on page reload

## Changes Made
1. **provider.tsx**: Disabled suspense mode (`suspense: false`) - this was the fix
2. **page.tsx**: Wrapped `handleMutationSuccess` in `useCallback` to prevent unnecessary re-renders (minor optimization)