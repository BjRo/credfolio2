---
# credfolio2-tmlf
title: Reference Letter Processing
status: draft
type: epic
priority: normal
created_at: 2026-01-20T11:24:33Z
updated_at: 2026-01-22T07:54:35Z
parent: credfolio2-tikg
blocking:
    - credfolio2-jijw
---

Extract structured data from reference letters using LLMs and display results.

## Goals
- Define schema for extracted reference letter data
- Use LLM to extract structured information
- Store extracted data in PostgreSQL
- Display extraction results in UI

## Data to Extract

Based on product decisions (2026-01-21):

**Author Details**
- Name, title, organization
- Relationship to candidate (manager, colleague, professor, etc.)

**Key Skills & Qualities**
- Technical skills (with category)
- Soft skills and personality traits
- Evidence/examples supporting each quality

**Accomplishments**
- Specific achievements and projects cited
- Impact descriptions where mentioned

**Recommendation**
- Overall strength (strong/moderate/reserved)
- Sentiment score
- Key endorsement quotes

**Metadata**
- Confidence scores for extracted fields
- Source file reference

**Multi-letter Aggregation**
- Schema must support combining data from multiple letters into a unified profile
- Skills need normalization (e.g., "JavaScript" vs "JS" → same skill)

## Checklist
- [ ] Define reference letter data model/schema
- [ ] Create database tables for extracted data
- [ ] Design LLM prompt for structured extraction
- [ ] Implement extraction River job
- [ ] Create GraphQL queries for reference letters
- [ ] Build UI to display extracted data
- [ ] Show extraction status (pending, processing, done, failed)
- [ ] Handle extraction errors gracefully