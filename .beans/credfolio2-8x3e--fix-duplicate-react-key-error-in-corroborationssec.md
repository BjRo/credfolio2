---
# credfolio2-8x3e
title: Fix duplicate React key error in CorroborationsSection
status: completed
type: bug
priority: normal
created_at: 2026-02-03T11:42:49Z
updated_at: 2026-02-03T11:44:59Z
---

The CorroborationsSection component uses profileExperienceId and profileSkillId as React keys, but the same skill/experience can appear multiple times with different quotes, causing 'Encountered two children with the same key' errors. Fix by using composite keys that include the index.