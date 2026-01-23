---
# credfolio2-aevf
title: Resume extraction schema and prompt
status: draft
type: task
priority: normal
created_at: 2026-01-23T16:27:30Z
updated_at: 2026-01-23T16:37:22Z
parent: credfolio2-6oza
blocking:
    - credfolio2-he9y
    - credfolio2-jijw
---

Define the structured output schema for resume parsing and create the extraction prompt.

## Schema Fields

### Personal Info
- name (required)
- email
- phone
- location (city, country)
- linkedin_url
- website

### Summary
- headline (short title, e.g., "Senior Software Engineer")
- summary (longer about text)

### Work Experience (array)
- company_name
- job_title
- location
- start_date
- end_date (null if current)
- is_current (boolean)
- description (responsibilities, achievements)
- highlights (array of bullet points)

### Education (array)
- institution
- degree
- field_of_study
- start_date
- end_date
- gpa (optional)
- highlights

### Skills
- technical (array of strings)
- soft (array of strings)
- languages (array with proficiency)

### Certifications (array)
- name
- issuer
- date
- url

## Checklist

- [ ] Define Go structs for resume data
- [ ] Add GraphQL types
- [ ] Create extraction prompt
- [ ] Test with sample resumes
- [ ] Handle partial/missing data gracefully