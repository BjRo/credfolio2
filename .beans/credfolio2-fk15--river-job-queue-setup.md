---
# credfolio2-fk15
title: River job queue setup
status: todo
type: task
priority: normal
created_at: 2026-01-22T09:38:56Z
updated_at: 2026-01-22T09:38:56Z
parent: credfolio2-k38n
---

Integrate River for background job processing.

## Goals
- Add River dependency and database migrations
- Configure River client in application startup
- Create base job infrastructure

## Checklist
- [ ] Add riverqueue/river dependency
- [ ] Create River schema migration
- [ ] Initialize River client in main.go
- [ ] Add River configuration to config package