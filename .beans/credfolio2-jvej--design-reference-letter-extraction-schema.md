---
# credfolio2-jvej
title: Design reference letter extraction schema
status: completed
type: task
priority: high
created_at: 2026-01-22T07:53:47Z
updated_at: 2026-01-22T09:34:14Z
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

- [x] Define TypeScript types for extracted data
- [x] Define Go structs matching the schema
- [x] Add GraphQL types to schema
- [x] Document field extraction guidelines for LLM prompt design
- [x] Create migration for new database tables if needed (not needed - using existing JSONB column)

## Implementation Summary

The extraction schema was implemented using a GraphQL-first approach:

1. **Go domain types**: `src/backend/internal/domain/extraction.go`
   - Type-safe structs with JSON tags
   - Enums for relationship, skill category, recommendation strength

2. **GraphQL schema**: `src/backend/internal/graphql/schema/schema.graphqls`
   - Full type definitions with documentation
   - Changed `extractedData: JSON` to `extractedData: ExtractedLetterData`

3. **TypeScript types**: Auto-generated via `pnpm codegen`
   - Types flow from GraphQL schema automatically

4. **LLM prompt documentation**: `src/backend/internal/domain/extraction_schema.md`
   - Full schema reference
   - Confidence scoring guidelines
   - Example output JSON

**No database migration needed** - the existing `extracted_data JSONB` column stores the typed data.

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