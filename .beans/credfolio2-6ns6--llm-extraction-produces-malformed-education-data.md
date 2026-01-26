---
# credfolio2-6ns6
title: LLM extraction produces malformed education data
status: todo
type: bug
created_at: 2026-01-26T11:16:38Z
updated_at: 2026-01-26T11:16:38Z
---

The resume extraction LLM sometimes produces corrupted data with JSON syntax characters embedded in fields. Example: startDate contains ')"}}],' and gpa contains '-degree-bachelorship_of Science'. This corrupted data gets stored and displayed incorrectly. Need to add post-extraction validation to catch and reject malformed data.