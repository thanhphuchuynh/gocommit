package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	configFileName = ".gocommit.json"
)

type DelayedCommitConfig struct {
	Enabled             bool `json:"enabled"`
	RestrictedStartHour int  `json:"restricted_start_hour"`
	RestrictedEndHour   int  `json:"restricted_end_hour"`
	SuggestionInterval  int  `json:"suggestion_interval"`
}

type Config struct {
	APIKey         string              `json:"api_key"`
	Provider       string              `json:"provider"` // "gemini", "openrouter", or "ollama", default: "gemini"
	Model          string              `json:"model"`    // Model name for provider (e.g., "anthropic/claude-3.5-sonnet" for OpenRouter, "codellama:7b" for Ollama)
	Endpoint       string              `json:"endpoint"` // Ollama endpoint (e.g., "http://localhost:11434"), optional
	LoggingEnabled bool                `json:"logging_enabled"`
	IconMode       bool                `json:"icon_mode"`
	DelayedCommit  DelayedCommitConfig `json:"delayed_commit"`
}

func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %v", err)
	}
	return filepath.Join(homeDir, configFileName), nil
}

func LoadConfig() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Return config with sensible defaults
		return &Config{
			DelayedCommit: DelayedCommitConfig{
				Enabled:             false,
				RestrictedStartHour: 9,
				RestrictedEndHour:   17,
				SuggestionInterval:  20,
			},
		}, nil
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	// Apply defaults for delayed commit if not set
	if config.DelayedCommit.RestrictedStartHour == 0 && config.DelayedCommit.RestrictedEndHour == 0 {
		config.DelayedCommit.RestrictedStartHour = 9
		config.DelayedCommit.RestrictedEndHour = 17
	}
	if config.DelayedCommit.SuggestionInterval == 0 {
		config.DelayedCommit.SuggestionInterval = 20
	}

	return &config, nil
}

func SaveConfig(config *Config) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

func GetAPIKey() (string, error) {
	config, err := LoadConfig()
	if err != nil {
		return "", err
	}

	if config.APIKey == "" {
		return "", fmt.Errorf("API key not configured. Please run 'gocommit --config' to set your API key")
	}

	return config.APIKey, nil
}

func SetAPIKey(apiKey string) error {
	config, err := LoadConfig()
	if err != nil {
		return err
	}

	config.APIKey = apiKey
	return SaveConfig(config)
}

func GetProvider() (string, error) {
	config, err := LoadConfig()
	if err != nil {
		return "", err
	}

	// Default to "gemini" if not set
	if config.Provider == "" {
		return "gemini", nil
	}

	return config.Provider, nil
}

func SetProvider(provider string) error {
	config, err := LoadConfig()
	if err != nil {
		return err
	}

	config.Provider = provider
	return SaveConfig(config)
}

func GetModel() (string, error) {
	config, err := LoadConfig()
	if err != nil {
		return "", err
	}

	return config.Model, nil
}

func SetModel(model string) error {
	config, err := LoadConfig()
	if err != nil {
		return err
	}

	config.Model = model
	return SaveConfig(config)
}

func GetEndpoint() (string, error) {
	config, err := LoadConfig()
	if err != nil {
		return "", err
	}

	return config.Endpoint, nil
}

func SetEndpoint(endpoint string) error {
	config, err := LoadConfig()
	if err != nil {
		return err
	}

	config.Endpoint = endpoint
	return SaveConfig(config)
}

func IsLoggingEnabled() (bool, error) {
	config, err := LoadConfig()
	if err != nil {
		return false, err
	}

	return config.LoggingEnabled, nil
}

func SetLoggingEnabled(enabled bool) error {
	config, err := LoadConfig()
	if err != nil {
		return err
	}

	config.LoggingEnabled = enabled
	return SaveConfig(config)
}

func IsIconModeEnabled() (bool, error) {
	config, err := LoadConfig()
	if err != nil {
		return false, err
	}

	return config.IconMode, nil
}

func SetIconMode(enabled bool) error {
	config, err := LoadConfig()
	if err != nil {
		return err
	}

	config.IconMode = enabled
	return SaveConfig(config)
}

func IsDelayedCommitEnabled() (bool, error) {
	config, err := LoadConfig()
	if err != nil {
		return false, err
	}

	return config.DelayedCommit.Enabled, nil
}

func SetDelayedCommitEnabled(enabled bool) error {
	config, err := LoadConfig()
	if err != nil {
		return err
	}

	config.DelayedCommit.Enabled = enabled
	return SaveConfig(config)
}

