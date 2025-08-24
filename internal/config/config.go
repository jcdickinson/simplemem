package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ApiKeyConfig represents either a direct API key string or a path to a file containing the key
type ApiKeyConfig struct {
	Value string `json:"-"`
	raw   interface{}
}

// UnmarshalJSON implements json.Unmarshaler to handle both string and object forms
func (a *ApiKeyConfig) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as string first
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		a.Value = str
		a.raw = str
		return nil
	}

	// Try to unmarshal as object with "path" field
	var obj struct {
		Path string `json:"path"`
	}
	if err := json.Unmarshal(data, &obj); err != nil {
		return fmt.Errorf("api_key must be either a string or an object with 'path' field")
	}

	if obj.Path == "" {
		return fmt.Errorf("api_key object must have non-empty 'path' field")
	}

	// Read the key from the file
	keyBytes, err := os.ReadFile(obj.Path)
	if err != nil {
		return fmt.Errorf("failed to read API key from file %s: %w", obj.Path, err)
	}

	a.Value = strings.TrimSpace(string(keyBytes))
	a.raw = obj
	return nil
}

// MarshalJSON implements json.Marshaler
func (a ApiKeyConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.raw)
}

// VoyageAIConfig holds configuration for VoyageAI integration
type VoyageAIConfig struct {
	ApiKey      ApiKeyConfig `json:"api_key"`
	Model       string       `json:"model,omitempty"`
	RerankModel string       `json:"rerank_model,omitempty"`
}

// Config represents the complete simplemem configuration
type Config struct {
	VoyageAI VoyageAIConfig `json:"voyage_ai"`
}

// getConfigDirs returns the list of configuration directories to check, in order of precedence
func getConfigDirs() []string {
	var dirs []string

	// 1. XDG_CONFIG_DIRS (system-wide configs)
	if xdgConfigDirs := os.Getenv("XDG_CONFIG_DIRS"); xdgConfigDirs != "" {
		for _, dir := range strings.Split(xdgConfigDirs, ":") {
			if dir != "" {
				dirs = append(dirs, filepath.Join(dir, "simplemem"))
			}
		}
	}

	// 2. XDG_CONFIG_HOME or ~/.config (user config)
	var userConfigDir string
	if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
		userConfigDir = xdgConfigHome
	} else {
		if homeDir, err := os.UserHomeDir(); err == nil {
			userConfigDir = filepath.Join(homeDir, ".config")
		}
	}
	if userConfigDir != "" {
		dirs = append(dirs, filepath.Join(userConfigDir, "simplemem"))
	}

	return dirs
}

// Load loads configuration from the standard config directories and project-specific overrides
func Load() (*Config, error) {
	config := &Config{}

	// Load from system/user config directories
	configDirs := getConfigDirs()
	for _, dir := range configDirs {
		configPath := filepath.Join(dir, "config.json")
		if err := loadConfigFile(configPath, config); err != nil {
			if !os.IsNotExist(err) {
				return nil, fmt.Errorf("failed to load config from %s: %w", configPath, err)
			}
		}
	}

	// Load project-specific config overlay
	projectConfigPath := ".config/simplemem/config.json"
	if err := loadConfigFile(projectConfigPath, config); err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to load project config from %s: %w", projectConfigPath, err)
		}
	}

	return config, nil
}

// loadConfigFile loads a single config file and merges it into the existing config
func loadConfigFile(path string, config *Config) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Parse into a temporary config to overlay on top of existing
	var overlay Config
	if err := json.Unmarshal(data, &overlay); err != nil {
		return fmt.Errorf("invalid JSON in config file: %w", err)
	}

	// Merge the overlay into the existing config
	mergeConfig(config, &overlay)
	return nil
}

// mergeConfig merges the overlay config into the base config
func mergeConfig(base, overlay *Config) {
	// For now, we only have VoyageAI config, so we can do a simple field-by-field merge
	// In the future, this could be made more generic using reflection
	
	if overlay.VoyageAI.ApiKey.raw != nil {
		base.VoyageAI.ApiKey = overlay.VoyageAI.ApiKey
	}
	if overlay.VoyageAI.Model != "" {
		base.VoyageAI.Model = overlay.VoyageAI.Model
	}
	if overlay.VoyageAI.RerankModel != "" {
		base.VoyageAI.RerankModel = overlay.VoyageAI.RerankModel
	}
}