//nolint:errcheck,revive // Test file - error checks and unused params are OK in test helpers
package llm_test

import (
	"context"
	"testing"

	"backend/internal/domain"
	"backend/internal/infrastructure/llm"
)

func TestDocumentExtractor_ExtractTextWithRequest(t *testing.T) {
	// Mock provider that returns extracted text
	inner := &mockProvider{
		response: &domain.LLMResponse{
			Content:      "This is the extracted text from the document. It contains important information about the candidate.",
			Model:        "claude-sonnet-4-20250514",
			InputTokens:  500,
			OutputTokens: 50,
			StopReason:   "end_turn",
		},
	}

	extractor := llm.NewDocumentExtractor(inner, llm.DocumentExtractorConfig{})

	// Fake JPEG image data
	imageData := []byte{0xFF, 0xD8, 0xFF, 0xE0} // JPEG magic bytes

	result, err := extractor.ExtractTextWithRequest(context.Background(), llm.ExtractionRequest{
		Document:  imageData,
		MediaType: domain.ImageMediaTypeJPEG,
	})

	if err != nil {
		t.Fatalf("ExtractTextWithRequest() error = %v", err)
	}

	if result.Text != inner.response.Content {
		t.Errorf("Text = %q, want %q", result.Text, inner.response.Content)
	}
	if result.InputTokens != 500 {
		t.Errorf("InputTokens = %d, want 500", result.InputTokens)
	}
	if result.OutputTokens != 50 {
		t.Errorf("OutputTokens = %d, want 50", result.OutputTokens)
	}
}

func TestDocumentExtractor_ExtractTextWithRequest_PDF(t *testing.T) {
	inner := &mockProvider{
		response: &domain.LLMResponse{
			Content:      "PDF content extracted",
			Model:        "claude-sonnet-4-20250514",
			InputTokens:  1000,
			OutputTokens: 100,
		},
	}

	extractor := llm.NewDocumentExtractor(inner, llm.DocumentExtractorConfig{})

	// Fake PDF data
	pdfData := []byte("%PDF-1.4") // PDF magic bytes

	result, err := extractor.ExtractTextWithRequest(context.Background(), llm.ExtractionRequest{
		Document:  pdfData,
		MediaType: domain.ImageMediaTypePDF,
	})

	if err != nil {
		t.Fatalf("ExtractTextWithRequest() error = %v", err)
	}

	if result.Text != "PDF content extracted" {
		t.Errorf("Text = %q, want %q", result.Text, "PDF content extracted")
	}
}

func TestDocumentExtractor_ExtractTextWithRequest_CustomPrompt(t *testing.T) {
	var capturedReq domain.LLMRequest
	inner := &capturingProvider{
		response: &domain.LLMResponse{
			Content: "Extracted with custom prompt",
		},
		captureReq: &capturedReq,
	}

	extractor := llm.NewDocumentExtractor(inner, llm.DocumentExtractorConfig{})

	customPrompt := "Extract the author name and date from this letter."
	_, err := extractor.ExtractTextWithRequest(context.Background(), llm.ExtractionRequest{
		Document:     []byte{0xFF, 0xD8, 0xFF},
		MediaType:    domain.ImageMediaTypeJPEG,
		CustomPrompt: customPrompt,
	})

	if err != nil {
		t.Fatalf("ExtractTextWithRequest() error = %v", err)
	}

	// Verify custom prompt was used
	if len(capturedReq.Messages) == 0 {
		t.Fatal("expected messages in request")
	}
	msg := capturedReq.Messages[0]
	found := false
	for _, block := range msg.Content {
		if block.Type == domain.ContentTypeText && block.Text == customPrompt {
			found = true
			break
		}
	}
	if !found {
		t.Error("custom prompt not found in request")
	}
}

