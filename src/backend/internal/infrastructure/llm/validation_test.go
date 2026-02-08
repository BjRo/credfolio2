package llm

import (
	"errors"
	"strings"
	"testing"
	"time"

	"backend/internal/domain"
)

// TestValidateResumeData_RequiredFields tests that required fields are validated.
func TestValidateResumeData_RequiredFields(t *testing.T) {
	validator := NewExtractedDataValidator()

	t.Run("rejects empty name", func(t *testing.T) {
		data := &domain.ResumeExtractedData{
			Name:        "",
			ExtractedAt: time.Now(),
		}

		err := validator.ValidateResumeData(data)
		if err == nil {
			t.Fatal("expected validation error for empty name")
		}

		var valErr *domain.ValidationError
		if !errors.As(err, &valErr) {
			t.Fatalf("expected ValidationError, got %T", err)
		}

		if valErr.Field != "name" {
			t.Errorf("expected field 'name', got %q", valErr.Field)
		}

		if !errors.Is(err, domain.ErrEmptyRequired) {
			t.Error("expected error to wrap ErrEmptyRequired")
		}
	})

	t.Run("accepts valid name", func(t *testing.T) {
		data := &domain.ResumeExtractedData{
			Name:        "John Doe",
			ExtractedAt: time.Now(),
		}

		err := validator.ValidateResumeData(data)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if data.Name != "John Doe" {
			t.Errorf("name should not be modified, got %q", data.Name)
		}
	})

	t.Run("trims whitespace from name", func(t *testing.T) {
		data := &domain.ResumeExtractedData{
			Name:        "  John Doe  ",
			ExtractedAt: time.Now(),
		}

		err := validator.ValidateResumeData(data)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if data.Name != "John Doe" {
			t.Errorf("expected trimmed name, got %q", data.Name)
		}
	})
}

// TestValidateResumeData_FieldLengthLimits tests that fields are truncated to max lengths.
func TestValidateResumeData_FieldLengthLimits(t *testing.T) {
	validator := NewExtractedDataValidator()

	t.Run("truncates name at 200 chars", func(t *testing.T) {
		longName := strings.Repeat("a", 250)
		data := &domain.ResumeExtractedData{
			Name:        longName,
			ExtractedAt: time.Now(),
		}

		err := validator.ValidateResumeData(data)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(data.Name) != maxNameLength {
			t.Errorf("expected name length %d, got %d", maxNameLength, len(data.Name))
		}
	})

	t.Run("truncates summary at 2000 chars", func(t *testing.T) {
		longSummary := strings.Repeat("a", 2500)
		data := &domain.ResumeExtractedData{
			Name:        "John Doe",
			Summary:     &longSummary,
			ExtractedAt: time.Now(),
		}

		err := validator.ValidateResumeData(data)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if data.Summary == nil {
			t.Fatal("summary should not be nil")
		}

		if len(*data.Summary) != maxSummaryLength {
			t.Errorf("expected summary length %d, got %d", maxSummaryLength, len(*data.Summary))
		}
	})

	t.Run("truncates experience description at 5000 chars", func(t *testing.T) {
		longDesc := strings.Repeat("a", 6000)
		data := &domain.ResumeExtractedData{
			Name: "John Doe",
			Experience: []domain.WorkExperience{
				{
					Company:     "Tech Corp",
					Title:       "Engineer",
					Description: &longDesc,
				},
			},
			ExtractedAt: time.Now(),
		}

		err := validator.ValidateResumeData(data)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if data.Experience[0].Description == nil {
			t.Fatal("description should not be nil")
		}

		if len(*data.Experience[0].Description) != maxDescriptionLength {
			t.Errorf("expected description length %d, got %d", maxDescriptionLength, len(*data.Experience[0].Description))
		}
	})
}

