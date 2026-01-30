// Package domain contains the core business entities and repository interfaces.
package domain

import "time"

// AuthorRelationship represents the relationship type between letter author and candidate.
type AuthorRelationship string

// Author relationship constants.
const (
	AuthorRelationshipManager      AuthorRelationship = "manager"
	AuthorRelationshipPeer         AuthorRelationship = "peer"
	AuthorRelationshipDirectReport AuthorRelationship = "direct_report"
	AuthorRelationshipClient       AuthorRelationship = "client"
	AuthorRelationshipMentor       AuthorRelationship = "mentor"
	AuthorRelationshipProfessor    AuthorRelationship = "professor"
	AuthorRelationshipColleague    AuthorRelationship = "colleague"
	AuthorRelationshipOther        AuthorRelationship = "other"
)

// SkillCategory classifies skills into categories.
type SkillCategory string

// Skill category constants.
const (
	SkillCategoryTechnical SkillCategory = "technical"
	SkillCategorySoft      SkillCategory = "soft"
	SkillCategoryDomain    SkillCategory = "domain"
)

// ExtractedAuthor contains author details extracted from a reference letter.
type ExtractedAuthor struct {
	Name         string             `json:"name"`
	Title        *string            `json:"title,omitempty"`
	Company      *string            `json:"company,omitempty"`
	Relationship AuthorRelationship `json:"relationship"`
}

// ExtractedTestimonial represents a full quote suitable for display on the profile.
type ExtractedTestimonial struct {
	Quote           string   `json:"quote"`
	SkillsMentioned []string `json:"skillsMentioned,omitempty"`
}

// ExtractedSkillMention represents a specific skill mentioned in the letter with context.
type ExtractedSkillMention struct { //nolint:govet // Field ordering prioritizes JSON serialization over memory alignment
	Skill   string  `json:"skill"`
	Quote   string  `json:"quote"`
	Context *string `json:"context,omitempty"`
}

// ExtractedExperienceMention represents a reference to a role/company in the letter.
type ExtractedExperienceMention struct {
	Company string `json:"company"`
	Role    string `json:"role"`
	Quote   string `json:"quote"`
}

// ExtractionMetadata contains information about the extraction process.
type ExtractionMetadata struct { //nolint:govet // Field ordering prioritizes JSON serialization over memory alignment
	ExtractedAt      time.Time `json:"extractedAt"`
	ModelVersion     string    `json:"modelVersion"`
	ProcessingTimeMs *int      `json:"processingTimeMs,omitempty"`
}

// ExtractedLetterData is the complete extracted data from a reference letter.
// This schema is designed for the credibility system, focusing on validations.
type ExtractedLetterData struct { //nolint:govet // Field ordering prioritizes JSON serialization over memory alignment
	Author             ExtractedAuthor              `json:"author"`
	Testimonials       []ExtractedTestimonial       `json:"testimonials"`
	SkillMentions      []ExtractedSkillMention      `json:"skillMentions"`
	ExperienceMentions []ExtractedExperienceMention `json:"experienceMentions"`
	DiscoveredSkills   []string                     `json:"discoveredSkills"`
	Metadata           ExtractionMetadata           `json:"metadata"`
}