func TestDocumentExtractor_ExtractTextWithRequest_CustomModel(t *testing.T) {
	var capturedReq domain.LLMRequest
	inner := &capturingProvider{
		response: &domain.LLMResponse{
			Content: "Result",
		},
		captureReq: &capturedReq,
	}

	extractor := llm.NewDocumentExtractor(inner, llm.DocumentExtractorConfig{})

	_, err := extractor.ExtractTextWithRequest(context.Background(), llm.ExtractionRequest{
		Document:  []byte{0xFF, 0xD8, 0xFF},
		MediaType: domain.ImageMediaTypeJPEG,
		Model:     "claude-3-haiku-20240307",
	})

	if err != nil {
		t.Fatalf("ExtractTextWithRequest() error = %v", err)
	}

	if capturedReq.Model != "claude-3-haiku-20240307" {
		t.Errorf("Model = %q, want %q", capturedReq.Model, "claude-3-haiku-20240307")
	}
}

func TestDocumentExtractor_ExtractTextWithRequest_Error(t *testing.T) {
	inner := &mockProvider{
		err: &domain.LLMError{
			Provider: "mock",
			Message:  "API error",
		},
	}

	extractor := llm.NewDocumentExtractor(inner, llm.DocumentExtractorConfig{})

	_, err := extractor.ExtractTextWithRequest(context.Background(), llm.ExtractionRequest{
		Document:  []byte{0xFF, 0xD8, 0xFF},
		MediaType: domain.ImageMediaTypeJPEG,
	})

	if err == nil {
		t.Fatal("expected error")
	}
}

// capturingProvider captures the request for inspection.
type capturingProvider struct {
	response   *domain.LLMResponse
	captureReq *domain.LLMRequest
}

func (p *capturingProvider) Complete(ctx context.Context, req domain.LLMRequest) (*domain.LLMResponse, error) {
	*p.captureReq = req
	return p.response, nil
}

func (p *capturingProvider) Name() string {
	return "capturing"
}

