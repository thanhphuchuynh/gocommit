package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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
