---
# credfolio2-nzoc
title: Improve date extraction instructions in LLM prompt
status: completed
type: task
priority: normal
created_at: 2026-01-26T15:23:03Z
updated_at: 2026-01-26T15:23:48Z
---

Extend the resume extraction prompt with clearer instructions for:
1. Extracting dates in ISO format (YYYY-MM-DD)
2. Handling partial dates (expand month-only to include year from context)
3. Treating unparsable dates (return null instead of malformed values)

The issue: dates like '-09-01T00:00Z' are being returned when the year is missing.