---
# credfolio2-2ex3
title: Testimonials UI Enhancements
status: in-progress
type: epic
created_at: 2026-01-30T13:46:56Z
updated_at: 2026-01-30T13:46:56Z
---

Enhance the testimonials section ("What Others Say") on the profile page with improved trust signals, external linking, and better organization.

## Context
Currently, testimonials are displayed as separate cards even when multiple quotes come from the same author. Each testimonial links to its source reference letter internally but this isn't surfaced in the UI. Author information is plain text with no external verification.

## Goals
- **Build trust**: Show provenance of testimonials (where they came from)
- **Enable verification**: Allow linking to author's LinkedIn profile
- **Reduce clutter**: Group multiple quotes from same author

## Child Features
1. Source badge - Show which reference letter each testimonial came from
2. LinkedIn integration - Allow linking author profiles to LinkedIn
3. Collapse by author - Group testimonials from the same author

## Technical Context
- Testimonial → ReferenceLetter relation already exists in GraphQL schema
- Author fields: `authorName`, `authorTitle`, `authorCompany`, `relationship`
- Frontend component: `src/frontend/src/components/profile/TestimonialsSection.tsx`