---
# credfolio2-wh63
title: Fix HTTP WriteTimeout causing context canceled errors during LLM detection
status: in-progress
type: bug
priority: high
created_at: 2026-02-05T21:59:46Z
updated_at: 2026-02-05T22:26:05Z
---

## Summary

The `detectDocumentContent` mutation synchronously uploads a file AND calls LLM APIs (text extraction + classification), blocking the HTTP request for 30-120+ seconds. Next.js rewrite proxy has a 30s default timeout, killing the connection → "anthropic: context canceled".

## Fix

Split into async pattern matching existing `processDocument` flow:
1. **Upload mutation** (fast): store file, queue detection job, return fileId
2. **Background worker** (async): extract text + classify via LLM, save results to DB
3. **Polling query** (fast): frontend polls for detection status/results

## Checklist

### Backend
- [x] Migration: add detection_status, detection_result, detection_error to files table
- [x] Domain: DetectionStatus type, File struct fields, FileRepository.Update, JobEnqueuer.EnqueueDocumentDetection
- [x] File repository: implement Update method
- [x] Detection worker (TDD): background job for text extraction + detection
- [x] Queue + server wiring: EnqueueDocumentDetection, register worker
- [x] GraphQL schema: rename mutation, add polling query + types, regenerate
- [x] Resolver: async upload mutation + detection status query + tests

### Frontend
- [x] Types: update DocumentUploadProps, add DetectionProgressProps
- [x] DocumentUpload (TDD): change to uploadForDetection, return fileId only
- [x] DetectionProgress (TDD): new polling component
- [x] UploadFlow: add "detect" step, wire components
- [x] Cleanup: revert proxyTimeout and WriteTimeout workarounds

### Definition of Done
- [x] Tests written (TDD: write tests before implementation)
- [x] \`pnpm lint\` passes with no errors
- [x] \`pnpm test\` passes with no failures
- [x] Visual verification with agent-browser (for UI changes)
- [x] All other checklist items above are completed
- [ ] Branch pushed and PR created for human review