// TestValidateResumeData_HTMLEscaping tests that HTML is escaped to prevent XSS.
func TestValidateResumeData_HTMLEscaping(t *testing.T) {
	validator := NewExtractedDataValidator()

	t.Run("escapes HTML in name", func(t *testing.T) {
		data := &domain.ResumeExtractedData{
			Name:        "<script>alert('XSS')</script>John",
			ExtractedAt: time.Now(),
		}

		err := validator.ValidateResumeData(data)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if strings.Contains(data.Name, "<script>") {
			t.Errorf("expected HTML to be escaped, got %q", data.Name)
		}

		if !strings.Contains(data.Name, "&lt;script&gt;") {
			t.Errorf("expected escaped HTML, got %q", data.Name)
		}
	})

	t.Run("escapes HTML in summary", func(t *testing.T) {
		summary := "<b>Bold</b> text with <img src=x onerror=alert(1)>"
		data := &domain.ResumeExtractedData{
			Name:        "John Doe",
			Summary:     &summary,
			ExtractedAt: time.Now(),
		}

		err := validator.ValidateResumeData(data)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if data.Summary == nil {
			t.Fatal("summary should not be nil")
		}

		if strings.Contains(*data.Summary, "<b>") || strings.Contains(*data.Summary, "<img") {
			t.Errorf("expected HTML to be escaped, got %q", *data.Summary)
		}

		if !strings.Contains(*data.Summary, "&lt;b&gt;") {
			t.Errorf("expected escaped HTML, got %q", *data.Summary)
		}
	})

	t.Run("escapes HTML in skills", func(t *testing.T) {
		data := &domain.ResumeExtractedData{
			Name:        "John Doe",
			Skills:      []string{"JavaScript", "<script>alert(1)</script>"},
			ExtractedAt: time.Now(),
		}

		err := validator.ValidateResumeData(data)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if strings.Contains(data.Skills[1], "<script>") {
			t.Errorf("expected HTML to be escaped in skills, got %q", data.Skills[1])
		}
	})
}

// TestValidateResumeData_ControlCharacters tests removal of control characters.
func TestValidateResumeData_ControlCharacters(t *testing.T) {
	validator := NewExtractedDataValidator()

	t.Run("removes null bytes", func(t *testing.T) {
		data := &domain.ResumeExtractedData{
			Name:        "John\x00Doe",
			ExtractedAt: time.Now(),
		}

		err := validator.ValidateResumeData(data)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if strings.Contains(data.Name, "\x00") {
			t.Errorf("expected null byte to be removed, got %q", data.Name)
		}

		if data.Name != "JohnDoe" {
			t.Errorf("expected 'JohnDoe', got %q", data.Name)
		}
	})

	t.Run("removes control characters except newlines and tabs", func(t *testing.T) {
		summary := "Line 1\nLine 2\tTab\x01\x02\x03End"
		data := &domain.ResumeExtractedData{
			Name:        "John Doe",
			Summary:     &summary,
			ExtractedAt: time.Now(),
		}

		err := validator.ValidateResumeData(data)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if data.Summary == nil {
			t.Fatal("summary should not be nil")
		}

		// Should keep \n and \t but remove \x01, \x02, \x03
		if !strings.Contains(*data.Summary, "\n") {
			t.Error("expected newline to be preserved")
		}
		if !strings.Contains(*data.Summary, "\t") {
			t.Error("expected tab to be preserved")
		}
		if strings.Contains(*data.Summary, "\x01") || strings.Contains(*data.Summary, "\x02") {
			t.Errorf("expected control characters to be removed, got %q", *data.Summary)
		}
	})
}

