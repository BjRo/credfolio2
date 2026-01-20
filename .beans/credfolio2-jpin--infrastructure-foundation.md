---
# credfolio2-jpin
title: Infrastructure Foundation
status: in-progress
type: epic
priority: normal
created_at: 2026-01-20T11:24:12Z
updated_at: 2026-01-20T13:51:58Z
parent: credfolio2-tikg
blocking:
    - credfolio2-k38n
    - credfolio2-5r8s
---

Set up the development infrastructure: Docker Compose, database migrations, project structure, and environment configuration.

## Process

**Before starting work on this epic or its tasks:**

1. Follow the **dev-workflow** skill (`/skill dev-workflow`)
2. Use **TDD** for all code (`/skill tdd`)
3. Create a feature branch from up-to-date main
4. Open a PR for human review when done - do NOT merge yourself

## Goals
- Reproducible local development environment via Docker Compose
- PostgreSQL database with migration tooling
- MinIO for S3-compatible file storage
- Clean Architecture folder structure for Go backend
- Environment-based configuration

## Checklist
- [ ] Create Docker Compose with PostgreSQL and MinIO services
- [ ] Set up golang-migrate for database migrations
- [ ] Create initial database schema (users, files, reference_letters tables)
- [ ] Establish Go backend folder structure (cmd, internal, pkg)
- [ ] Configure environment variables loading
- [ ] Update CLAUDE.md with new infrastructure details
- [ ] Add Makefile/scripts for common operations