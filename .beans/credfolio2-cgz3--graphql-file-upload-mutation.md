---
# credfolio2-cgz3
title: GraphQL file upload mutation
status: completed
type: task
priority: normal
created_at: 2026-01-22T09:38:54Z
updated_at: 2026-01-22T09:48:33Z
parent: credfolio2-k38n
---

Add GraphQL mutation for uploading files with validation.

## Goals
- Accept multipart file uploads via GraphQL
- Validate file types (PDF, DOCX, TXT only)
- Store file in MinIO and metadata in database
- Return file metadata on success

## Checklist
- [x] Add uploadFile mutation to GraphQL schema
- [x] Implement file type validation
- [x] Implement mutation resolver
- [x] Wire up storage and repository
- [x] Write resolver tests (updated existing tests)