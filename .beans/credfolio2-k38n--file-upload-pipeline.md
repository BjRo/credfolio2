---
# credfolio2-k38n
title: File Upload Pipeline
status: todo
type: epic
priority: normal
created_at: 2026-01-20T11:24:19Z
updated_at: 2026-01-21T14:24:38Z
parent: credfolio2-tikg
blocking:
    - credfolio2-tmlf
---

Enable users to upload reference letter documents (PDF, DOCX, TXT) and store them for processing.

## Goals
- Accept file uploads through the backend API
- Store files in MinIO (S3-compatible)
- Queue files for async LLM processing via River
- Basic upload UI in Next.js

## Checklist
- [ ] Add River job queue dependency and setup
- [ ] Create file storage abstraction layer (supports MinIO/S3/local)
- [ ] Implement file upload GraphQL mutation
- [ ] Create file metadata database table and model
- [ ] Implement River job for document processing
- [ ] Build upload UI component with drag-and-drop
- [ ] Add file type validation (PDF, DOCX, TXT only)
- [ ] Show upload progress and status