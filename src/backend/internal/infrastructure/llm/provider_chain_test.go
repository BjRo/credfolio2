//nolint:errcheck,revive,goconst // Test file - error checks, unused params, and string constants are OK
package llm_test

import (
	"context"
	"testing"

	"backend/internal/domain"
	"backend/internal/infrastructure/llm"
)

func TestProviderChain_Validate(t *testing.T) {
	tests := []struct {
		name    string
		chain   llm.ProviderChain
		wantErr bool
	}{
		{
			name:    "empty chain",
			chain:   llm.ProviderChain{},
			wantErr: true,
		},
		{
			name:    "valid single entry",
			chain:   llm.ProviderChain{{Provider: "openai", Model: "gpt-4o"}},
			wantErr: false,
		},
		{
			name: "valid multi entry",
			chain: llm.ProviderChain{
				{Provider: "openai", Model: "gpt-4o"},
				{Provider: "anthropic", Model: "claude-sonnet-4-20250514"},
			},
			wantErr: false,
		},
		{
			name:    "entry with empty provider",
			chain:   llm.ProviderChain{{Provider: "", Model: "gpt-4o"}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.chain.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProviderChain_Primary(t *testing.T) {
	chain := llm.ProviderChain{
		{Provider: "openai", Model: "gpt-4o"},
		{Provider: "anthropic"},
	}
	primary := chain.Primary()
	if primary.Provider != "openai" || primary.Model != "gpt-4o" {
		t.Errorf("Primary() = %+v, want {openai, gpt-4o}", primary)
	}

	// Empty chain
	empty := llm.ProviderChain{}
	p := empty.Primary()
	if p.Provider != "" {
		t.Errorf("Primary() on empty chain = %+v, want zero value", p)
	}
}

func TestNewChainedProvider_UnregisteredProvider(t *testing.T) {
	registry := llm.NewProviderRegistry()
	registry.Register("anthropic", &mockProvider{})

	// Chain references "openai" which is NOT registered
	chain := llm.ProviderChain{{Provider: "openai", Model: "gpt-4o"}}

	_, err := llm.NewChainedProvider(registry, chain)
	if err == nil {
		t.Fatal("expected error for unregistered provider, got nil")
	}
}

func TestNewChainedProvider_RegisteredProvider(t *testing.T) {
	registry := llm.NewProviderRegistry()
	registry.Register("openai", &mockProvider{
		response: &domain.LLMResponse{Content: "test"},
	})

	chain := llm.ProviderChain{{Provider: "openai", Model: "gpt-4o"}}

	cp, err := llm.NewChainedProvider(registry, chain)
	if err != nil {
		t.Fatalf("NewChainedProvider() error = %v", err)
	}
	if cp == nil {
		t.Fatal("expected non-nil ChainedProvider")
	}
}

func TestChainedProvider_Complete_OverridesModel(t *testing.T) {
	var capturedReq domain.LLMRequest
	registry := llm.NewProviderRegistry()
	registry.Register("openai", &capturingProvider{
		response:   &domain.LLMResponse{Content: "result"},
		captureReq: &capturedReq,
	})

	chain := llm.ProviderChain{{Provider: "openai", Model: "gpt-4o"}}
	cp, err := llm.NewChainedProvider(registry, chain)
	if err != nil {
		t.Fatalf("NewChainedProvider() error = %v", err)
	}

	// Request with empty model — chain should override
	_, err = cp.Complete(context.Background(), domain.LLMRequest{
		Messages:  []domain.Message{domain.NewTextMessage(domain.RoleUser, "hi")},
		MaxTokens: 100,
	})
	if err != nil {
		t.Fatalf("Complete() error = %v", err)
	}

	if capturedReq.Model != "gpt-4o" {
		t.Errorf("Model = %q, want %q", capturedReq.Model, "gpt-4o")
	}
}

func TestChainedProvider_Complete_DoesNotOverrideExplicitModel(t *testing.T) {
	var capturedReq domain.LLMRequest
	registry := llm.NewProviderRegistry()
	registry.Register("openai", &capturingProvider{
		response:   &domain.LLMResponse{Content: "result"},
		captureReq: &capturedReq,
	})

	chain := llm.ProviderChain{{Provider: "openai", Model: "gpt-4o"}}
	cp, err := llm.NewChainedProvider(registry, chain)
	if err != nil {
		t.Fatalf("NewChainedProvider() error = %v", err)
	}

	// Request with explicit model — chain should NOT override
	_, err = cp.Complete(context.Background(), domain.LLMRequest{
		Messages:  []domain.Message{domain.NewTextMessage(domain.RoleUser, "hi")},
		Model:     "gpt-4o-mini",
		MaxTokens: 100,
	})
	if err != nil {
		t.Fatalf("Complete() error = %v", err)
	}

	if capturedReq.Model != "gpt-4o-mini" {
		t.Errorf("Model = %q, want %q (explicit model should be preserved)", capturedReq.Model, "gpt-4o-mini")
	}
}

// TestGetProviderForChain_FallsBackToDefault verifies that when chain creation
// fails (e.g. unregistered provider), the extractor falls back to the default provider.
// This is the core bug scenario — the fallback should happen but should be LOGGED.
func TestGetProviderForChain_FallsBackToDefault(t *testing.T) {
	defaultProvider := &capturingProvider{
		response:   &domain.LLMResponse{Content: "from-default"},
		captureReq: &domain.LLMRequest{},
	}

	registry := llm.NewProviderRegistry()
	registry.Register("anthropic", defaultProvider)
	// Note: "openai" is NOT registered

	extractor := llm.NewDocumentExtractor(defaultProvider, llm.DocumentExtractorConfig{
		ProviderRegistry: registry,
		ResumeExtractionChain: llm.ProviderChain{
			{Provider: "openai", Model: "gpt-4o"}, // unregistered!
		},
	})

	jsonResponse := `{
		"name": "Test", "email": "", "phone": "", "location": "",
		"summary": "", "experience": [], "education": [], "skills": [], "confidence": 0.9
	}`
	defaultProvider.response = &domain.LLMResponse{Content: jsonResponse}

	_, err := extractor.ExtractResumeData(context.Background(), "Resume text")
	if err != nil {
		t.Fatalf("ExtractResumeData() error = %v", err)
	}

	// The request went to the default provider (the fallback path)
	// This test documents the current behavior — the chain silently fell back
}

// TestExtractResumeData_UsesChainProvider verifies that when a chain is properly
// configured with a registered provider, the chain provider is used — NOT the default.
func TestExtractResumeData_UsesChainProvider(t *testing.T) {
	var defaultReq domain.LLMRequest
	defaultProvider := &capturingProvider{
		response:   &domain.LLMResponse{Content: "from-default"},
		captureReq: &defaultReq,
	}

	var chainReq domain.LLMRequest
	chainProvider := &capturingProvider{
		response:   &domain.LLMResponse{Content: "from-chain"},
		captureReq: &chainReq,
	}

	registry := llm.NewProviderRegistry()
	registry.Register("anthropic", defaultProvider)
	registry.Register("openai", chainProvider)

	jsonResponse := `{
		"name": "Test", "email": "", "phone": "", "location": "",
		"summary": "", "experience": [], "education": [], "skills": [], "confidence": 0.9
	}`
	chainProvider.response = &domain.LLMResponse{Content: jsonResponse}

	extractor := llm.NewDocumentExtractor(defaultProvider, llm.DocumentExtractorConfig{
		ProviderRegistry: registry,
		ResumeExtractionChain: llm.ProviderChain{
			{Provider: "openai", Model: "gpt-4o"},
		},
	})

	_, err := extractor.ExtractResumeData(context.Background(), "Resume text")
	if err != nil {
		t.Fatalf("ExtractResumeData() error = %v", err)
	}

	// The chain provider should have been used, and the model should be set by the chain
	if chainReq.Model != "gpt-4o" {
		t.Errorf("chain provider got Model = %q, want %q", chainReq.Model, "gpt-4o")
	}
}

// TestExtractResumeData_DoesNotSetDeprecatedDefaultModel verifies that
// ExtractResumeData does NOT set Model from config.DefaultModel (which is deprecated).
// The model should come from the chain, not from the deprecated field.
func TestExtractResumeData_DoesNotSetDeprecatedDefaultModel(t *testing.T) {
	var capturedReq domain.LLMRequest
	provider := &capturingProvider{
		captureReq: &capturedReq,
	}

	jsonResponse := `{
		"name": "Test", "email": "", "phone": "", "location": "",
		"summary": "", "experience": [], "education": [], "skills": [], "confidence": 0.9
	}`
	provider.response = &domain.LLMResponse{Content: jsonResponse}

	extractor := llm.NewDocumentExtractor(provider, llm.DocumentExtractorConfig{
		DefaultModel: "should-not-be-used",
	})

	_, err := extractor.ExtractResumeData(context.Background(), "Resume text")
	if err != nil {
		t.Fatalf("ExtractResumeData() error = %v", err)
	}

	// The request should NOT have Model set from DefaultModel
	// (the chain or empty model should be used instead)
	if capturedReq.Model == "should-not-be-used" {
		t.Error("ExtractResumeData still uses deprecated DefaultModel field; it should let the chain handle model selection")
	}
}

// TestExtractLetterData_DoesNotSetDeprecatedDefaultModel verifies that
// ExtractLetterData does NOT set Model from config.DefaultModel.
func TestExtractLetterData_DoesNotSetDeprecatedDefaultModel(t *testing.T) {
	var capturedReq domain.LLMRequest
	provider := &capturingProvider{
		captureReq: &capturedReq,
	}

	jsonResponse := `{
		"author": {"name": "Author", "title": "", "company": "", "relationship": "peer"},
		"testimonials": [], "skillMentions": [], "experienceMentions": [], "discoveredSkills": []
	}`
	provider.response = &domain.LLMResponse{Content: jsonResponse}

	extractor := llm.NewDocumentExtractor(provider, llm.DocumentExtractorConfig{
		DefaultModel: "should-not-be-used",
	})

	_, err := extractor.ExtractLetterData(context.Background(), "Letter text", nil)
	if err != nil {
		t.Fatalf("ExtractLetterData() error = %v", err)
	}

	if capturedReq.Model == "should-not-be-used" {
		t.Error("ExtractLetterData still uses deprecated DefaultModel field; it should let the chain handle model selection")
	}
}

// TestGetProviderForChain_LogsWarningOnFallback verifies that when chain creation fails,
// a warning is logged with details about the failed chain.
func TestGetProviderForChain_LogsWarningOnFallback(t *testing.T) {
	log := &mockLogger{}
	defaultProvider := &capturingProvider{
		captureReq: &domain.LLMRequest{},
	}

	registry := llm.NewProviderRegistry()
	registry.Register("anthropic", defaultProvider)
	// Note: "openai" is NOT registered

	jsonResponse := `{
		"name": "Test", "email": "", "phone": "", "location": "",
		"summary": "", "experience": [], "education": [], "skills": [], "confidence": 0.9
	}`
	defaultProvider.response = &domain.LLMResponse{Content: jsonResponse}

	extractor := llm.NewDocumentExtractor(defaultProvider, llm.DocumentExtractorConfig{
		ProviderRegistry: registry,
		ResumeExtractionChain: llm.ProviderChain{
			{Provider: "openai", Model: "gpt-4o"},
		},
		Logger: log,
	})

	_, err := extractor.ExtractResumeData(context.Background(), "Resume text")
	if err != nil {
		t.Fatalf("ExtractResumeData() error = %v", err)
	}

	// Verify warning was logged
	entries := log.getEntries()
	var found bool
	for _, entry := range entries {
		if entry.Message == "Provider chain creation failed, falling back to default provider" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected warning log about chain creation failure, got none")
	}
}

func TestProviderRegistry_Names(t *testing.T) {
	registry := llm.NewProviderRegistry()
	registry.Register("anthropic", &mockProvider{})
	registry.Register("openai", &mockProvider{})

	names := registry.Names()
	if len(names) != 2 {
		t.Errorf("Names() = %v, want 2 entries", names)
	}

	// Check both names are present (order is non-deterministic for maps)
	nameSet := make(map[string]bool)
	for _, n := range names {
		nameSet[n] = true
	}
	if !nameSet["anthropic"] || !nameSet["openai"] {
		t.Errorf("Names() = %v, want [anthropic, openai]", names)
	}
}

func TestChainedProvider_Name(t *testing.T) {
	registry := llm.NewProviderRegistry()
	registry.Register("openai", &mockProvider{})

	chain := llm.ProviderChain{{Provider: "openai", Model: "gpt-4o"}}
	cp, _ := llm.NewChainedProvider(registry, chain)

	if got := cp.Name(); got != "chained(openai)" {
		t.Errorf("Name() = %q, want %q", got, "chained(openai)")
	}
}
