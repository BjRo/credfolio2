---
# credfolio2-jijw
title: End-to-end vertical slice integration
status: draft
type: task
created_at: 2026-01-20T15:31:10Z
updated_at: 2026-01-20T15:31:10Z
parent: credfolio2-tikg
---

Wire together all components for the first vertical slice: upload a reference letter, extract text via LLM, display raw results.

## Integration Points
- Frontend upload form → Backend upload endpoint
- Upload endpoint → MinIO storage
- Processing trigger → LLM Gateway for text extraction
- Extracted data → Database storage
- GraphQL API → Serve results to frontend
- Frontend → Display extracted text

## Acceptance Criteria
- [ ] User can upload a PDF/image of a reference letter
- [ ] File is stored in MinIO
- [ ] LLM extracts text from the document
- [ ] Extracted text is stored in the database
- [ ] Frontend displays the raw extracted text
- [ ] Error states are handled gracefully

## Dependencies
This task should be worked on after the individual epics are complete:
- Infrastructure Foundation (credfolio2-jpin)
- File Upload Pipeline (credfolio2-k38n)
- LLM Gateway Service (credfolio2-5r8s)
- GraphQL API (credfolio2-zbqk)