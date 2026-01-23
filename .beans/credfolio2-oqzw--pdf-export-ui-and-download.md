---
# credfolio2-oqzw
title: PDF export UI and download
status: draft
type: feature
created_at: 2026-01-23T16:29:33Z
updated_at: 2026-01-23T16:29:33Z
parent: credfolio2-nix0
---

UI for triggering and downloading PDF export.

## Trigger

- "Export PDF" button on profile page
- Dropdown for options (future: template selection)

## Flow

1. User clicks "Export PDF"
2. Button shows loading state
3. PDF generates server-side
4. Browser downloads file
5. Success toast

## Preview (Nice-to-have)

- Modal showing PDF preview
- Confirm before download
- Uses react-pdf or iframe

## Download Behavior

- Filename: "{name}-resume.pdf"
- Direct download (Content-Disposition: attachment)
- Or: Opens in new tab as PDF

## Checklist

- [ ] Add "Export PDF" button to profile
- [ ] Implement loading state during generation
- [ ] Trigger download on completion
- [ ] Handle generation errors
- [ ] Add success notification
- [ ] (Optional) Add preview modal