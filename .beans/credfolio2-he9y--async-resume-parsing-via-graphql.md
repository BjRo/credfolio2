---
# credfolio2-he9y
title: Async resume parsing via GraphQL
status: todo
type: feature
priority: normal
created_at: 2026-01-23T15:56:45Z
updated_at: 2026-01-23T16:37:31Z
parent: credfolio2-6oza
blocking:
    - credfolio2-jijw
---

Implement async document processing for resume/CV parsing with structured output, replacing the synchronous test endpoint.

## Goals

1. **Remove test extraction endpoint** - Document upload should go through GraphQL, not a separate REST endpoint
2. **Async processing with River** - Use River job queue for LLM extraction (avoids timeout issues)
3. **Resume/CV parsing** - Target structured output schema for resume data extraction
4. **Database persistence** - Store parsed resume data in the database

## Implementation Plan

### 1. Define Resume Schema

Create structured output schema for resume parsing:
- Personal info (name, email, phone, location)
- Work experience (company, role, dates, description)
- Education (institution, degree, dates)
- Skills (list of skills, possibly categorized)
- Languages
- Certifications

### 2. Database Changes

- Add `parsed_resumes` table or extend existing schema
- Store structured JSON output from LLM
- Link to original uploaded file

### 3. River Job for Extraction

- Create `ResumeExtractionJob` worker
- Job receives file ID, fetches from storage, calls LLM
- Stores result in database on completion
- Handle errors gracefully (mark as failed, allow retry)

### 4. GraphQL Integration

- Add mutation: `parseResume(fileId: ID!): ResumeParseJob!`
- Add query: `resumeParseJob(id: ID!): ResumeParseJob`
- Add subscription or polling for job status
- Return parsed data when complete

### 5. Cleanup

- Remove `/api/extract` endpoint (backend)
- Remove `/api/extract` Next.js proxy route
- Remove `/extract-test` page (or repurpose for demo)
- Clean up debug console.logs

## Files to Change

**Backend:**
- `src/backend/internal/domain/` - Add resume domain types
- `src/backend/internal/infrastructure/llm/extraction.go` - Update prompt for resume schema
- `src/backend/internal/worker/` - Add resume extraction job
- `src/backend/internal/graphql/schema/` - Add mutations/queries
- `src/backend/internal/handler/extract.go` - Remove (or deprecate)

**Frontend:**
- `src/frontend/src/app/extract-test/` - Remove or repurpose
- `src/frontend/src/app/api/extract/` - Remove

**Database:**
- New migration for parsed resume storage

## Future Work

- Reference letter parsing (different schema, similar pattern)
- Batch processing for multiple documents
- Re-extraction with updated prompts
