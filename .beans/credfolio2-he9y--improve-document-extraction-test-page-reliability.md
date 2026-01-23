---
# credfolio2-he9y
title: Improve document extraction test page reliability
status: todo
type: task
priority: normal
created_at: 2026-01-23T15:56:45Z
updated_at: 2026-01-23T15:57:09Z
parent: credfolio2-5r8s
---

Follow-up improvements for the /extract-test page after initial LLM gateway work.

## Context
The extraction endpoint works but has reliability issues with long-running requests (~60s for large PDFs).

## Issues to Address

1. **Browser timeout/CORS issues** - Long requests (60s+) cause browser-level failures even with server-side proxy
2. **No progress indication** - Users see only a spinner during the long wait
3. **Debug console.logs left in code** - Clean up temporary debugging statements

## Potential Solutions

- Add elapsed time indicator during extraction
- Consider async/polling pattern: POST returns job ID, client polls for completion
- Add streaming response support for real-time progress
- Investigate why browser still times out with Next.js proxy

## Files Involved
- src/frontend/src/app/extract-test/page.tsx
- src/frontend/src/app/api/extract/route.ts
- src/backend/internal/handler/extract.go