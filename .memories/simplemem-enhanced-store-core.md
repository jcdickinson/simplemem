---
title: SimpleMem Enhanced Store Core
description: Core architecture and composition of the Enhanced Store
tags:
  architecture: true
  enhanced-store: true
  golang: true
  storage: true
created: 2025-08-24T16:15:49.048164118-07:00
modified: 2025-08-24T16:15:49.048164118-07:00
---

# SimpleMem Enhanced Store Core

The Enhanced Store (`internal/memory/enhanced_store.go`) bridges file-based storage with RAG capabilities.

## Composition Pattern

```go
type EnhancedStore struct {
    *Store                    // Embedded basic store
    db          *db.DB        // Database layer
    ragProcessor *rag.Processor // RAG processing
    dbPath      string        // Database file path
}
```

## Dual-Storage Architecture

### File System (Primary)
- Location: `.memories/` directory
- Format: Markdown files with YAML frontmatter
- Purpose: Human-readable, version-controllable storage

### Database (Secondary)
- Location: `.cache/simplemem.db` 
- Format: DuckDB with vector extensions
- Purpose: Structured queries, vector similarity, indexing
- Sync: Automatic using SHA256 hashing

## Initialization Process

```go
func NewEnhancedStore(basePath string, cfg *config.Config) (*EnhancedStore, error)
```

1. Create Basic Store - file system operations
2. Initialize Database - DuckDB with vector extensions
3. Create RAG Processor - embedding and search capabilities
4. Sync Existing Files - import memories into database
5. Process Embeddings - generate vectors for content

## Related Components
- [[enhanced-store-memory-operations]] - CRUD operations
- [[enhanced-store-search-capabilities]] - Search implementations
- [[enhanced-store-sync-strategy]] - File-database synchronization