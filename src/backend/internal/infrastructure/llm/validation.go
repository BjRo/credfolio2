package llm

import (
	"html"
	"strings"
	"unicode"
	"unicode/utf8"

	"backend/internal/domain"
)

// Field length limits to prevent database overflow and XSS/injection attacks.
const (
	maxNameLength        = 200
	maxEmailLength       = 320 // RFC 5321 max
	maxPhoneLength       = 50
	maxLocationLength    = 200
	maxSummaryLength     = 2000
	maxDescriptionLength = 5000
	maxQuoteLength       = 2000
	maxSkillNameLength   = 100
	maxCompanyLength     = 200
	maxTitleLength       = 200

	maxSkillsCount      = 50
	maxExperienceCount  = 20
	maxEducationCount   = 10
	maxTestimonialCount = 10
)

// ExtractedDataValidator validates and sanitizes extracted LLM data.
type ExtractedDataValidator struct{}

// NewExtractedDataValidator creates a new validator instance.
func NewExtractedDataValidator() *ExtractedDataValidator {
	return &ExtractedDataValidator{}
}

// ValidateResumeData validates and sanitizes resume extraction results.
func (v *ExtractedDataValidator) ValidateResumeData(data *domain.ResumeExtractedData) error {
	// Validate and sanitize name (required)
	if strings.TrimSpace(data.Name) == "" {
		return &domain.ValidationError{Field: "name", Message: "name is required", Err: domain.ErrEmptyRequired}
	}
	data.Name = sanitizeString(data.Name, maxNameLength)

	// Validate optional fields
	data.Email = sanitizeOptionalString(data.Email, maxEmailLength)
	data.Phone = sanitizeOptionalString(data.Phone, maxPhoneLength)
	data.Location = sanitizeOptionalString(data.Location, maxLocationLength)
	data.Summary = sanitizeOptionalString(data.Summary, maxSummaryLength)

	// Validate array lengths
	if len(data.Skills) > maxSkillsCount {
		data.Skills = data.Skills[:maxSkillsCount]
	}
	if len(data.Experience) > maxExperienceCount {
		data.Experience = data.Experience[:maxExperienceCount]
	}
	if len(data.Education) > maxEducationCount {
		data.Education = data.Education[:maxEducationCount]
	}

	// Sanitize each skill
	for i := range data.Skills {
		data.Skills[i] = sanitizeString(data.Skills[i], maxSkillNameLength)
	}

	// Sanitize experience entries
	for i := range data.Experience {
		exp := &data.Experience[i]
		exp.Company = sanitizeString(exp.Company, maxCompanyLength)
		exp.Title = sanitizeString(exp.Title, maxTitleLength)
		exp.Location = sanitizeOptionalString(exp.Location, maxLocationLength)
		exp.Description = sanitizeOptionalString(exp.Description, maxDescriptionLength)
	}

	// Sanitize education entries
	for i := range data.Education {
		edu := &data.Education[i]
		edu.Institution = sanitizeString(edu.Institution, maxCompanyLength)
		edu.Degree = sanitizeOptionalString(edu.Degree, maxTitleLength)
		edu.Field = sanitizeOptionalString(edu.Field, maxTitleLength)
		edu.Achievements = sanitizeOptionalString(edu.Achievements, maxDescriptionLength)
	}

	return nil
}

