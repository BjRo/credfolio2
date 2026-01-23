---
# credfolio2-nix0
title: PDF Export
status: draft
type: epic
created_at: 2026-01-23T16:27:01Z
updated_at: 2026-01-23T16:27:01Z
parent: credfolio2-dwid
---

Generate a professional PDF resume/CV from the profile data.

## Design Goals

- Clean, ATS-friendly format
- Professional typography
- Proper page breaks
- Print-optimized colors

## PDF Content

Based on profile data:
- Header with name, contact, summary
- Experience section with details
- Skills section (categorized or flat)
- Education section
- Optional: Testimonials/quotes section

## Generation Approach

Options to evaluate:
1. **Server-side** - Go generates PDF (e.g., gofpdf, chromedp)
2. **Client-side** - React-pdf or similar
3. **Headless Chrome** - Render HTML to PDF

Recommendation: Start with server-side for reliability

## Features

- Preview before download
- Multiple template options (future)
- Customizable sections to include
- Download as PDF button
- Optional: Direct share link