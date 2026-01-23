package llm

import (
	"context"

	"backend/internal/domain"
)

// FallbackProvider wraps multiple providers and falls back to secondary providers
// when the primary fails. This is a stub implementation for future expansion.
type FallbackProvider struct {
	primary   domain.LLMProvider
	fallbacks []domain.LLMProvider
}

// NewFallbackProvider creates a fallback provider with a primary and optional fallbacks.
// Currently, this is a stub that only uses the primary provider.
// TODO: Implement actual fallback logic when multiple providers are needed.
func NewFallbackProvider(primary domain.LLMProvider, fallbacks ...domain.LLMProvider) *FallbackProvider {
	return &FallbackProvider{
		primary:   primary,
		fallbacks: fallbacks,
	}
}

// Complete tries the primary provider. Fallback logic is not yet implemented.
func (p *FallbackProvider) Complete(ctx context.Context, req domain.LLMRequest) (*domain.LLMResponse, error) {
	// For now, just use primary provider.
	// Future implementation will try fallbacks on specific error conditions.
	return p.primary.Complete(ctx, req)
}

// Name returns the primary provider's name.
func (p *FallbackProvider) Name() string {
	return p.primary.Name()
}

// Verify FallbackProvider implements domain.LLMProvider.
var _ domain.LLMProvider = (*FallbackProvider)(nil)