// TestValidateResumeData_ArrayLimits tests that arrays are limited in size.
func TestValidateResumeData_ArrayLimits(t *testing.T) {
	validator := NewExtractedDataValidator()

	t.Run("limits skills to 50", func(t *testing.T) {
		skills := make([]string, 60)
		for i := range skills {
			skills[i] = "Skill" + string(rune(i))
		}

		data := &domain.ResumeExtractedData{
			Name:        "John Doe",
			Skills:      skills,
			ExtractedAt: time.Now(),
		}

		err := validator.ValidateResumeData(data)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(data.Skills) != maxSkillsCount {
			t.Errorf("expected %d skills, got %d", maxSkillsCount, len(data.Skills))
		}
	})

	t.Run("limits experience to 20", func(t *testing.T) {
		experiences := make([]domain.WorkExperience, 25)
		for i := range experiences {
			experiences[i] = domain.WorkExperience{
				Company: "Company" + string(rune(i)),
				Title:   "Title" + string(rune(i)),
			}
		}

		data := &domain.ResumeExtractedData{
			Name:        "John Doe",
			Experience:  experiences,
			ExtractedAt: time.Now(),
		}

		err := validator.ValidateResumeData(data)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(data.Experience) != maxExperienceCount {
			t.Errorf("expected %d experiences, got %d", maxExperienceCount, len(data.Experience))
		}
	})

	t.Run("limits education to 10", func(t *testing.T) {
		education := make([]domain.Education, 15)
		for i := range education {
			education[i] = domain.Education{
				Institution: "School" + string(rune(i)),
			}
		}

		data := &domain.ResumeExtractedData{
			Name:        "John Doe",
			Education:   education,
			ExtractedAt: time.Now(),
		}

		err := validator.ValidateResumeData(data)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(data.Education) != maxEducationCount {
			t.Errorf("expected %d education entries, got %d", maxEducationCount, len(data.Education))
		}
	})
}

// TestValidateResumeData_OptionalFields tests handling of optional fields.
func TestValidateResumeData_OptionalFields(t *testing.T) {
	validator := NewExtractedDataValidator()

	t.Run("sets empty optional string pointers to nil", func(t *testing.T) {
		emptyStr := ""
		data := &domain.ResumeExtractedData{
			Name:        "John Doe",
			Email:       &emptyStr,
			Phone:       &emptyStr,
			Location:    &emptyStr,
			Summary:     &emptyStr,
			ExtractedAt: time.Now(),
		}

		err := validator.ValidateResumeData(data)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if data.Email != nil {
			t.Error("expected empty email to be nil")
		}
		if data.Phone != nil {
			t.Error("expected empty phone to be nil")
		}
		if data.Location != nil {
			t.Error("expected empty location to be nil")
		}
		if data.Summary != nil {
			t.Error("expected empty summary to be nil")
		}
	})

	t.Run("preserves valid optional fields", func(t *testing.T) {
		email := "john@example.com"
		phone := "+1-555-0100"
		data := &domain.ResumeExtractedData{
			Name:        "John Doe",
			Email:       &email,
			Phone:       &phone,
			ExtractedAt: time.Now(),
		}

		err := validator.ValidateResumeData(data)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if data.Email == nil || *data.Email != "john@example.com" {
			t.Error("expected email to be preserved")
		}
		if data.Phone == nil || *data.Phone != "+1-555-0100" {
			t.Error("expected phone to be preserved")
		}
	})
}

