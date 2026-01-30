// Package model contains GraphQL model types.
package model

import (
	"time"

	"backend/internal/domain"
)

// ExtractedLetterData is the GraphQL model for extracted letter data (credibility-focused).
type ExtractedLetterData struct { //nolint:govet // Field ordering prioritizes JSON serialization over memory alignment
	Author             *ExtractedAuthor              `json:"author"`
	Testimonials       []*ExtractedTestimonial       `json:"testimonials"`
	SkillMentions      []*ExtractedSkillMention      `json:"skillMentions"`
	ExperienceMentions []*ExtractedExperienceMention `json:"experienceMentions"`
	DiscoveredSkills   []string                      `json:"discoveredSkills"`
	Metadata           *ExtractionMetadata           `json:"metadata"`
}

// ExtractedAuthor is the GraphQL model for author details.
type ExtractedAuthor struct {
	Name         string                    `json:"name"`
	Title        *string                   `json:"title,omitempty"`
	Company      *string                   `json:"company,omitempty"`
	Relationship domain.AuthorRelationship `json:"relationship"`
}

// ExtractedTestimonial is the GraphQL model for a testimonial quote.
type ExtractedTestimonial struct {
	Quote           string   `json:"quote"`
	SkillsMentioned []string `json:"skillsMentioned,omitempty"`
}

// ExtractedSkillMention is the GraphQL model for a skill mention with context.
type ExtractedSkillMention struct { //nolint:govet // Field ordering prioritizes JSON serialization over memory alignment
	Skill   string  `json:"skill"`
	Quote   string  `json:"quote"`
	Context *string `json:"context,omitempty"`
}

// ExtractedExperienceMention is the GraphQL model for an experience mention.
type ExtractedExperienceMention struct {
	Company string `json:"company"`
	Role    string `json:"role"`
	Quote   string `json:"quote"`
}

// ExtractionMetadata is the GraphQL model for extraction metadata.
type ExtractionMetadata struct { //nolint:govet // Field ordering prioritizes JSON serialization over memory alignment
	ExtractedAt      time.Time `json:"extractedAt"`
	ModelVersion     string    `json:"modelVersion"`
	ProcessingTimeMs *int      `json:"processingTimeMs,omitempty"`
}
