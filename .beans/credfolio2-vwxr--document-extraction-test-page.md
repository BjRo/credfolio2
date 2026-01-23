---
# credfolio2-vwxr
title: Document extraction test page
status: completed
type: feature
priority: normal
created_at: 2026-01-23T15:07:45Z
updated_at: 2026-01-23T16:04:06Z
---

Build a simple UI page to manually test document extraction via the LLM gateway. Upload a document image/PDF and see the extracted text.

## Goals
- Provide a way to manually test the DocumentExtractor end-to-end
- Useful for demos and exploratory testing
- Simple implementation focused on functionality over polish

## Checklist
- [x] Create backend endpoint for document extraction (POST /api/extract)
- [x] Create frontend test page with file upload
- [x] Display extraction results with loading/error states
- [x] Support common image formats (JPEG, PNG) and PDF