func GetDelayedCommitConfig() (*DelayedCommitConfig, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	return &config.DelayedCommit, nil
}

func SetDelayedCommitConfig(startHour, endHour, interval int) error {
	// Validate hours
	if startHour < 0 || startHour > 23 {
		return fmt.Errorf("start hour must be between 0 and 23")
	}
	if endHour < 0 || endHour > 23 {
		return fmt.Errorf("end hour must be between 0 and 23")
	}

	// Validate interval
	if interval <= 0 {
		return fmt.Errorf("suggestion interval must be positive")
	}

	config, err := LoadConfig()
	if err != nil {
		return err
	}

	config.DelayedCommit.RestrictedStartHour = startHour
	config.DelayedCommit.RestrictedEndHour = endHour
	config.DelayedCommit.SuggestionInterval = interval
	return SaveConfig(config)
}

func GetRestrictedHours() (int, int, error) {
	config, err := LoadConfig()
	if err != nil {
		return 0, 0, err
	}

	return config.DelayedCommit.RestrictedStartHour, config.DelayedCommit.RestrictedEndHour, nil
}

func GetSuggestionInterval() (int, error) {
	config, err := LoadConfig()
	if err != nil {
		return 0, err
	}

	return config.DelayedCommit.SuggestionInterval, nil
}

// ValidateGeminiAPIKey validates the Gemini API key format
func ValidateGeminiAPIKey(apiKey string) error {
	// Google AI API keys typically start with "AIza" and are 39 characters long
	if len(apiKey) != 39 || !strings.HasPrefix(apiKey, "AIza") {
		return fmt.Errorf("invalid API key format. Google AI API keys should start with 'AIza' and be 39 characters long")
	}
	return nil
}

// ValidateOpenRouterAPIKey validates the OpenRouter API key format
func ValidateOpenRouterAPIKey(apiKey string) error {
	// OpenRouter API keys start with "sk-or-" or "sk-"
	if !strings.HasPrefix(apiKey, "sk-or-") && !strings.HasPrefix(apiKey, "sk-") {
		return fmt.Errorf("invalid API key format. OpenRouter API keys should start with 'sk-or-' or 'sk-'")
	}
	return nil
}

