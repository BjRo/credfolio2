# Reference Letter Extraction Schema

This document defines the JSON schema for extracted reference letter data. The LLM service should output data in this format.

## Overview

When processing a reference letter, the LLM should extract structured data and output it as JSON conforming to the `ExtractedLetterData` schema. The extraction includes:

- **Author details**: Information about the letter writer
- **Skills**: Technical, soft, and domain skills mentioned
- **Qualities**: Personal traits and characteristics
- **Accomplishments**: Specific achievements cited
- **Recommendation**: Overall assessment and sentiment
- **Metadata**: Extraction process information

## Schema Definition

### ExtractedLetterData (root object)

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `author` | ExtractedAuthor | Yes | Details about the letter's author |
| `skills` | ExtractedSkill[] | Yes | Skills mentioned (empty array if none) |
| `qualities` | ExtractedQuality[] | Yes | Qualities/traits (empty array if none) |
| `accomplishments` | ExtractedAccomplishment[] | Yes | Accomplishments cited (empty array if none) |
| `recommendation` | ExtractedRecommendation | Yes | Overall recommendation assessment |
| `metadata` | ExtractionMetadata | Yes | Extraction process metadata |

### ExtractedAuthor

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | The author's full name |
| `title` | string | No | Job title or position |
| `organization` | string | No | Organization where author works |
| `relationship` | enum | Yes | Relationship type (see values below) |
| `relationshipDetails` | string | No | Additional context (e.g., "Direct supervisor for 3 years") |
| `confidence` | float | Yes | Confidence score 0.0-1.0 |

**Relationship values**: `manager`, `colleague`, `professor`, `client`, `mentor`, `other`

### ExtractedSkill

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | Skill name as it appears in the letter |
| `normalizedName` | string | Yes | Lowercase normalized name for aggregation |
| `category` | enum | Yes | `technical`, `soft`, or `domain` |
| `mentions` | int | Yes | Number of times mentioned in letter |
| `context` | string[] | No | Relevant quotes showing skill usage |
| `confidence` | float | Yes | Confidence score 0.0-1.0 |

**Normalization examples**:
- "JavaScript", "JS", "Javascript" → `"javascript"`
- "React.js", "ReactJS", "React" → `"react"`
- "Problem Solving", "problem-solving" → `"problem solving"`

### ExtractedQuality

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `trait` | string | Yes | The quality or trait name |
| `evidence` | string[] | No | Supporting quotes from the letter |
| `confidence` | float | Yes | Confidence score 0.0-1.0 |

### ExtractedAccomplishment

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `description` | string | Yes | What was accomplished |
| `impact` | string | No | Quantifiable impact or result |
| `timeframe` | string | No | When it occurred (e.g., "2024", "Q3 2023") |
| `confidence` | float | Yes | Confidence score 0.0-1.0 |

### ExtractedRecommendation

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `strength` | enum | Yes | `strong`, `moderate`, or `reserved` |
| `sentiment` | float | Yes | Score from -1.0 (negative) to 1.0 (positive) |
| `keyQuotes` | string[] | No | Quotes demonstrating recommendation strength |
| `summary` | string | No | Brief summary of the recommendation |
| `confidence` | float | Yes | Confidence score 0.0-1.0 |

### ExtractionMetadata

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `extractedAt` | ISO8601 | Yes | Extraction timestamp |
| `modelVersion` | string | Yes | LLM model identifier |
| `overallConfidence` | float | Yes | Aggregate confidence 0.0-1.0 |
| `processingTimeMs` | int | No | Processing time in milliseconds |
| `warningsCount` | int | Yes | Number of warnings |
| `warnings` | string[] | No | Warning messages |

## Confidence Score Guidelines

Confidence scores (0.0 to 1.0) should reflect:

- **0.9-1.0**: Information is explicitly stated and unambiguous
- **0.7-0.9**: Information is clearly implied or slightly ambiguous
- **0.5-0.7**: Information is inferred with moderate certainty
- **0.3-0.5**: Information is a reasonable guess with uncertainty
- **< 0.3**: Should not be included; mark as unknown instead

## Example Output

```json
{
  "author": {
    "name": "Dr. Jane Smith",
    "title": "Director of Engineering",
    "organization": "Acme Technologies",
    "relationship": "manager",
    "relationshipDetails": "Direct supervisor for 3 years on the Platform team",
    "confidence": 0.95
  },
  "skills": [
    {
      "name": "JavaScript",
      "normalizedName": "javascript",
      "category": "technical",
      "mentions": 3,
      "context": [
        "Expert knowledge of modern JavaScript",
        "Led the team's transition to ES6+"
      ],
      "confidence": 0.92
    },
    {
      "name": "Leadership",
      "normalizedName": "leadership",
      "category": "soft",
      "mentions": 2,
      "context": [
        "Naturally stepped into leadership roles",
        "Mentored three junior developers"
      ],
      "confidence": 0.88
    }
  ],
  "qualities": [
    {
      "trait": "Self-motivated",
      "evidence": [
        "Proactively identified and resolved technical debt",
        "Never needed to be reminded of deadlines"
      ],
      "confidence": 0.85
    },
    {
      "trait": "Collaborative",
      "evidence": [
        "Always willing to help teammates"
      ],
      "confidence": 0.82
    }
  ],
  "accomplishments": [
    {
      "description": "Led migration from legacy monolith to microservices architecture",
      "impact": "Reduced deployment time by 80% and improved system reliability",
      "timeframe": "2024",
      "confidence": 0.90
    },
    {
      "description": "Designed and implemented real-time notification system",
      "impact": "Handles 1M+ notifications per day",
      "confidence": 0.85
    }
  ],
  "recommendation": {
    "strength": "strong",
    "sentiment": 0.92,
    "keyQuotes": [
      "I highly recommend John without any reservation",
      "One of the best engineers I have had the pleasure of working with"
    ],
    "summary": "Enthusiastic endorsement for senior engineering roles with emphasis on technical leadership",
    "confidence": 0.95
  },
  "metadata": {
    "extractedAt": "2026-01-22T10:30:00Z",
    "modelVersion": "claude-3-opus",
    "overallConfidence": 0.88,
    "processingTimeMs": 2340,
    "warningsCount": 0,
    "warnings": []
  }
}
```

## Handling Missing Information

- Use `null` for optional fields when information is not available
- Use empty arrays `[]` for list fields when no items are found
- Do **not** invent information; only extract what is present
- Lower confidence scores for inferred information
- Add warnings for ambiguous or potentially incorrect extractions

## LLM Prompt Integration

When prompting the LLM for extraction:

1. Provide the full letter text
2. Include this schema as the expected output format
3. Instruct to return valid JSON only
4. Specify confidence scoring guidelines
5. Request skill normalization

Example prompt structure:
```
Extract structured data from the following reference letter.
Output must be valid JSON matching the ExtractedLetterData schema.

[Schema definition]

[Confidence scoring guidelines]

Letter content:
---
[Letter text here]
---

Output the extraction result as JSON:
```
