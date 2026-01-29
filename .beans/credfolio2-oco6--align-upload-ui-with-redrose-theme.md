---
# credfolio2-oco6
title: Align upload UI with red/rose theme
status: in-progress
type: task
created_at: 2026-01-29T11:26:21Z
updated_at: 2026-01-29T11:26:21Z
---

The upload dialog and page currently use blue/green colors while the rest of the app uses a red/rose suite. Need to update the upload components to match the consistent theme.

## Changes Made

- [x] Updated `ResumeUpload.tsx` - replaced blue/green colors with `primary`, `muted`, `foreground`, `destructive` semantic colors
- [x] Updated `FileUpload.tsx` - same semantic color updates
- [x] Updated `upload/page.tsx` - semantic colors for page background, text, and status badges
- [x] Updated `upload-resume/page.tsx` - semantic colors for page styling
- [x] Updated `page.tsx` (home) - spinner and text colors now use theme

## Definition of Done
- [x] Tests written (TDD: write tests before implementation)
- [x] `pnpm lint` passes with no errors
- [x] `pnpm test` passes with no failures
- [x] Visual verification with agent-browser (for UI changes)
- [x] All other checklist items above are completed
- [x] Branch pushed and PR created for human review

## PR
https://github.com/BjRo/credfolio2/pull/39
