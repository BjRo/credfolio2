package llm

import (
	"context"
	"fmt"

	"backend/internal/domain"
)

// ProviderModelConfig represents a single provider + model combination.
type ProviderModelConfig struct {
	// Provider is the provider name (e.g., "anthropic", "openai").
	Provider string

	// Model is the model ID to use. If empty, uses the provider's default.
	Model string
}

// ProviderChain represents an ordered list of provider+model combinations.
// The first entry is the primary choice; subsequent entries are fallbacks.
// At least one entry must be present.
type ProviderChain []ProviderModelConfig

// Validate ensures the chain has at least one entry.
func (c ProviderChain) Validate() error {
	if len(c) == 0 {
		return fmt.Errorf("provider chain must have at least one entry")
	}
	for i, cfg := range c {
		if cfg.Provider == "" {
			return fmt.Errorf("provider chain entry %d: provider name is required", i)
		}
	}
	return nil
}

// Primary returns the first (primary) provider+model configuration.
func (c ProviderChain) Primary() ProviderModelConfig {
	if len(c) == 0 {
		return ProviderModelConfig{}
	}
	return c[0]
}

// Fallbacks returns the fallback configurations (all entries after the first).
func (c ProviderChain) Fallbacks() []ProviderModelConfig {
	if len(c) <= 1 {
		return nil
	}
	return c[1:]
}

// ProviderRegistry holds named LLM providers for lookup.
type ProviderRegistry struct {
	providers map[string]domain.LLMProvider
}

// NewProviderRegistry creates a new provider registry.
func NewProviderRegistry() *ProviderRegistry {
	return &ProviderRegistry{
		providers: make(map[string]domain.LLMProvider),
	}
}

// Register adds a provider to the registry.
func (r *ProviderRegistry) Register(name string, provider domain.LLMProvider) {
	r.providers[name] = provider
}

// Get retrieves a provider by name.
func (r *ProviderRegistry) Get(name string) (domain.LLMProvider, bool) {
	p, ok := r.providers[name]
	return p, ok
}

// Names returns all registered provider names.
func (r *ProviderRegistry) Names() []string {
	names := make([]string, 0, len(r.providers))
	for name := range r.providers {
		names = append(names, name)
	}
	return names
}

// ChainedProvider wraps a provider registry and chain to execute requests.
// Currently uses only the primary provider; fallback logic can be added later.
type ChainedProvider struct {
	registry *ProviderRegistry
	chain    ProviderChain
}

// NewChainedProvider creates a provider that uses the given chain.
func NewChainedProvider(registry *ProviderRegistry, chain ProviderChain) (*ChainedProvider, error) {
	if err := chain.Validate(); err != nil {
		return nil, err
	}

	// Verify all providers in chain are registered
	for _, cfg := range chain {
		if _, ok := registry.Get(cfg.Provider); !ok {
			return nil, fmt.Errorf("provider %q not registered", cfg.Provider)
		}
	}

	return &ChainedProvider{
		registry: registry,
		chain:    chain,
	}, nil
}

// Name returns a descriptive name for the chained provider.
func (p *ChainedProvider) Name() string {
	if len(p.chain) == 0 {
		return "chained(empty)"
	}
	return fmt.Sprintf("chained(%s)", p.chain.Primary().Provider)
}

// Complete executes the request using the provider chain.
// Currently uses only the primary provider. Fallback logic will be added later.
func (p *ChainedProvider) Complete(ctx context.Context, req domain.LLMRequest) (*domain.LLMResponse, error) {
	primary := p.chain.Primary()
	provider, _ := p.registry.Get(primary.Provider)

	// Override model if specified in chain config
	if primary.Model != "" && req.Model == "" {
		req.Model = primary.Model
	}

	return provider.Complete(ctx, req)
}

// Verify ChainedProvider implements domain.LLMProvider.
var _ domain.LLMProvider = (*ChainedProvider)(nil)
