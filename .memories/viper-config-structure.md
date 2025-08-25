---
title: Viper Config Structure
description: Configuration structures with mapstructure tags for Viper
tags:
  configuration: true
  golang: true
  viper: true
created: 2025-08-24T16:23:35.098201929-07:00
modified: 2025-08-24T16:23:35.098201929-07:00
---

# Viper Config Structure

Configuration structure implementation in `internal/config/config.go`.

## ApiKeyConfig
```go
type ApiKeyConfig struct {
    Value string `mapstructure:"-"`
    Path  string `mapstructure:"path"`
}
```

Supports loading API keys from:
- Direct value in config
- File path (with `~/` expansion)
- Environment variables

## VoyageAIConfig
```go
type VoyageAIConfig struct {
    ApiKey      ApiKeyConfig `mapstructure:"api_key"`
    Model       string       `mapstructure:"model"`
    RerankModel string       `mapstructure:"rerank_model"`
}
```

Default values:
- `model`: `"voyage-3.5"`
- `rerank_model`: `"rerank-lite-1"`

## TOML Configuration Example
```toml
[voyage_ai]
model = "voyage-3.5"
rerank_model = "rerank-lite-1"

[voyage_ai.api_key]
path = "~/.config/simplemem/voyage_ai_key"
```

## Loading Process
1. Viper unmarshals config into structs
2. API key resolution (file path → value)
3. Environment variable override
4. Validation and defaults applied

## Related Components
- [[viper-toml-migration-overview]] - Migration overview
- [[viper-cobra-integration]] - CLI integration
- [[viper-migration-examples]] - Usage examples