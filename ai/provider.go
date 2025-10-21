package ai

import (
	"context"
	"fmt"
)

// Provider defines the interface for AI commit message generators
type Provider interface {
	GenerateMessages(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error)
	Name() string
	ValidateConfig() error
}

// GenerateRequest contains all parameters for message generation
type GenerateRequest struct {
	Diff       string
	LastCommit string
	Detailed   bool
	UseIcons   bool
}

// GenerateResponse contains generated commit messages
type GenerateResponse struct {
	Messages []string
	Provider string
}

// Config holds AI provider configuration
type Config struct {
	Provider string // "gemini" or "openrouter"
	APIKey   string
	Model    string // optional, provider-specific
}

// NewProvider creates an AI provider based on configuration
func NewProvider(cfg Config) (Provider, error) {
	switch cfg.Provider {
	case "gemini", "":
		return NewGeminiProvider(cfg.APIKey), nil
	case "openrouter":
		if cfg.APIKey == "" {
			return nil, fmt.Errorf("OpenRouter API key is required")
		}
		return NewOpenRouterProvider(cfg.APIKey, cfg.Model), nil
	default:
		return nil, fmt.Errorf("unknown provider: %s", cfg.Provider)
	}
}
