---
# credfolio2-k38n
title: File Upload Pipeline
status: in-progress
type: epic
priority: normal
created_at: 2026-01-20T11:24:19Z
updated_at: 2026-01-22T09:37:20Z
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
- [x] Add River job queue dependency and setup
- [x] Create file storage abstraction layer (supports MinIO/S3/local)
- [x] Implement file upload GraphQL mutation
- [x] Create file metadata database table and model (already exists)
- [x] Implement River job for document processing
- [x] Build upload UI component with drag-and-drop
- [x] Add file type validation (PDF, DOCX, TXT only)
- [x] Show upload progress and status