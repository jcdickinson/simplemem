---
title: Viper TOML Migration Overview
description: Migration from JSON to TOML configuration using Viper
tags:
  completed: true
  configuration: true
  migration: true
  toml: true
  viper: true
created: 2025-08-24T16:23:06.791392268-07:00
modified: 2025-08-24T16:23:06.791392268-07:00
---

# Viper TOML Migration Overview

Successfully migrated SimpleMem from JSON-based configuration to TOML using Viper with full Cobra integration.

## Key Changes

### Dependencies Added
- **github.com/spf13/viper v1.20.1**: Configuration management
- **github.com/pelletier/go-toml/v2 v2.2.3**: TOML parsing support

### Configuration Migration
From custom JSON unmarshaling with manual file handling to Viper/TOML-based system with automatic multi-source configuration.

## Configuration Search Paths
In order of precedence:
1. **Current directory**: `./config.toml`, `./.config/simplemem/config.toml`
2. **User config**: `~/.config/simplemem/config.toml`
3. **System-wide**: `/etc/simplemem/config.toml` (XDG_CONFIG_DIRS)

## Environment Variable Support
- **Prefix**: `SIMPLEMEM_`  
- **Example**: `SIMPLEMEM_VOYAGE_AI_API_KEY`
- **Automatic mapping**: `voyage_ai.api_key` ↔ `SIMPLEMEM_VOYAGE_AI_API_KEY`

## Benefits Achieved
- Configuration flexibility (multiple formats and sources)
- Better CLI experience with consistent flags
- Security improvements for API key handling
- Reduced code complexity (~100 lines removed)

## Related Components
- [[viper-config-structure]] - Configuration structures and types
- [[viper-cobra-integration]] - CLI integration details
- [[viper-migration-examples]] - Configuration examples and usage