//nolint:gocyclo // Test function with many assertions is intentionally thorough
func TestDocumentExtractor_ExtractLetterData(t *testing.T) {
	// Mock provider that returns structured letter extraction JSON
	jsonResponse := `{
		"author": {
			"name": "John Smith",
			"title": "Engineering Manager",
			"company": "Acme Corp",
			"relationship": "manager"
		},
		"testimonials": [
			{
				"quote": "Jane's leadership during our cloud migration was exceptional.",
				"skillsMentioned": ["leadership", "cloud architecture"]
			},
			{
				"quote": "She consistently delivered high-quality solutions under tight deadlines.",
				"skillsMentioned": ["problem solving", "time management"]
			}
		],
		"skillMentions": [
			{
				"skill": "Go",
				"quote": "Her expertise in Go helped us build a highly performant backend.",
				"context": "technical skills"
			},
			{
				"skill": "Kubernetes",
				"quote": "She led the migration to Kubernetes, reducing our deployment time by 80%.",
				"context": "infrastructure"
			}
		],
		"experienceMentions": [
			{
				"company": "Acme Corp",
				"role": "Senior Engineer",
				"quote": "During her time as Senior Engineer at Acme Corp, Jane..."
			}
		],
		"discoveredSkills": [
			{
				"skill": "mentoring",
				"quote": "Jane mentored several junior developers on the team.",
				"context": "leadership"
			},
			{
				"skill": "system design",
				"quote": "She designed the architecture for our new microservices platform.",
				"context": "technical skills"
			},
			{
				"skill": "cross-team collaboration",
				"quote": "Jane excelled at working across teams to deliver complex projects.",
				"context": "soft skills"
			}
		]
	}`

	inner := &mockProvider{
		response: &domain.LLMResponse{
			Content:      jsonResponse,
			Model:        "claude-sonnet-4-20250514",
			InputTokens:  2000,
			OutputTokens: 500,
			StopReason:   "end_turn",
		},
	}

	extractor := llm.NewDocumentExtractor(inner, llm.DocumentExtractorConfig{})

	result, err := extractor.ExtractLetterData(context.Background(), "Reference letter text here...", nil)

	if err != nil {
		t.Fatalf("ExtractLetterData() error = %v", err)
	}

	// Verify author
	if result.Author.Name != "John Smith" {
		t.Errorf("Author.Name = %q, want %q", result.Author.Name, "John Smith")
	}
	if result.Author.Title == nil || *result.Author.Title != "Engineering Manager" {
		t.Errorf("Author.Title = %v, want %q", result.Author.Title, "Engineering Manager")
	}
	if result.Author.Company == nil || *result.Author.Company != "Acme Corp" {
		t.Errorf("Author.Company = %v, want %q", result.Author.Company, "Acme Corp")
	}
	if result.Author.Relationship != domain.AuthorRelationshipManager {
		t.Errorf("Author.Relationship = %q, want %q", result.Author.Relationship, domain.AuthorRelationshipManager)
	}

	// Verify testimonials
	if len(result.Testimonials) != 2 {
		t.Errorf("len(Testimonials) = %d, want 2", len(result.Testimonials))
	}
	if result.Testimonials[0].Quote != "Jane's leadership during our cloud migration was exceptional." {
		t.Errorf("Testimonials[0].Quote unexpected: %q", result.Testimonials[0].Quote)
	}
	if len(result.Testimonials[0].SkillsMentioned) != 2 {
		t.Errorf("len(Testimonials[0].SkillsMentioned) = %d, want 2", len(result.Testimonials[0].SkillsMentioned))
	}

	// Verify skill mentions
	if len(result.SkillMentions) != 2 {
		t.Errorf("len(SkillMentions) = %d, want 2", len(result.SkillMentions))
	}
	if result.SkillMentions[0].Skill != "Go" {
		t.Errorf("SkillMentions[0].Skill = %q, want %q", result.SkillMentions[0].Skill, "Go")
	}
	if result.SkillMentions[0].Context == nil || *result.SkillMentions[0].Context != "technical skills" {
		t.Errorf("SkillMentions[0].Context = %v, want %q", result.SkillMentions[0].Context, "technical skills")
	}

	// Verify experience mentions
	if len(result.ExperienceMentions) != 1 {
		t.Errorf("len(ExperienceMentions) = %d, want 1", len(result.ExperienceMentions))
	}
	if result.ExperienceMentions[0].Company != "Acme Corp" {
		t.Errorf("ExperienceMentions[0].Company = %q, want %q", result.ExperienceMentions[0].Company, "Acme Corp")
	}
	if result.ExperienceMentions[0].Role != "Senior Engineer" {
		t.Errorf("ExperienceMentions[0].Role = %q, want %q", result.ExperienceMentions[0].Role, "Senior Engineer")
	}

	// Verify discovered skills (now objects with skill, quote, context)
	if len(result.DiscoveredSkills) != 3 {
		t.Errorf("len(DiscoveredSkills) = %d, want 3", len(result.DiscoveredSkills))
	}
	if result.DiscoveredSkills[0].Skill != "mentoring" {
		t.Errorf("DiscoveredSkills[0].Skill = %q, want %q", result.DiscoveredSkills[0].Skill, "mentoring")
	}
	if result.DiscoveredSkills[0].Quote != "Jane mentored several junior developers on the team." {
		t.Errorf("DiscoveredSkills[0].Quote = %q, want %q", result.DiscoveredSkills[0].Quote, "Jane mentored several junior developers on the team.")
	}
	if result.DiscoveredSkills[0].Context == nil || *result.DiscoveredSkills[0].Context != "leadership" {
		t.Errorf("DiscoveredSkills[0].Context = %v, want %q", result.DiscoveredSkills[0].Context, "leadership")
	}
}

func TestDocumentExtractor_ExtractLetterData_MarkdownCodeBlock(t *testing.T) {
	// Test that markdown code blocks are stripped
	jsonResponse := "```json\n" + `{
		"author": {
			"name": "Jane Doe",
			"title": "",
			"company": "",
			"relationship": "peer"
		},
		"testimonials": [],
		"skillMentions": [],
		"experienceMentions": [],
		"discoveredSkills": []
	}` + "\n```"

	inner := &mockProvider{
		response: &domain.LLMResponse{
			Content: jsonResponse,
		},
	}

	extractor := llm.NewDocumentExtractor(inner, llm.DocumentExtractorConfig{})

	result, err := extractor.ExtractLetterData(context.Background(), "Letter text", nil)

	if err != nil {
		t.Fatalf("ExtractLetterData() error = %v", err)
	}

	if result.Author.Name != "Jane Doe" {
		t.Errorf("Author.Name = %q, want %q", result.Author.Name, "Jane Doe")
	}
	if result.Author.Relationship != domain.AuthorRelationshipPeer {
		t.Errorf("Author.Relationship = %q, want %q", result.Author.Relationship, domain.AuthorRelationshipPeer)
	}
	// Empty strings should result in nil pointers
	if result.Author.Title != nil {
		t.Errorf("Author.Title = %v, want nil", result.Author.Title)
	}
	if result.Author.Company != nil {
		t.Errorf("Author.Company = %v, want nil", result.Author.Company)
	}
}

