---
# credfolio2-ue4q
title: Reference letter extraction schema
status: draft
type: task
priority: normal
created_at: 2026-01-23T16:28:38Z
updated_at: 2026-01-23T16:29:49Z
parent: credfolio2-1kt0
blocking:
    - credfolio2-6dty
---

Extraction schema specifically for reference letters (extends existing schema from credfolio2-jvej).

## Schema (from previous work)

- Author details (name, title, org, relationship)
- Skills mentioned (with category)
- Qualities and traits (with evidence)
- Accomplishments cited
- Recommendation strength/sentiment
- Key quotes

## Additions for Enhancement Flow

- Match skills to existing profile skills (normalization)
- Identify truly new vs confirming information
- Confidence scores for extraction quality

## Merge Logic

Define how extracted data merges into profile:
- New skills: Add to profile with source_type=reference_letter
- Existing skills: Mark as "validated", add source
- Testimonials: Add to new testimonials table
- Accomplishments: Add to work experience highlights

## Checklist

- [ ] Review existing schema from credfolio2-jvej
- [ ] Add merge/enhancement metadata
- [ ] Create skill normalization logic
- [ ] Define testimonials table
- [ ] Document merge rules