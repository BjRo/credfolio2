// Package domain contains the core business entities and repository interfaces.
package domain

import "time"

// AuthorRelationship represents the relationship type between letter author and candidate.
type AuthorRelationship string

// Author relationship constants.
const (
	AuthorRelationshipManager   AuthorRelationship = "manager"
	AuthorRelationshipColleague AuthorRelationship = "colleague"
	AuthorRelationshipProfessor AuthorRelationship = "professor"
	AuthorRelationshipClient    AuthorRelationship = "client"
	AuthorRelationshipMentor    AuthorRelationship = "mentor"
	AuthorRelationshipOther     AuthorRelationship = "other"
)

// SkillCategory classifies skills into categories.
type SkillCategory string

// Skill category constants.
const (
	SkillCategoryTechnical SkillCategory = "technical"
	SkillCategorySoft      SkillCategory = "soft"
	SkillCategoryDomain    SkillCategory = "domain"
)

// RecommendationStrength represents the strength of a recommendation.
type RecommendationStrength string

// Recommendation strength constants.
const (
	RecommendationStrengthStrong   RecommendationStrength = "strong"
	RecommendationStrengthModerate RecommendationStrength = "moderate"
	RecommendationStrengthReserved RecommendationStrength = "reserved"
)

// ExtractedAuthor contains author details extracted from a reference letter.
type ExtractedAuthor struct { //nolint:govet // Field ordering prioritizes JSON serialization over memory alignment
	Name                string             `json:"name"`
	Title               *string            `json:"title,omitempty"`
	Organization        *string            `json:"organization,omitempty"`
	Relationship        AuthorRelationship `json:"relationship"`
	RelationshipDetails *string            `json:"relationshipDetails,omitempty"`
	Confidence          float64            `json:"confidence"`
}

// ExtractedSkill represents a skill mentioned in a reference letter.
type ExtractedSkill struct { //nolint:govet // Field ordering prioritizes JSON serialization over memory alignment
	Name           string        `json:"name"`
	NormalizedName string        `json:"normalizedName"`
	Category       SkillCategory `json:"category"`
	Mentions       int           `json:"mentions"`
	Context        []string      `json:"context,omitempty"`
	Confidence     float64       `json:"confidence"`
}

// ExtractedQuality represents a quality or trait from a reference letter.
type ExtractedQuality struct {
	Trait      string   `json:"trait"`
	Evidence   []string `json:"evidence,omitempty"`
	Confidence float64  `json:"confidence"`
}

// ExtractedAccomplishment represents an accomplishment cited in a letter.
type ExtractedAccomplishment struct { //nolint:govet // Field ordering prioritizes JSON serialization over memory alignment
	Description string  `json:"description"`
	Impact      *string `json:"impact,omitempty"`
	Timeframe   *string `json:"timeframe,omitempty"`
	Confidence  float64 `json:"confidence"`
}

// ExtractedRecommendation represents the overall recommendation assessment.
type ExtractedRecommendation struct { //nolint:govet // Field ordering prioritizes JSON serialization over memory alignment
	Strength   RecommendationStrength `json:"strength"`
	Sentiment  float64                `json:"sentiment"` // -1.0 to 1.0
	KeyQuotes  []string               `json:"keyQuotes,omitempty"`
	Summary    *string                `json:"summary,omitempty"`
	Confidence float64                `json:"confidence"`
}

// ExtractionMetadata contains information about the extraction process.
type ExtractionMetadata struct { //nolint:govet // Field ordering prioritizes JSON serialization over memory alignment
	ExtractedAt       time.Time `json:"extractedAt"`
	ModelVersion      string    `json:"modelVersion"`
	OverallConfidence float64   `json:"overallConfidence"`
	ProcessingTimeMs  *int      `json:"processingTimeMs,omitempty"`
	WarningsCount     int       `json:"warningsCount"`
	Warnings          []string  `json:"warnings,omitempty"`
}

// ExtractedLetterData is the complete extracted data from a reference letter.
type ExtractedLetterData struct {
	Author          ExtractedAuthor           `json:"author"`
	Skills          []ExtractedSkill          `json:"skills"`
	Qualities       []ExtractedQuality        `json:"qualities"`
	Accomplishments []ExtractedAccomplishment `json:"accomplishments"`
	Recommendation  ExtractedRecommendation   `json:"recommendation"`
	Metadata        ExtractionMetadata        `json:"metadata"`
}
