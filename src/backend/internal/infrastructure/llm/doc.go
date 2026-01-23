// Package llm provides LLM provider implementations and fault tolerance patterns.
//
// This package implements the domain.LLMProvider interface using the official
// anthropic-sdk-go and failsafe-go libraries for:
//   - Anthropic Claude API (with vision support for document extraction)
//   - Circuit breaker pattern for fault tolerance
//   - Retry with exponential backoff
//   - Request timeout handling
//   - Request/response logging
//
// # Usage
//
// Create a resilient Anthropic provider with all fault tolerance features:
//
//	provider := llm.NewAnthropicProvider(llm.AnthropicConfig{
//		APIKey: os.Getenv("ANTHROPIC_API_KEY"),
//	})
//
//	resilient := llm.NewResilientProvider(provider, llm.ResilientConfig{
//		MaxAttempts:      3,
//		BaseDelay:        500 * time.Millisecond,
//		MaxDelay:         30 * time.Second,
//		FailureThreshold: 5,
//		ResetTimeout:     60 * time.Second,
//		RequestTimeout:   120 * time.Second,
//	})
//
//	logged := llm.NewLoggingProvider(resilient, logger)
//
// # Structured Outputs
//
// Request structured JSON output using the OutputSchema field:
//
//	resp, err := provider.Complete(ctx, domain.LLMRequest{
//		Messages: []domain.Message{
//			domain.NewTextMessage(domain.RoleUser, "Extract the key points"),
//		},
//		OutputSchema: map[string]any{
//			"type": "object",
//			"properties": map[string]any{
//				"points": map[string]any{
//					"type": "array",
//					"items": map[string]any{"type": "string"},
//				},
//			},
//			"required": []string{"points"},
//		},
//	})
//
// # Document Extraction
//
// Extract text from documents using Claude's vision capabilities:
//
//	extractor := llm.NewDocumentExtractor(provider, llm.DocumentExtractorConfig{})
//
//	result, err := extractor.ExtractText(ctx, llm.ExtractionRequest{
//		Document:  imageBytes,
//		MediaType: domain.ImageMediaTypeJPEG,
//	})
package llm
