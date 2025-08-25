package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// ApiKeyConfig represents either a direct API key string or a path to a file containing the key
type ApiKeyConfig struct {
	Value string `mapstructure:"-"`
	Path  string `mapstructure:"path"`
}

// VoyageAIConfig holds configuration for VoyageAI integration
type VoyageAIConfig struct {
	ApiKey      ApiKeyConfig `mapstructure:"api_key"`
	Model       string       `mapstructure:"model"`
	RerankModel string       `mapstructure:"rerank_model"`
}

// Config represents the complete simplemem configuration
type Config struct {
	VoyageAI        VoyageAIConfig `mapstructure:"voyage_ai"`
	MaxMemoryLength int            `mapstructure:"max_memory_length"`
}

// InitializeViper sets up Viper configuration with proper search paths and defaults
func InitializeViper() error {
	// Set config name and format
	viper.SetConfigName("config")
	viper.SetConfigType("toml")

	// Add config search paths in order of precedence
	// 1. Current directory (project-specific config)
	viper.AddConfigPath(".")
	viper.AddConfigPath(".config/simplemem")

	// 2. XDG_CONFIG_HOME or ~/.config (user config)
	if xdgConfigHome := os.Getenv("XDG_CONFIG_HOME"); xdgConfigHome != "" {
		viper.AddConfigPath(filepath.Join(xdgConfigHome, "simplemem"))
	} else {
		if homeDir, err := os.UserHomeDir(); err == nil {
			viper.AddConfigPath(filepath.Join(homeDir, ".config", "simplemem"))
		}
	}

	// 3. XDG_CONFIG_DIRS (system-wide configs)
	if xdgConfigDirs := os.Getenv("XDG_CONFIG_DIRS"); xdgConfigDirs != "" {
		for dir := range strings.SplitSeq(xdgConfigDirs, ":") {
			if dir != "" {
				viper.AddConfigPath(filepath.Join(dir, "simplemem"))
			}
		}
	}

	// Set defaults
	viper.SetDefault("voyage_ai.model", "voyage-3.5")
	viper.SetDefault("voyage_ai.rerank_model", "rerank-lite-1")
	viper.SetDefault("max_memory_length", 2500)

	// Enable environment variable support
	viper.SetEnvPrefix("SIMPLEMEM")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Try to read config file (it's okay if it doesn't exist)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("failed to read config file: %w", err)
		}
	}

	return nil
}

// stringToApiKeyConfigHookFunc returns a mapstructure.DecodeHookFunc that converts
// string values to ApiKeyConfig structs
func stringToApiKeyConfigHookFunc() mapstructure.DecodeHookFunc {
	return func(f, t reflect.Type, data interface{}) (interface{}, error) {
		if t != reflect.TypeOf(ApiKeyConfig{}) {
			return data, nil
		}

		switch f.Kind() {
		case reflect.String:
			// If it's a string, treat it as a direct API key value
			return ApiKeyConfig{Value: data.(string)}, nil
		case reflect.Map:
			// If it's a map, let the normal unmarshaling handle it
			return data, nil
		default:
			return data, nil
		}
	}
}

// Load loads configuration using Viper and resolves API keys
func Load() (*Config, error) {
	if err := InitializeViper(); err != nil {
		return nil, err
	}

	var config Config
	
	// Configure decoder with custom hook for ApiKeyConfig
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: stringToApiKeyConfigHookFunc(),
		Result:     &config,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create decoder: %w", err)
	}

	if err := decoder.Decode(viper.AllSettings()); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Resolve API key from file if path is specified
	if err := resolveApiKey(&config.VoyageAI.ApiKey); err != nil {
		return nil, fmt.Errorf("failed to resolve VoyageAI API key: %w", err)
	}

	return &config, nil
}

// resolveApiKey resolves the API key from file if path is specified, otherwise uses direct value from environment
func resolveApiKey(apiKey *ApiKeyConfig) error {
	// First, try to get directly from environment variable
	if envKey := viper.GetString("voyage_ai.api_key"); envKey != "" {
		// If it's not a path (doesn't start with / or ./ or ~/), treat as direct value
		if !strings.HasPrefix(envKey, "/") && !strings.HasPrefix(envKey, "./") && !strings.HasPrefix(envKey, "~/") {
			apiKey.Value = envKey
			return nil
		}
		// Otherwise treat as path
		apiKey.Path = envKey
	}

	// If we have a path, read the key from file
	if apiKey.Path != "" {
		// Handle ~ expansion
		if strings.HasPrefix(apiKey.Path, "~/") {
			if homeDir, err := os.UserHomeDir(); err == nil {
				apiKey.Path = filepath.Join(homeDir, apiKey.Path[2:])
			}
		}

		keyBytes, err := os.ReadFile(apiKey.Path)
		if err != nil {
			return fmt.Errorf("failed to read API key from file %s: %w", apiKey.Path, err)
		}
		apiKey.Value = strings.TrimSpace(string(keyBytes))
	}

	return nil
}
