---
title: Viper Cobra Integration
description: Integration of Viper configuration with Cobra CLI framework
tags:
  cli: true
  cobra: true
  configuration: true
  viper: true
created: 2025-08-24T16:24:01.576084907-07:00
modified: 2025-08-24T16:24:01.576084907-07:00
---

# Viper Cobra Integration

How Viper configuration is integrated with Cobra CLI in `cmd/root.go`.

## New CLI Flags
- `--config`: Specify custom config file path
- `--db`: Database path (existing, now Viper-bound)

## Initialization Flow
```go
func init() {
    cobra.OnInitialize(initConfig)
    // Flag binding to Viper
    viper.BindPFlag("db", rootCmd.PersistentFlags().Lookup("db"))
}

func initConfig() {
    if configFile != "" {
        viper.SetConfigFile(configFile)
    }
    config.InitializeViper()
}
```

## Configuration Precedence
1. CLI flags (highest priority)
2. Environment variables
3. Config file values
4. Default values (lowest priority)

## Key Features
- Automatic flag-to-config binding
- Environment variable support
- Multiple config file locations
- Help text integration

## Usage Examples
```bash
# Use custom config file
simplemem --config /path/to/config.toml

# Override database path
simplemem --db /custom/path/simplemem.db

# Environment variable override
SIMPLEMEM_VOYAGE_AI_API_KEY="key" simplemem
```

## Related Components
- [[viper-toml-migration-overview]] - Migration overview
- [[viper-config-structure]] - Configuration structures
- [[viper-migration-examples]] - More usage examples