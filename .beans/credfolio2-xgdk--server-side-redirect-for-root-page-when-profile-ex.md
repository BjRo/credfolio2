---
# credfolio2-xgdk
title: Server-side redirect for root page when profile exists
status: todo
type: task
priority: normal
created_at: 2026-01-27T16:38:45Z
updated_at: 2026-02-05T14:37:44Z
parent: credfolio2-abtx
---

The current root page (/) uses a client-side redirect via `useEffect` + `router.push()` when a completed resume exists. This causes a brief flash of the loading spinner before redirecting.

Refactor to use a Next.js server-side redirect instead, so the user is redirected before any client-side rendering occurs. This would involve:

- Moving the resume existence check to a Server Component or `redirect()` call in the page
- Fetching the user's resumes server-side (SSR or RSC) instead of via client-side `useQuery`
- Using Next.js `redirect()` from `next/navigation` in a server context

This improves perceived performance and eliminates the loading flash for returning users.

## Checklist
- [ ] Move resume existence check to server-side (RSC or middleware)
- [ ] Use Next.js `redirect()` for server-side redirect
- [ ] Keep client-side upload form rendering for new users
- [ ] Update tests to cover the new server-side behavior

## Definition of Done
- [ ] Tests written (TDD: write tests before implementation)
- [ ] `pnpm lint` passes with no errors
- [ ] `pnpm test` passes with no failures
- [ ] Visual verification with agent-browser (for UI changes)
- [ ] All other checklist items above are completed