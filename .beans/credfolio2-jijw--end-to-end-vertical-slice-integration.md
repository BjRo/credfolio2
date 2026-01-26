---
# credfolio2-jijw
title: 'Minimal E2E integration: Resume upload to profile view'
status: completed
type: task
priority: normal
created_at: 2026-01-20T15:31:10Z
updated_at: 2026-01-26T10:01:01Z
parent: credfolio2-tikg
---

Wire together all components for the first working end-to-end flow: **Resume upload → LLM extraction → Profile display**.

## Summary

- **Document**: Resume/CV (PDF or DOCX)
- **Extraction**: Standard profile (name, contact, work experience, education, skills)
- **UX**: Upload → async processing with polling → auto-redirect to /profile
- **Auth**: Demo user (existing hardcoded user, auth out of scope)

## Technical Approach

### Data Model

Create a new `resumes` table (separate from existing `reference_letters`):
- `id`, `user_id`, `file_id` (links to uploaded file)
- `status`: pending → processing → completed/failed
- `extracted_data` (JSONB): structured profile data
- `error_message`: capture failures
- Timestamps

### Extraction Schema (ResumeExtractedData)

```go
type ResumeExtractedData struct {
    Name        string
    Email       string
    Phone       string
    Location    string
    Summary     string
    Experience  []WorkExperience  // company, title, dates, description
    Education   []Education       // institution, degree, field, dates
    Skills      []string
    ExtractedAt time.Time
    Confidence  float64
}
```

### Processing Flow

1. Frontend uploads file via existing `uploadFile` mutation
2. Backend stores in MinIO, creates File record
3. **New**: Create Resume record (status: pending), enqueue ResumeProcessingJob
4. Frontend polls `resume(id)` query for status changes
5. Job worker: download file → extract text → call LLM with structured output → save to `extracted_data`
6. Frontend detects `completed` status → auto-redirect to `/profile`

### API Changes

**New GraphQL types:**
- `Resume` type with `id`, `status`, `extractedData`, `createdAt`, `updatedAt`
- `ResumeExtractedData` type (matches Go struct above)
- `WorkExperience`, `Education` types

**New mutations:**
- `uploadResume(userId: ID!, file: Upload!): UploadResumeResult!` (or extend existing uploadFile)

**New queries:**
- `resume(id: ID!): Resume`
- `resumes(userId: ID!): [Resume!]!`

## Checklist

### Backend: Data Model
- [x] Create `ResumeExtractedData` struct in `domain/resume.go`
- [x] Create `Resume` domain model with status enum
- [x] Create migration: `resumes` table with JSONB `extracted_data`
- [x] Add Resume repository interface and Bun implementation

### Backend: GraphQL API
- [x] Add Resume types to GraphQL schema
- [x] Add `uploadResume` mutation (creates File + Resume + enqueues job)
- [x] Add `resume(id)` and `resumes(userId)` queries
- [x] Implement resolvers

### Backend: Job Processing
- [x] Create `ResumeProcessingWorker` (River job)
- [x] Implement file download from MinIO in worker
- [x] Create LLM extraction prompt for resumes (structured output)
- [x] Call LLM and parse response into `ResumeExtractedData`
- [x] Update Resume record with extracted data or error

### Frontend: Upload Flow
- [x] Update upload page to call `uploadResume` mutation
- [x] Implement polling for resume status (useQuery with pollInterval or manual)
- [x] Show processing indicator while status is pending/processing
- [x] Auto-redirect to `/profile/{resumeId}` when status is completed

### Frontend: Profile Page
- [x] Create `/profile/[id]/page.tsx` route
- [x] Query resume by ID with extracted data
- [x] Display: name, contact info, summary
- [x] Display: work experience list (company, title, dates, description)
- [x] Display: education list
- [x] Display: skills as tags/chips
- [x] Handle loading and error states

### Integration & Testing
- [ ] Manual E2E test: upload real resume PDF → see profile
- [ ] Verify error handling: upload invalid file, LLM failure
- [ ] Check job retry behavior on transient failures

## Out of Scope

Explicitly NOT part of this task (future work):
- Polished UI/design (skeleton loading, animations)
- Reference letter processing
- Profile editing
- Multiple resumes per user
- Real authentication
- PDF export
- Error retry UI

## Definition of Done

- [ ] Can upload a PDF resume through the UI
- [ ] Processing completes within ~30 seconds
- [ ] Auto-redirected to profile page showing extracted data
- [ ] Name, work experience, education, and skills are displayed
- [ ] Errors are logged (no crash on failure)