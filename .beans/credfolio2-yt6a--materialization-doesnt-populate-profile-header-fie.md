---
# credfolio2-yt6a
title: Materialization doesn't populate profile header fields
status: in-progress
type: bug
created_at: 2026-02-06T09:40:33Z
updated_at: 2026-02-06T09:40:33Z
---

## Problem
MaterializeResumeData calls GetOrCreateByUserID to get/create the profile, but never copies the extracted name, email, phone, location, summary from ResumeExtractedData into the Profile row. This results in an empty profile header showing "Unknown".

## Fix
After GetOrCreateByUserID, update the profile with extracted header fields (only if not already set, to respect user edits).

## Checklist
- [ ] Update MaterializeResumeData to populate profile header fields from extracted data
- [ ] Add test for header field population
- [ ] pnpm lint passes
- [ ] pnpm test passes