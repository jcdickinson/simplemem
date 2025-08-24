---
title: 'Knowledge: Project Dependencies Analysis'
tags:
  area: dependencies
  component: go-modules
  knowledge: true
created: 2025-08-24T00:29:29.179475232-07:00
modified: 2025-08-24T00:29:29.179475232-07:00
---

# SimpleMem Dependency Analysis

## Current Dependencies (from go.mod)

### Core MCP Framework
- `github.com/ThinkInAIXYZ/go-mcp v0.2.20` - MCP protocol implementation

### Markdown Processing  
- `github.com/gomarkdown/markdown v0.0.0-20250810172220-2e2c11897d1a` - Markdown parsing

### Database
- `github.com/marcboeker/go-duckdb v1.8.5` - DuckDB Go driver for vector storage

### Configuration
- `gopkg.in/yaml.v3 v3.0.1` - YAML processing for frontmatter

## Notable Dependencies
### Apache Arrow (Indirect)
- `github.com/apache/arrow-go/v18 v18.1.0` - Used by DuckDB for columnar operations

### JSON Processing
- `github.com/goccy/go-json v0.10.5` - High-performance JSON
- `github.com/tidwall/gjson v1.18.0` - Fast JSON queries

### UUID Generation
- `github.com/google/uuid v1.6.0` - For unique identifiers

## Missing Dependencies
Based on plan.md Phase 2 requirements:
- **VoyageAI Client**: No `github.com/voyage-ai/voyageai-go` dependency found
- May need to implement custom HTTP client for VoyageAI API

## Build System
- Uses Go 1.24 (cutting edge)
- Just-based build system (justfile)
- Nix development environment (flake.nix)

## Related  
- [[simplemem-project-overview]]
- [[todo-implement-missing-rag-components]]