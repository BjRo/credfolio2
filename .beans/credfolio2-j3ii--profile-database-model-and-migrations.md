---
# credfolio2-j3ii
title: Profile database model and migrations
status: draft
type: task
priority: normal
created_at: 2026-01-23T16:27:32Z
updated_at: 2026-01-23T16:37:22Z
parent: credfolio2-6oza
blocking:
    - credfolio2-he9y
    - credfolio2-jijw
---

Create the database schema for storing extracted profile data.

## Tables

### profiles
- id (UUID, PK)
- user_id (FK, nullable initially)
- name
- email
- phone
- location_city
- location_country
- linkedin_url
- website_url
- headline
- summary
- created_at
- updated_at

### work_experiences
- id (UUID, PK)
- profile_id (FK)
- company_name
- job_title
- location
- start_date
- end_date
- is_current
- description
- highlights (text[] or JSONB)
- display_order
- source_type (resume, reference_letter)
- source_id (file_id that contributed this)
- created_at
- updated_at

### education
- id (UUID, PK)
- profile_id (FK)
- institution
- degree
- field_of_study
- start_date
- end_date
- gpa
- highlights (JSONB)
- display_order
- created_at
- updated_at

### profile_skills
- id (UUID, PK)
- profile_id (FK)
- name
- category (technical, soft, language)
- proficiency_level (for languages)
- source_type
- source_id
- created_at

### certifications
- id (UUID, PK)
- profile_id (FK)
- name
- issuer
- date_obtained
- url
- created_at

## Checklist

- [ ] Create migration files
- [ ] Add indexes for common queries
- [ ] Create Bun models
- [ ] Create repository interfaces and implementations
- [ ] Write repository tests