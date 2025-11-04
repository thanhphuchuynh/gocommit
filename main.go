// Package main provides the entry point for gocommit, an AI-powered Git commit message generator.
//
// Architecture:
// - config/: Configuration management and CLI command handlers
// - ai/: AI provider abstraction (Gemini, future: OpenRouter)
// - git/: Git operations (diff, commit, history)
// - ui/: Terminal UI components (message selector, editor, time selector)
// - commit/: Commit workflow orchestration and delayed commit logic
// - internal/: Internal utilities (time parsing, validation)
// - logger/: Request logging functionality
//
// This main.go serves as a thin orchestration layer that:
// 1. Parses command-line flags
// 2. Handles configuration commands (--config, --enable-*, --disable-*)
// 3. Sets up the AI provider
// 4. Creates and executes the commit workflow
package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/thanhphuchuynh/ai"
	"github.com/thanhphuchuynh/commit"
	"github.com/thanhphuchuynh/config"
)

func main() {
	// Parse command-line flags
	configFlag := flag.Bool("config", false, "Configure AI provider and API key interactively")
	detailedFlag := flag.Bool("d", false, "Generate detailed commit messages with body text")
	enableLoggingFlag := flag.Bool("enable-logging", false, "Enable request logging to file")
	disableLoggingFlag := flag.Bool("disable-logging", false, "Disable request logging to file")
	iconFlag := flag.Bool("icon", false, "Use emoji icons in commit messages")
	configDelayedFlag := flag.Bool("config-delayed", false, "Configure delayed commit settings")
	enableDelayedFlag := flag.Bool("enable-delayed", false, "Enable delayed commit feature")
	disableDelayedFlag := flag.Bool("disable-delayed", false, "Disable delayed commit feature")
	setModelFlag := flag.String("set-model", "", "Set AI model for OpenRouter (e.g., anthropic/claude-3.5-sonnet)")
	autoFlag := flag.Bool("auto", false, "Auto-commit without interactive prompt (also -y)")
	yFlag := flag.Bool("y", false, "Auto-commit without interactive prompt (alias for --auto)")
	flag.Parse()

	// Handle configuration commands
	if *configFlag {
		if err := config.HandleConfigureAPIKey(); err != nil {
			log.Fatalf("Error: %v", err)
		}
		return
	}

	if *configDelayedFlag {
		if err := config.HandleConfigureDelayedCommit(); err != nil {
			log.Fatalf("Error: %v", err)
		}
		return
	}

	if *enableLoggingFlag {
		if err := config.SetLoggingEnabled(true); err != nil {
			log.Fatalf("Failed to enable logging: %v", err)
		}
		fmt.Println("Logging enabled successfully!")
		return
	}

	if *disableLoggingFlag {
		if err := config.SetLoggingEnabled(false); err != nil {
			log.Fatalf("Failed to disable logging: %v", err)
		}
		fmt.Println("Logging disabled successfully!")
		return
	}

	if *enableDelayedFlag {
		if err := config.SetDelayedCommitEnabled(true); err != nil {
			log.Fatalf("Failed to enable delayed commit: %v", err)
		}
		fmt.Println("Delayed commit feature enabled successfully!")
		return
	}

	if *disableDelayedFlag {
		if err := config.SetDelayedCommitEnabled(false); err != nil {
			log.Fatalf("Failed to disable delayed commit: %v", err)
		}
		fmt.Println("Delayed commit feature disabled successfully!")
		return
	}

	if *setModelFlag != "" {
		model := *setModelFlag
		if err := config.SetModel(model); err != nil {
			log.Fatalf("Failed to set model: %v", err)
		}
		fmt.Printf("AI model set to '%s' successfully!\n", model)
		return
	}

	// Get provider configuration (defaults to "gemini")
	providerName, err := config.GetProvider()
	if err != nil {
		log.Fatalf("Error getting provider: %v", err)
	}

	// Load API key (not needed for Ollama)
	var apiKey string
	if providerName != "ollama" {
		apiKey, err = config.GetAPIKey()
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
	}

	// Get model configuration (optional, used by OpenRouter and Ollama)
	model, err := config.GetModel()
	if err != nil {
		log.Fatalf("Error getting model: %v", err)
	}

	// Get endpoint configuration (optional, used by Ollama)
	endpoint, err := config.GetEndpoint()
	if err != nil {
		log.Fatalf("Error getting endpoint: %v", err)
	}

	// Create AI provider
	provider, err := ai.NewProvider(ai.Config{
		Provider: providerName,
		APIKey:   apiKey,
		Model:    model,
		Endpoint: endpoint,
	})
	if err != nil {
		log.Fatalf("Failed to create AI provider: %v", err)
	}

	// Create and execute commit workflow
	workflow := commit.NewWorkflow(commit.WorkflowConfig{
		Detailed: *detailedFlag,
		UseIcons: *iconFlag,
		Auto:     *autoFlag || *yFlag, // Enable auto mode if either flag is set
	}, provider)

	if err := workflow.Execute(context.Background()); err != nil {
		log.Fatalf("Commit workflow failed: %v", err)
	}
}
