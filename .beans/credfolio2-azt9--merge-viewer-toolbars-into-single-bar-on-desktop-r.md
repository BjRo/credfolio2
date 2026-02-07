---
# credfolio2-azt9
title: Merge viewer toolbars into single bar on desktop, responsive two-row on mobile
status: completed
type: task
priority: normal
created_at: 2026-02-07T15:11:32Z
updated_at: 2026-02-07T15:21:07Z
parent: credfolio2-klgo
---

## Context

The PDF viewer currently has two separate white bars:
1. Page toolbar (back button + document title/subtitle)
2. Component toolbar (page navigation + zoom controls)

This wastes vertical space on desktop. On mobile, the default zoom (100%) makes the PDF page wider than the viewport.

## Changes

1. **Desktop**: Merge both bars into a single row — back button + title on left, page nav + zoom on right
2. **Mobile**: Keep a two-row layout using responsive Tailwind classes
3. **Mobile**: Lower default zoom factor so the page fits the viewport width

## Checklist

- [ ] Merge the two toolbar bars into one on desktop (responsive breakpoint)
- [ ] Keep two-row layout on mobile
- [ ] Lower default zoom on mobile so page fits viewport
- [ ] Tests pass
- [ ] Visual verification