---
# credfolio2-i03i
title: PDF generation service
status: draft
type: feature
created_at: 2026-01-23T16:29:32Z
updated_at: 2026-01-23T16:29:32Z
parent: credfolio2-nix0
---

Backend service to generate PDF from profile data.

## Approach Options

1. **Go PDF library** (gofpdf, pdfcpu) - Direct PDF generation
2. **HTML to PDF** (chromedp, wkhtmltopdf) - Render HTML template
3. **External service** (DocRaptor, etc.) - API call

Recommendation: HTML to PDF via chromedp
- Easier styling with HTML/CSS
- Chrome rendering = consistent results
- chromedp already works headlessly

## PDF Structure

1. Header: Name, contact, summary
2. Experience section
3. Education section
4. Skills section
5. Optional: Testimonials/quotes section

## API

```
POST /api/profile/{id}/export/pdf
Response: application/pdf binary
```

Or GraphQL mutation returning download URL:
```graphql
mutation GeneratePDF(profileId: ID!): String! # URL
```

## Checklist

- [ ] Evaluate PDF generation approach
- [ ] Create HTML template for resume
- [ ] Style template for print
- [ ] Set up chromedp or chosen tool
- [ ] Implement PDF generation endpoint
- [ ] Handle page breaks properly
- [ ] Test with various profile sizes