func TestDocumentExtractor_ExtractLetterData_EmptyArrays(t *testing.T) {
	// Test that empty arrays are properly initialized (not nil)
	jsonResponse := `{
		"author": {
			"name": "Bob",
			"title": "",
			"company": "",
			"relationship": "colleague"
		},
		"testimonials": [],
		"skillMentions": [],
		"experienceMentions": [],
		"discoveredSkills": []
	}`

	inner := &mockProvider{
		response: &domain.LLMResponse{
			Content: jsonResponse,
		},
	}

	extractor := llm.NewDocumentExtractor(inner, llm.DocumentExtractorConfig{})

	result, err := extractor.ExtractLetterData(context.Background(), "Letter text", nil)

	if err != nil {
		t.Fatalf("ExtractLetterData() error = %v", err)
	}

	// Arrays should be initialized, not nil
	if result.Testimonials == nil {
		t.Error("Testimonials should not be nil")
	}
	if result.SkillMentions == nil {
		t.Error("SkillMentions should not be nil")
	}
	if result.ExperienceMentions == nil {
		t.Error("ExperienceMentions should not be nil")
	}
	if result.DiscoveredSkills == nil {
		t.Error("DiscoveredSkills should not be nil")
	}
}

func TestDocumentExtractor_ExtractLetterData_AllRelationshipTypes(t *testing.T) {
	relationships := []domain.AuthorRelationship{
		domain.AuthorRelationshipManager,
		domain.AuthorRelationshipPeer,
		domain.AuthorRelationshipDirectReport,
		domain.AuthorRelationshipClient,
		domain.AuthorRelationshipMentor,
		domain.AuthorRelationshipProfessor,
		domain.AuthorRelationshipColleague,
		domain.AuthorRelationshipOther,
	}

	for _, rel := range relationships {
		t.Run(string(rel), func(t *testing.T) {
			jsonResponse := `{
				"author": {
					"name": "Test Author",
					"title": "",
					"company": "",
					"relationship": "` + string(rel) + `"
				},
				"testimonials": [],
				"skillMentions": [],
				"experienceMentions": [],
				"discoveredSkills": []
			}`

			inner := &mockProvider{
				response: &domain.LLMResponse{
					Content: jsonResponse,
				},
			}

			extractor := llm.NewDocumentExtractor(inner, llm.DocumentExtractorConfig{})

			result, err := extractor.ExtractLetterData(context.Background(), "Letter text", nil)

			if err != nil {
				t.Fatalf("ExtractLetterData() error = %v", err)
			}

			if result.Author.Relationship != rel {
				t.Errorf("Author.Relationship = %q, want %q", result.Author.Relationship, rel)
			}
		})
	}
}

func TestDocumentExtractor_ExtractLetterData_Error(t *testing.T) {
	inner := &mockProvider{
		err: &domain.LLMError{
			Provider: "mock",
			Message:  "API error",
		},
	}

	extractor := llm.NewDocumentExtractor(inner, llm.DocumentExtractorConfig{})

	_, err := extractor.ExtractLetterData(context.Background(), "Letter text", nil)

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestDocumentExtractor_ExtractLetterData_InvalidJSON(t *testing.T) {
	inner := &mockProvider{
		response: &domain.LLMResponse{
			Content: "not valid json",
		},
	}

	extractor := llm.NewDocumentExtractor(inner, llm.DocumentExtractorConfig{})

	_, err := extractor.ExtractLetterData(context.Background(), "Letter text", nil)

	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
