---
# credfolio2-jijw
title: End-to-end vertical slice integration
status: draft
type: task
created_at: 2026-01-20T15:31:10Z
updated_at: 2026-01-20T15:31:10Z
parent: credfolio2-tikg
---

Wire together all components for the first vertical slice: upload a reference letter PDF, extract structured data via LLM, display results.

## Integration Points
- Frontend upload form → Backend upload endpoint
- Upload endpoint → MinIO storage
- Processing trigger → River job → LLM Gateway
- LLM extraction → Structured data using defined schema
- Extracted data → PostgreSQL storage
- GraphQL API → Serve structured results to frontend
- Frontend → Display extracted profile data

## Acceptance Criteria
- [ ] User can upload a PDF file through the UI
- [ ] File is stored in MinIO with metadata in PostgreSQL
- [ ] Background job triggers LLM extraction
- [ ] Extracted **structured data** (author, skills, accomplishments, recommendation) is saved
- [ ] User sees extracted profile data in the UI (not raw text)
- [ ] Error states are handled (upload failure, extraction failure)

## Demo Scenario

1. Navigate to upload page
2. Select a reference letter PDF
3. See upload progress
4. After processing, view extracted:
   - Author information (name, title, relationship)
   - Skills identified (categorized)
   - Key accomplishments
   - Recommendation strength and key quote

## Dependencies
This task should be worked on after the individual epics are complete:
- Infrastructure Foundation (credfolio2-jpin)
- File Upload Pipeline (credfolio2-k38n)
- LLM Gateway Service (credfolio2-5r8s)
- GraphQL API (credfolio2-zbqk)