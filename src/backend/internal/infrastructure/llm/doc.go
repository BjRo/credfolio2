// Package llm provides LLM provider implementations and fault tolerance patterns.
//
// This package implements the domain.LLMProvider interface with support for:
//   - Anthropic Claude API (with vision support for document extraction)
//   - Circuit breaker pattern for fault tolerance
//   - Retry with exponential backoff
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
//		RetryConfig: llm.RetrierConfig{
//			MaxAttempts: 3,
//			BaseDelay:   500 * time.Millisecond,
//		},
//		CircuitBreakerConfig: llm.CircuitBreakerConfig{
//			FailureThreshold: 5,
//			ResetTimeout:     60 * time.Second,
//		},
//	})
//
//	logged := llm.NewLoggingProvider(resilient, logger)
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
