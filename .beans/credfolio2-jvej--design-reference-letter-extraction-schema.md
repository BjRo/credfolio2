---
# credfolio2-jvej
title: Design reference letter extraction schema
status: todo
type: task
priority: high
created_at: 2026-01-22T07:53:47Z
updated_at: 2026-01-22T07:53:47Z
parent: credfolio2-tmlf
---

Define the JSON schema for extracted reference letter data.

## Requirements

Based on product decisions:
- Must capture author details (name, title, organization, relationship to candidate)
- Must extract key skills and qualities (technical skills, soft skills, traits)
- Must identify specific accomplishments and examples
- Must assess recommendation strength/sentiment

## Design Considerations

- Schema should support aggregation across multiple letters
- Skills need normalization (e.g., "JavaScript" vs "JS" → same skill)
- Consider confidence scores for extracted fields
- Handle partial/missing information gracefully

## Checklist

- [ ] Define TypeScript types for extracted data
- [ ] Define Go structs matching the schema
- [ ] Add GraphQL types to schema
- [ ] Document field extraction guidelines for LLM prompt design
- [ ] Create migration for new database tables if needed

## Example Structure (starting point)

```typescript
interface ExtractedReferenceLetter {
  author: {
    name: string;
    title?: string;
    organization?: string;
    relationship: string; // "manager", "colleague", "professor", etc.
  };
  skills: Array<{
    name: string;
    category: "technical" | "soft" | "domain";
    mentions: number; // how many times referenced
  }>;
  qualities: Array<{
    trait: string;
    evidence?: string; // supporting quote/example
  }>;
  accomplishments: Array<{
    description: string;
    impact?: string;
  }>;
  recommendation: {
    strength: "strong" | "moderate" | "reserved";
    sentiment: number; // -1 to 1
    keyQuote?: string;
  };
  metadata: {
    extractedAt: string;
    confidence: number;
    sourceFileId: string;
  };
}
```