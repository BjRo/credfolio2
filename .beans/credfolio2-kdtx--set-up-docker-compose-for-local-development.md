---
# credfolio2-kdtx
title: Set up Docker Compose for local development
status: todo
type: task
priority: normal
created_at: 2026-01-20T11:26:05Z
updated_at: 2026-01-20T11:26:38Z
parent: credfolio2-jpin
blocking:
    - credfolio2-3gnq
---

Create Docker Compose configuration with PostgreSQL and MinIO services.

## Requirements
- PostgreSQL 16 container with health check
- MinIO container for S3-compatible storage
- Shared network for service communication
- Volume mounts for data persistence
- Environment variable configuration

## Acceptance Criteria
- Running `docker compose up` starts all services
- PostgreSQL is accessible on localhost:5432
- MinIO console is accessible on localhost:9001
- Data persists across restarts

## Technical Notes
- Use official postgres:16 image
- Use official minio/minio image
- Create .env.example with required variables