// ValidateLetterData validates and sanitizes reference letter extraction results.
func (v *ExtractedDataValidator) ValidateLetterData(data *domain.ExtractedLetterData) error {
	// Validate author name (required)
	// Note: We allow "unknown" to support German-style reference letters where
	// the author name may not be explicitly stated. Users can edit post-import.
	authorName := strings.TrimSpace(data.Author.Name)
	if authorName == "" {
		return &domain.ValidationError{Field: "author.name", Message: "author name is required", Err: domain.ErrEmptyRequired}
	}
	data.Author.Name = sanitizeString(data.Author.Name, maxNameLength)
	data.Author.Title = sanitizeOptionalString(data.Author.Title, maxTitleLength)
	data.Author.Company = sanitizeOptionalString(data.Author.Company, maxCompanyLength)

	// Validate array lengths
	if len(data.Testimonials) > maxTestimonialCount {
		data.Testimonials = data.Testimonials[:maxTestimonialCount]
	}
	if len(data.SkillMentions) > maxSkillsCount {
		data.SkillMentions = data.SkillMentions[:maxSkillsCount]
	}
	if len(data.DiscoveredSkills) > maxSkillsCount {
		data.DiscoveredSkills = data.DiscoveredSkills[:maxSkillsCount]
	}

	// Sanitize testimonials
	for i := range data.Testimonials {
		t := &data.Testimonials[i]
		t.Quote = sanitizeString(t.Quote, maxQuoteLength)
		for j := range t.SkillsMentioned {
			t.SkillsMentioned[j] = sanitizeString(t.SkillsMentioned[j], maxSkillNameLength)
		}
	}

	// Sanitize skill mentions
	for i := range data.SkillMentions {
		s := &data.SkillMentions[i]
		s.Skill = sanitizeString(s.Skill, maxSkillNameLength)
		s.Quote = sanitizeString(s.Quote, maxQuoteLength)
		s.Context = sanitizeOptionalString(s.Context, maxDescriptionLength)
	}

	// Sanitize experience mentions
	for i := range data.ExperienceMentions {
		e := &data.ExperienceMentions[i]
		e.Company = sanitizeString(e.Company, maxCompanyLength)
		e.Role = sanitizeString(e.Role, maxTitleLength)
		e.Quote = sanitizeString(e.Quote, maxQuoteLength)
	}

	// Sanitize discovered skills
	for i := range data.DiscoveredSkills {
		s := &data.DiscoveredSkills[i]
		s.Skill = sanitizeString(s.Skill, maxSkillNameLength)
		s.Quote = sanitizeString(s.Quote, maxQuoteLength)
		s.Context = sanitizeOptionalString(s.Context, maxDescriptionLength)
	}

	return nil
}

// sanitizeString cleans and truncates a string.
func sanitizeString(s string, maxLen int) string {
	// HTML escape to prevent XSS
	s = html.EscapeString(s)

	// Remove null bytes and control characters (except newlines/tabs)
	s = strings.Map(func(r rune) rune {
		if r == '\n' || r == '\t' || r == '\r' {
			return r
		}
		if unicode.IsControl(r) || r == 0 {
			return -1
		}
		return r
	}, s)

	// Trim whitespace
	s = strings.TrimSpace(s)

	// Truncate if too long - use UTF-8 aware truncation
	if len(s) > maxLen {
		s = truncateUTF8(s, maxLen)
	}

	return s
}

// truncateUTF8 truncates a string to a maximum byte length while preserving UTF-8 character boundaries.
// It ensures the resulting string is valid UTF-8 by not splitting multi-byte characters.
func truncateUTF8(s string, maxBytes int) string {
	if maxBytes <= 0 {
		return ""
	}
	if len(s) <= maxBytes {
		return s
	}

	// Walk through the string and find the last complete character that fits
	byteCount := 0
	for i := range s {
		_, size := utf8.DecodeRuneInString(s[i:])
		if byteCount+size > maxBytes {
			// Adding this character would exceed the limit
			return s[:i]
		}
		byteCount += size
	}

	// If we got here, the whole string fits
	return s
}

// sanitizeOptionalString sanitizes an optional string pointer.
func sanitizeOptionalString(sp *string, maxLen int) *string {
	if sp == nil || *sp == "" {
		return nil
	}
	sanitized := sanitizeString(*sp, maxLen)
	if sanitized == "" {
		return nil
	}
	return &sanitized
}

// Ensure ExtractedDataValidator implements domain.ExtractedDataValidator
var _ domain.ExtractedDataValidator = (*ExtractedDataValidator)(nil)
