---
# credfolio2-tmlf
title: Reference Letter Processing
status: draft
type: epic
created_at: 2026-01-20T11:24:33Z
updated_at: 2026-01-20T11:24:33Z
parent: credfolio2-tikg
---

Extract structured data from reference letters using LLMs and display results.

## Goals
- Define schema for extracted reference letter data
- Use LLM to extract structured information
- Store extracted data in PostgreSQL
- Display extraction results in UI

## Data to Extract
- Company name and industry
- Job title/role
- Employment dates
- Key responsibilities
- Skills and technologies mentioned
- Methodologies (Agile, Scrum, etc.)
- Testimonials/praise quotes
- Overall sentiment

## Checklist
- [ ] Define reference letter data model/schema
- [ ] Create database tables for extracted data
- [ ] Design LLM prompt for structured extraction
- [ ] Implement extraction River job
- [ ] Create GraphQL queries for reference letters
- [ ] Build UI to display extracted data
- [ ] Show extraction status (pending, processing, done, failed)
- [ ] Handle extraction errors gracefully