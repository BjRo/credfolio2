---
# credfolio2-63w0
title: Skeleton loading UI during extraction
status: draft
type: feature
created_at: 2026-01-23T16:27:29Z
updated_at: 2026-01-23T16:27:29Z
parent: credfolio2-6oza
---

Show an animated skeleton that mimics the profile layout while processing.

## Design Goals

- Skeleton should look like the actual profile
- Animated shimmer effect
- Makes wait feel shorter (perceived performance)
- Communicates that something is happening

## Skeleton Structure

Match the profile layout:
- Header skeleton (name, title placeholders)
- Summary skeleton (paragraph lines)
- Experience skeletons (card shapes)
- Skills skeleton (tag shapes)

## Behavior

1. Show immediately after upload starts
2. Animate continuously during processing
3. Poll for completion status (or use subscription)
4. Smooth transition when ready (fade skeleton, reveal content)

## Technical

- Use shadcn/ui Skeleton component
- Create SkeletonProfile component
- Polling interval: 1-2 seconds
- GraphQL query for job status

## Checklist

- [ ] Create SkeletonProfile component
- [ ] Add shimmer animation
- [ ] Implement status polling
- [ ] Handle extraction failure gracefully
- [ ] Smooth reveal animation when complete