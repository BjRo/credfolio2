---
# credfolio2-kbx9
title: Profile page shows stale data after redirect from upload
status: in-progress
type: bug
created_at: 2026-01-28T16:58:30Z
updated_at: 2026-01-28T16:58:30Z
---

## Problem
After uploading a resume, the user is redirected to the profile page but the extracted data isn't displayed. A page reload shows the data correctly.

## Likely Cause
Race condition - the frontend redirects before the extraction job completes, or doesn't poll/wait for the extraction to finish.

## Investigation
- [ ] Check the upload flow and redirect logic
- [ ] Check if the profile page waits for extraction to complete
- [ ] Determine if polling or real-time updates are needed

## Fix
- [ ] Implement solution
- [ ] Test the upload-to-profile flow