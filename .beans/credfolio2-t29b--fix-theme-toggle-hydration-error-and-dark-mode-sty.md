---
# credfolio2-t29b
title: Fix theme toggle hydration error and dark mode styling
status: completed
type: bug
priority: normal
created_at: 2026-01-27T17:19:14Z
updated_at: 2026-01-27T17:20:13Z
---

Two issues:

1. **Hydration error in ThemeToggle**: Server renders Monitor icon but client renders Moon/Sun based on resolved theme. Need to defer icon rendering until after mount.

2. **Dark mode styling issues**:
   - ProfileActions.tsx: hardcoded `bg-white` should be `bg-card`
   - WorkExperienceSection.tsx RoleCard: hardcoded `text-gray-600` should be `text-muted-foreground`

## Checklist
- [x] Fix ThemeToggle hydration mismatch
- [x] Fix hardcoded bg-white in ProfileActions
- [x] Fix hardcoded text-gray-600 in WorkExperienceSection RoleCard
- [x] `pnpm lint` passes
- [x] `pnpm test` passes