// TestValidateLetterData_RequiredFields tests required field validation for letters.
func TestValidateLetterData_RequiredFields(t *testing.T) {
	validator := NewExtractedDataValidator()

	t.Run("rejects empty author name", func(t *testing.T) {
		data := &domain.ExtractedLetterData{
			Author: domain.ExtractedAuthor{
				Name:         "",
				Relationship: domain.AuthorRelationshipManager,
			},
			Metadata: domain.ExtractionMetadata{
				ExtractedAt:  time.Now(),
				ModelVersion: "test",
			},
		}

		err := validator.ValidateLetterData(data)
		if err == nil {
			t.Fatal("expected validation error for empty author name")
		}

		var valErr *domain.ValidationError
		if !errors.As(err, &valErr) {
			t.Fatalf("expected ValidationError, got %T", err)
		}

		if valErr.Field != "author.name" {
			t.Errorf("expected field 'author.name', got %q", valErr.Field)
		}

		if !errors.Is(err, domain.ErrEmptyRequired) {
			t.Error("expected error to wrap ErrEmptyRequired")
		}
	})

	t.Run("accepts valid author name", func(t *testing.T) {
		data := &domain.ExtractedLetterData{
			Author: domain.ExtractedAuthor{
				Name:         "Jane Smith",
				Relationship: domain.AuthorRelationshipManager,
			},
			Metadata: domain.ExtractionMetadata{
				ExtractedAt:  time.Now(),
				ModelVersion: "test",
			},
		}

		err := validator.ValidateLetterData(data)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if data.Author.Name != "Jane Smith" {
			t.Errorf("author name should not be modified, got %q", data.Author.Name)
		}
	})
}

// TestValidateLetterData_Testimonials tests testimonial validation.
func TestValidateLetterData_Testimonials(t *testing.T) {
	validator := NewExtractedDataValidator()

	t.Run("truncates testimonial quote at 2000 chars", func(t *testing.T) {
		longQuote := strings.Repeat("a", 2500)
		data := &domain.ExtractedLetterData{
			Author: domain.ExtractedAuthor{
				Name:         "Jane Smith",
				Relationship: domain.AuthorRelationshipManager,
			},
			Testimonials: []domain.ExtractedTestimonial{
				{Quote: longQuote},
			},
			Metadata: domain.ExtractionMetadata{
				ExtractedAt:  time.Now(),
				ModelVersion: "test",
			},
		}

		err := validator.ValidateLetterData(data)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(data.Testimonials[0].Quote) != maxQuoteLength {
			t.Errorf("expected quote length %d, got %d", maxQuoteLength, len(data.Testimonials[0].Quote))
		}
	})

	t.Run("escapes HTML in testimonial quote", func(t *testing.T) {
		data := &domain.ExtractedLetterData{
			Author: domain.ExtractedAuthor{
				Name:         "Jane Smith",
				Relationship: domain.AuthorRelationshipManager,
			},
			Testimonials: []domain.ExtractedTestimonial{
				{Quote: "<script>alert(1)</script>Great work!"},
			},
			Metadata: domain.ExtractionMetadata{
				ExtractedAt:  time.Now(),
				ModelVersion: "test",
			},
		}

		err := validator.ValidateLetterData(data)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if strings.Contains(data.Testimonials[0].Quote, "<script>") {
			t.Errorf("expected HTML to be escaped, got %q", data.Testimonials[0].Quote)
		}

		if !strings.Contains(data.Testimonials[0].Quote, "&lt;script&gt;") {
			t.Errorf("expected escaped HTML, got %q", data.Testimonials[0].Quote)
		}
	})

	t.Run("limits testimonials to 10", func(t *testing.T) {
		testimonials := make([]domain.ExtractedTestimonial, 15)
		for i := range testimonials {
			testimonials[i] = domain.ExtractedTestimonial{
				Quote: "Quote " + string(rune(i)),
			}
		}

		data := &domain.ExtractedLetterData{
			Author: domain.ExtractedAuthor{
				Name:         "Jane Smith",
				Relationship: domain.AuthorRelationshipManager,
			},
			Testimonials: testimonials,
			Metadata: domain.ExtractionMetadata{
				ExtractedAt:  time.Now(),
				ModelVersion: "test",
			},
		}

		err := validator.ValidateLetterData(data)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(data.Testimonials) != maxTestimonialCount {
			t.Errorf("expected %d testimonials, got %d", maxTestimonialCount, len(data.Testimonials))
		}
	})
}

