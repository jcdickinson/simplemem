---
title: Viper Migration Examples
description: Configuration examples and migration path for Viper TOML setup
tags:
  configuration: true
  examples: true
  migration: true
  viper: true
created: 2025-08-24T16:24:34.420291762-07:00
modified: 2025-08-24T16:24:34.420291762-07:00
---

# Viper Migration Examples

Configuration examples and migration guidance for the Viper TOML system.

## Configuration File Examples

### Basic TOML Structure
```toml
[voyage_ai]
model = "voyage-3.5"
rerank_model = "rerank-lite-1"

[voyage_ai.api_key]
path = "~/.config/simplemem/voyage_ai_key"
```

### Environment Variables
```bash
export SIMPLEMEM_VOYAGE_AI_API_KEY="your-api-key"
export SIMPLEMEM_VOYAGE_AI_MODEL="voyage-3.5"
```

## Migration Path for Existing Users

### Backward Compatibility
1. **JSON configs still work**: Viper supports multiple formats
2. **Environment variables**: Same prefix pattern maintained  
3. **Existing workflows**: Unchanged functionality

### Recommended Migration Steps
1. Copy `config.toml.example` to desired location
2. Configure API keys via environment variables or file paths
3. Use `--config` flag for project-specific configs

## Files Modified

### New Files
- `config.toml.example`: Example TOML configuration
- Enhanced CLI with `--config` flag support

### Modified Files
- `internal/config/config.go`: Complete Viper rewrite
- `cmd/root.go`: Cobra-Viper integration
- `go.mod`: Added Viper dependencies

## Testing Results
✅ Configuration loading with Viper  
✅ CLI flag functionality  
✅ Environment variable override  
✅ Build and compilation success

## Related Components
- [[viper-toml-migration-overview]] - Migration overview
- [[viper-config-structure]] - Configuration structures
- [[viper-cobra-integration]] - CLI integration