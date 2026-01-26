---
# credfolio2-03qd
title: Profile page layout and routing
status: completed
type: task
priority: normal
created_at: 2026-01-23T16:28:03Z
updated_at: 2026-01-26T10:41:46Z
parent: credfolio2-umxd
blocking:
    - credfolio2-zhnh
    - credfolio2-kndk
    - credfolio2-oqzw
---

Create the profile page that assembles all components.

## Route

- `/profile/[id]` - View a specific profile
- After upload, redirect to `/profile/{new-profile-id}`

## Layout

```
┌─────────────────────────────────────┐
│         ProfileHeader               │
│         ProfileSummary              │
├─────────────────────────────────────┤
│  Experience        │  Skills        │
│  Timeline/Cards    │  Education     │
│                    │  Certifications│
│                    │                │
├─────────────────────────────────────┤
│  Actions: Add Reference | Export    │
└─────────────────────────────────────┘
```

## Technical

- Server component with data fetching
- Client components for interactive parts
- Loading state with skeleton
- Error boundary for fetch failures

## Checklist

- [x] Create profile/[id]/page.tsx
- [x] Set up data fetching via GraphQL
- [x] Compose all profile components
- [x] Add action buttons area
- [x] Handle profile not found
- [x] Add loading skeleton