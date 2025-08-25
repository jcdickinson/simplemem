---
title: Enhanced Store Sync Strategy
description: File-to-database synchronization in SimpleMem Enhanced Store
tags:
  database: true
  enhanced-store: true
  synchronization: true
created: 2025-08-24T16:17:30.825021439-07:00
modified: 2025-08-24T16:17:30.825021439-07:00
---

# Enhanced Store Sync Strategy

How the Enhanced Store keeps file system and database in sync.

## Change Detection

### Content Hashing
```go
// Calculate SHA256 hash for content change detection
hash := fmt.Sprintf("%x", sha256.Sum256([]byte(memInfo.Content)))
```

## Sync Process

1. **Read File**: Load content and parse frontmatter
2. **Calculate Hash**: Generate SHA256 of content
3. **Check Existing**: Query database for existing record
4. **Compare Hashes**: Skip if content unchanged
5. **Upsert Memory**: Update or insert memory record
6. **Sync Tags**: Update normalized tag table
7. **Process RAG**: Generate embeddings and extract links

## Batch Operations

### Initial Sync
```go
func (es *EnhancedStore) syncFilesToDatabase() error
```
- Processes all existing files on startup
- Handles errors gracefully (warns but continues)
- Used during Enhanced Store initialization

## Performance Optimizations

### Caching Strategies
- Content Hash Checking: Avoid unnecessary processing
- Database Query Optimization: Indexed searches
- Embedding Caching: Avoid regeneration of unchanged content

### Concurrent Operations
- Thread-Safe: All operations safe for concurrent use
- Background Processing: RAG operations don't block file operations
- Batch Processing: Efficient initial setup and bulk operations

## Related Components
- [[simplemem-enhanced-store-core]] - Core architecture
- [[enhanced-store-memory-operations]] - CRUD operations
- [[simplemem-database-schema]] - Database structure