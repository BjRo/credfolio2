---
# credfolio2-ze4o
title: Document processing River job
status: completed
type: task
priority: normal
created_at: 2026-01-22T09:38:56Z
updated_at: 2026-01-22T10:01:53Z
parent: credfolio2-k38n
---

Implement River job that queues uploaded files for LLM processing.

## Goals
- Create job type for document processing
- Enqueue job when file is uploaded
- Update reference letter status during processing
- Handle job failures gracefully

## Checklist
- [x] Define DocumentProcessingJob struct (DocumentProcessingArgs and DocumentProcessingWorker)
- [x] Implement job worker
- [x] Enqueue job after successful upload
- [x] Add status update logic
- [x] Write job tests