// HandleConfigureAPIKey handles the --config flag for configuring AI provider and API key
func HandleConfigureAPIKey() error {
	fmt.Print("\n=== GoCommit Configuration ===\n")

	// Step 1: Choose provider
	fmt.Println("Select AI Provider:")
	fmt.Println("  1. Gemini (Google)")
	fmt.Println("  2. OpenRouter (Claude, GPT-4, Llama, etc.)")
	fmt.Println("  3. Ollama (Local AI)")
	fmt.Print("\nEnter choice (1, 2, or 3) [1]: ")

	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	// Default to Gemini if no choice or invalid choice
	provider := "gemini"
	if choice == "2" {
		provider = "openrouter"
	} else if choice == "3" {
		provider = "ollama"
	}

	// Step 2: Enter API key (skip for Ollama)
	var apiKey string
	if provider != "ollama" {
		providerName := "Gemini"
		if provider == "openrouter" {
			providerName = "OpenRouter"
		}
		fmt.Printf("\nEnter %s API key: ", providerName)
		apiKey, _ = reader.ReadString('\n')
		apiKey = strings.TrimSpace(apiKey)

		if apiKey == "" {
			return fmt.Errorf("API key cannot be empty")
		}

		// Validate API key based on provider
		if provider == "gemini" {
			if err := ValidateGeminiAPIKey(apiKey); err != nil {
				return err
			}
		} else if provider == "openrouter" {
			if err := ValidateOpenRouterAPIKey(apiKey); err != nil {
				return err
			}
		}
	}

	// Step 3: Configure provider-specific options
	if provider == "openrouter" {
		fmt.Println("\nConfigure OpenRouter model (optional):")
		fmt.Println("  Popular models:")
		fmt.Println("    - anthropic/claude-3.5-sonnet (default)")
		fmt.Println("    - anthropic/claude-3-opus")
		fmt.Println("    - openai/gpt-4-turbo")
		fmt.Println("    - openai/gpt-4")
		fmt.Println("    - meta-llama/llama-3.1-70b-instruct")
		fmt.Print("\nEnter model name (or press Enter for default): ")

		model, _ := reader.ReadString('\n')
		model = strings.TrimSpace(model)

		if model != "" {
			if err := SetModel(model); err != nil {
				return fmt.Errorf("error setting model: %w", err)
			}
		}
	} else if provider == "ollama" {
		// Configure Ollama endpoint
		fmt.Println("\nConfigure Ollama:")
		fmt.Print("Enter Ollama endpoint [http://localhost:11434]: ")
		endpoint, _ := reader.ReadString('\n')
		endpoint = strings.TrimSpace(endpoint)

		if endpoint == "" {
			endpoint = "http://localhost:11434"
		}

		if err := SetEndpoint(endpoint); err != nil {
			return fmt.Errorf("error setting endpoint: %w", err)
		}

		// Configure Ollama model
		fmt.Println("\nConfigure Ollama model:")
		fmt.Println("  Recommended models:")
		fmt.Println("    - codellama:7b (default, 3.8GB)")
		fmt.Println("    - llama3:8b (4.7GB)")
		fmt.Println("    - mistral:7b (4.1GB)")
		fmt.Println("    - deepseek-coder:6.7b (3.8GB)")
		fmt.Println("    - codellama:13b (7.4GB, better quality)")
		fmt.Print("\nEnter model name [codellama:7b]: ")

		model, _ := reader.ReadString('\n')
		model = strings.TrimSpace(model)

		if model == "" {
			model = "codellama:7b"
		}

		if err := SetModel(model); err != nil {
			return fmt.Errorf("error setting model: %w", err)
		}

		fmt.Println("\nℹ️  Make sure Ollama is running and the model is pulled:")
		fmt.Printf("   ollama pull %s\n", model)
	}

	// Save configuration
	if err := SetProvider(provider); err != nil {
		return fmt.Errorf("error setting provider: %w", err)
	}

	if provider != "ollama" {
		if err := SetAPIKey(apiKey); err != nil {
			return fmt.Errorf("error saving API key: %w", err)
		}
	}

	fmt.Printf("\n✓ Configuration saved successfully!\n")
	fmt.Printf("  Provider: %s\n", provider)
	if provider == "openrouter" {
		model, _ := GetModel()
		if model != "" {
			fmt.Printf("  Model: %s\n", model)
		} else {
			fmt.Printf("  Model: anthropic/claude-3.5-sonnet (default)\n")
		}
	} else if provider == "ollama" {
		endpoint, _ := GetEndpoint()
		model, _ := GetModel()
		fmt.Printf("  Endpoint: %s\n", endpoint)
		fmt.Printf("  Model: %s\n", model)
	}
	fmt.Println("\nYou can now use 'gocommit' to generate commit messages.")

	return nil
}

// HandleConfigureDelayedCommit handles the --config-delayed flag
func HandleConfigureDelayedCommit() error {
	fmt.Println("Configure Delayed Commit Settings")
	fmt.Println("===================================")
	fmt.Println()

	// Get restricted start hour
	fmt.Print("Enter restricted start hour (0-23, e.g., 9 for 9 AM): ")
	var startHourStr string
	fmt.Scanln(&startHourStr)
	startHour, err := strconv.Atoi(strings.TrimSpace(startHourStr))
	if err != nil || startHour < 0 || startHour > 23 {
		return fmt.Errorf("invalid start hour. Must be between 0 and 23")
	}

	// Get restricted end hour
	fmt.Print("Enter restricted end hour (0-23, e.g., 17 for 5 PM): ")
	var endHourStr string
	fmt.Scanln(&endHourStr)
	endHour, err := strconv.Atoi(strings.TrimSpace(endHourStr))
	if err != nil || endHour < 0 || endHour > 23 {
		return fmt.Errorf("invalid end hour. Must be between 0 and 23")
	}

	// Get suggestion interval
	fmt.Print("Enter suggestion interval in minutes (e.g., 20, 30, 60): ")
	var intervalStr string
	fmt.Scanln(&intervalStr)
	interval, err := strconv.Atoi(strings.TrimSpace(intervalStr))
	if err != nil || interval <= 0 || interval > 1440 {
		return fmt.Errorf("invalid interval. Must be between 1 and 1440 minutes")
	}

	// Save configuration
	if err := SetDelayedCommitConfig(startHour, endHour, interval); err != nil {
		return fmt.Errorf("failed to save delayed commit configuration: %w", err)
	}

	fmt.Println()
	fmt.Println("Delayed commit settings configured successfully!")
	fmt.Printf("Restricted hours: %02d:00 - %02d:00\n", startHour, endHour)
	fmt.Printf("Suggestion interval: %d minutes\n", interval)
	return nil
}