// TestValidateLetterData_SkillMentions tests skill mention validation.
func TestValidateLetterData_SkillMentions(t *testing.T) {
	validator := NewExtractedDataValidator()

	t.Run("sanitizes skill mentions", func(t *testing.T) {
		context := "Some <b>context</b>"
		data := &domain.ExtractedLetterData{
			Author: domain.ExtractedAuthor{
				Name:         "Jane Smith",
				Relationship: domain.AuthorRelationshipManager,
			},
			SkillMentions: []domain.ExtractedSkillMention{
				{
					Skill:   "<script>alert(1)</script>Python",
					Quote:   "Great at Python!",
					Context: &context,
				},
			},
			Metadata: domain.ExtractionMetadata{
				ExtractedAt:  time.Now(),
				ModelVersion: "test",
			},
		}

		err := validator.ValidateLetterData(data)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if strings.Contains(data.SkillMentions[0].Skill, "<script>") {
			t.Error("expected HTML to be escaped in skill")
		}
		if data.SkillMentions[0].Context == nil || strings.Contains(*data.SkillMentions[0].Context, "<b>") {
			t.Error("expected HTML to be escaped in context")
		}
	})

	t.Run("limits skill mentions to 50", func(t *testing.T) {
		mentions := make([]domain.ExtractedSkillMention, 60)
		for i := range mentions {
			mentions[i] = domain.ExtractedSkillMention{
				Skill: "Skill" + string(rune(i)),
				Quote: "Quote",
			}
		}

		data := &domain.ExtractedLetterData{
			Author: domain.ExtractedAuthor{
				Name:         "Jane Smith",
				Relationship: domain.AuthorRelationshipManager,
			},
			SkillMentions: mentions,
			Metadata: domain.ExtractionMetadata{
				ExtractedAt:  time.Now(),
				ModelVersion: "test",
			},
		}

		err := validator.ValidateLetterData(data)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(data.SkillMentions) != maxSkillsCount {
			t.Errorf("expected %d skill mentions, got %d", maxSkillsCount, len(data.SkillMentions))
		}
	})
}

// TestValidateLetterData_DiscoveredSkills tests discovered skill validation.
func TestValidateLetterData_DiscoveredSkills(t *testing.T) {
	validator := NewExtractedDataValidator()

	t.Run("sanitizes discovered skills", func(t *testing.T) {
		data := &domain.ExtractedLetterData{
			Author: domain.ExtractedAuthor{
				Name:         "Jane Smith",
				Relationship: domain.AuthorRelationshipManager,
			},
			DiscoveredSkills: []domain.DiscoveredSkill{
				{
					Skill:    "<img src=x>Leadership",
					Quote:    "Excellent leader!",
					Category: domain.SkillCategorySoft,
				},
			},
			Metadata: domain.ExtractionMetadata{
				ExtractedAt:  time.Now(),
				ModelVersion: "test",
			},
		}

		err := validator.ValidateLetterData(data)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if strings.Contains(data.DiscoveredSkills[0].Skill, "<img") {
			t.Error("expected HTML to be escaped in discovered skill")
		}
	})

	t.Run("limits discovered skills to 50", func(t *testing.T) {
		skills := make([]domain.DiscoveredSkill, 60)
		for i := range skills {
			skills[i] = domain.DiscoveredSkill{
				Skill:    "Skill" + string(rune(i)),
				Quote:    "Quote",
				Category: domain.SkillCategoryTechnical,
			}
		}

		data := &domain.ExtractedLetterData{
			Author: domain.ExtractedAuthor{
				Name:         "Jane Smith",
				Relationship: domain.AuthorRelationshipManager,
			},
			DiscoveredSkills: skills,
			Metadata: domain.ExtractionMetadata{
				ExtractedAt:  time.Now(),
				ModelVersion: "test",
			},
		}

		err := validator.ValidateLetterData(data)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(data.DiscoveredSkills) != maxSkillsCount {
			t.Errorf("expected %d discovered skills, got %d", maxSkillsCount, len(data.DiscoveredSkills))
		}
	})
}
