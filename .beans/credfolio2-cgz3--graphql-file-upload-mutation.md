---
# credfolio2-cgz3
title: GraphQL file upload mutation
status: in-progress
type: task
priority: normal
created_at: 2026-01-22T09:38:54Z
updated_at: 2026-01-22T09:41:46Z
parent: credfolio2-k38n
---

Add GraphQL mutation for uploading files with validation.

## Goals
- Accept multipart file uploads via GraphQL
- Validate file types (PDF, DOCX, TXT only)
- Store file in MinIO and metadata in database
- Return file metadata on success

## Checklist
- [ ] Add uploadFile mutation to GraphQL schema
- [ ] Implement file type validation
- [ ] Implement mutation resolver
- [ ] Wire up storage and repository
- [ ] Write resolver tests