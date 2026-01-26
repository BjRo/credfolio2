---
# credfolio2-umxd
title: Profile Display (LinkedIn-style)
status: completed
type: epic
priority: normal
created_at: 2026-01-23T16:26:56Z
updated_at: 2026-01-26T11:39:11Z
parent: credfolio2-dwid
---

Display extracted profile data in a polished, LinkedIn-inspired layout.

## Design Goals

- Clean, professional appearance
- Work experience as timeline or cards
- Skills displayed as categorized tags
- Personal summary prominent at top
- Mobile-responsive

## Components

### Header Section
- Name, headline/title
- Contact info (if extracted)
- Summary/about section

### Experience Section
- Company, role, dates
- Responsibilities/achievements
- Expandable details

### Skills Section
- Categorized skill tags (technical, soft, domain)
- Possibly with proficiency indicators

### Education Section (if present)
- Institution, degree, dates

## Technical Notes

- Use shadcn/ui components (already set up)
- Ensure fast initial render
- Support deep linking to sections