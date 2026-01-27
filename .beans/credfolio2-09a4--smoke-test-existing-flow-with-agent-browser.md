---
# credfolio2-09a4
title: Smoke-test existing flow with agent-browser
status: completed
type: task
priority: normal
created_at: 2026-01-27T11:07:20Z
updated_at: 2026-01-27T11:23:17Z
parent: credfolio2-oo6v
---

Before starting implementation on oo6v (Education entry form), verify we can execute the existing resume upload and profile editing flow using agent-browser. This validates the dev environment is working and helps understand the UI patterns to follow.

## Checklist
- [x] Start dev servers (frontend + backend)
- [x] Upload fixture resume (CV_TEMPLATE_0004.pdf) via the UI
- [x] Navigate to the profile/edit page
- [x] Verify existing experience form works
- [x] Document findings for education form implementation

## Findings

### Upload Flow
- Upload page at `/upload-resume` with drag-and-drop area
- File uploads trigger backend processing with LLM extraction
- Auto-redirects to `/profile/{id}` when processing completes

### Profile Page
- Shows: header (name, email, phone, location), summary, work experience, education, skills
- Work Experience section has "+" add button and "..." menu on each entry
- Education section shows parsed data but has NO add/edit/delete controls yet (this is what oo6v adds)
- Skills shown as tags

### Work Experience Form Pattern (to replicate for Education)
- Modal dialog with title "Add Work Experience" / description text
- Fields: Company* (text), Job Title* (text), Location (text), Start Date (month+year dropdowns), End Date (month+year dropdowns) with "I currently work here" checkbox, Description (textarea), Key Achievements (dynamic list with add/remove)
- Cancel + Submit buttons
- Edit mode pre-fills the same form with existing data
- "More actions" dropdown on each entry shows "Edit" option

### Issues Noticed
- Education date display is broken: shows `-01-01endDate2018-01-01 - Present` instead of proper dates
- Edit form title says "Add Work Experience" even in edit mode (minor bug)