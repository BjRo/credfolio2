// Package model contains GraphQL model types.
package model

import "time"

// ExtractedLetterData is the GraphQL model for extracted letter data.
type ExtractedLetterData struct { //nolint:govet // Field ordering prioritizes readability over memory alignment
	Author          *ExtractedAuthor           `json:"author"`
	Skills          []string                   `json:"skills"`
	Qualities       []*ExtractedQuality        `json:"qualities"`
	Accomplishments []*ExtractedAccomplishment `json:"accomplishments"`
	Recommendation  *ExtractedRecommendation   `json:"recommendation"`
	Metadata        *ExtractionMetadata        `json:"metadata"`
}

// ExtractedAuthor is the GraphQL model for author details.
type ExtractedAuthor struct { //nolint:govet // Field ordering prioritizes readability over memory alignment
	Name                string  `json:"name"`
	Title               *string `json:"title,omitempty"`
	Organization        *string `json:"organization,omitempty"`
	Relationship        string  `json:"relationship"`
	RelationshipDetails *string `json:"relationshipDetails,omitempty"`
	Confidence          float64 `json:"confidence"`
}

// ExtractedQuality is the GraphQL model for an extracted quality.
type ExtractedQuality struct {
	Trait      string   `json:"trait"`
	Evidence   []string `json:"evidence,omitempty"`
	Confidence float64  `json:"confidence"`
}

// ExtractedAccomplishment is the GraphQL model for an accomplishment.
type ExtractedAccomplishment struct { //nolint:govet // Field ordering prioritizes readability over memory alignment
	Description string  `json:"description"`
	Impact      *string `json:"impact,omitempty"`
	Timeframe   *string `json:"timeframe,omitempty"`
	Confidence  float64 `json:"confidence"`
}

// ExtractedRecommendation is the GraphQL model for recommendation assessment.
type ExtractedRecommendation struct { //nolint:govet // Field ordering prioritizes readability over memory alignment
	Strength   string   `json:"strength"`
	Sentiment  float64  `json:"sentiment"`
	KeyQuotes  []string `json:"keyQuotes,omitempty"`
	Summary    *string  `json:"summary,omitempty"`
	Confidence float64  `json:"confidence"`
}

// ExtractionMetadata is the GraphQL model for extraction metadata.
type ExtractionMetadata struct { //nolint:govet // Field ordering prioritizes readability over memory alignment
	ExtractedAt       time.Time `json:"extractedAt"`
	ModelVersion      string    `json:"modelVersion"`
	OverallConfidence float64   `json:"overallConfidence"`
	ProcessingTimeMs  *int      `json:"processingTimeMs,omitempty"`
	WarningsCount     int       `json:"warningsCount"`
	Warnings          []string  `json:"warnings,omitempty"`
}
