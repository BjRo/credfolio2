---
# credfolio2-03qd
title: Profile page layout and routing
status: draft
type: task
priority: normal
created_at: 2026-01-23T16:28:03Z
updated_at: 2026-01-23T16:29:48Z
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

- [ ] Create profile/[id]/page.tsx
- [ ] Set up data fetching via GraphQL
- [ ] Compose all profile components
- [ ] Add action buttons area
- [ ] Handle profile not found
- [ ] Add loading skeleton