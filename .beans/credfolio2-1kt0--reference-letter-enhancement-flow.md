---
# credfolio2-1kt0
title: Reference Letter Credibility System
status: completed
type: epic
priority: normal
created_at: 2026-01-23T16:26:58Z
updated_at: 2026-01-30T13:36:22Z
parent: credfolio2-dwid
---

Add reference letters to build credibility through corroboration and testimonials.

## Core Concept

Reference letters **validate** existing profile data. The goal is **overlap**:
- Skills mentioned in both resume AND reference letters = high credibility
- Experiences confirmed by references = verified claims
- Multiple references confirming the same skill = stronger signal

## User Experience

### Upload Flow
1. From profile view, click "Add Reference Letter"
2. Upload letter (PDF, image, DOCX)
3. See processing indicator
4. View "Validation Preview" showing:
   - **Corroborations**: Existing skills/experiences that will be validated (with quotes)
   - **New Testimonials**: Full quotes to add to profile
   - **Discovered Skills**: Skills your reference mentioned that aren't in your profile yet
5. Select which validations to apply (granular checkboxes)
6. Click "Apply Selected" to save
7. Return to profile with updated credibility indicators

### Profile View (Default)
- Subtle credibility indicators (dots) on skills and experiences
- Visual mark on experiences that have validations
- Credibility score bar: "72% backed by references"
- Testimonials section with full quotes and attribution

### Profile View (Hover/Focus)
- Hover any validated item -> popover showing:
  - Which sources mention this item
  - Quote snippets from each reference
  - Author name and relationship

### Credibility Score Breakdown (Click)
- Expandable panel showing:
  - % of skills validated by references
  - % of experiences validated
  - List of sources with contribution counts
  - CTA: "Add more reference letters to increase credibility"

## What Gets Extracted

1. **Author info**: Name, title, company, relationship to candidate
2. **Testimonials**: Full quotes suitable for display
3. **Skill mentions**: Skills referenced + surrounding quote context
4. **Experience mentions**: References to roles/companies + quotes
5. **Discovered skills**: Skills mentioned that aren't in the profile

## Data Model

### Tables
- `reference_letters` - Uploaded documents with extracted_data JSONB
- `testimonials` - Full quotes with attribution
- `skill_validations` - Links profile_skill to reference_letter with quote_snippet
- `experience_validations` - Links profile_experience to reference_letter with quote_snippet

### Credibility Calculation
- Per-skill: COUNT of validations
- Per-experience: COUNT of validations
- Overall: (validated_skills + validated_experiences) / (total_skills + total_experiences)

## Visual Design

### Credibility Indicators
- **Skills**: Single dot = 1 source, double dot = 2 sources, triple dot = 3+ sources
- **Experiences**: Subtle badge/icon when has validations
- **Hover**: Popover with source list, quotes, "View full" link

### Credibility Score Bar
- Progress bar with percentage
- Click to expand breakdown
- Shows path to 100% (encourages more uploads)
