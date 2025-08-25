---
title: SimpleMem Database Schema
description: DuckDB schema design for memories, tags, and embeddings
tags:
  database: true
  duckdb: true
  schema: true
created: 2025-08-24T16:19:05.710697182-07:00
modified: 2025-08-24T18:48:30.288151539-07:00
---

# SimpleMem Database Schema

DuckDB schema for structured storage and vector operations.

## Core Tables

### memories table
- `id`: Primary key (auto-increment)
- `name`: Unique memory identifier/filename
- `title`: Optional display title from frontmatter
- `description`: Optional description from frontmatter  
- `content`: Full markdown content including frontmatter
- `body`: Content without frontmatter
- `file_hash`: SHA-256 hash for change detection
- `created`: Creation timestamp from frontmatter
- `modified`: Last modification timestamp from frontmatter

### tags table
- `memory_id`: Foreign key to memories table
- `key`: Tag key/name
- `value`: Tag value (boolean or string)

### chunks table
- `id`: Primary key (auto-increment)
- `memory_id`: Foreign key to memories table
- `content`: Chunk text content
- `chunk_index`: Position within memory (0-based)
- `embedding`: Vector embedding (FLOAT[1536])

## Relationships
- One memory → many tags (one-to-many)
- One memory → many chunks (one-to-many)
- Each chunk belongs to exactly one memory

## Index Strategy

### Primary Indexes
- `memories.name`: Unique index for fast lookup
- `memories.file_hash`: Index for change detection
- `tags.memory_id`: Foreign key index for filtering
- `chunks.memory_id`: Foreign key index for retrieval

### Vector Indexes
- Automatic vector similarity indexes on `chunks.embedding`
- Optimized for cosine similarity operations
- Background index maintenance

## Related Components
- [[database-vector-search]] - Vector search implementation
- [[database-operations]] - CRUD operations
- [[simplemem-enhanced-store-core]] - Enhanced store integration