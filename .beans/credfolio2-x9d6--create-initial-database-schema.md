---
# credfolio2-x9d6
title: Create initial database schema
status: todo
type: task
created_at: 2026-01-20T11:26:08Z
updated_at: 2026-01-20T11:26:08Z
parent: credfolio2-jpin
---

Design and create the initial PostgreSQL schema for core entities.

## Tables to Create
1. **profiles** - User profile metadata
   - id (UUID, PK)
   - name, headline, summary
   - created_at, updated_at

2. **reference_files** - Uploaded document metadata  
   - id (UUID, PK)
   - profile_id (FK)
   - filename, content_type, storage_key
   - status (pending, processing, processed, failed)
   - created_at, updated_at

3. **positions** - Extracted job positions
   - id (UUID, PK)
   - profile_id (FK)
   - reference_file_id (FK)
   - company_name, job_title
   - start_date, end_date
   - responsibilities (text[])
   - created_at, updated_at

4. **skills** - Extracted skills
   - id (UUID, PK)
   - profile_id (FK)
   - name, category
   - proficiency_level (optional)

5. **testimonials** - Extracted praise/quotes
   - id (UUID, PK)
   - position_id (FK)
   - quote_text, context
   - sentiment_score (optional)

## Acceptance Criteria
- All tables created via migrations
- Foreign key relationships established
- Appropriate indexes added
- Can rollback cleanly