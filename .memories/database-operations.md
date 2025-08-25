---
title: Database Operations
description: CRUD operations and data management in SimpleMem database layer
tags:
  crud: true
  database: true
  operations: true
created: 2025-08-24T16:20:39.334127778-07:00
modified: 2025-08-24T16:20:39.334127778-07:00
---

# Database Operations

Core database operations for memory management in SimpleMem.

## Memory Management

### Memory Storage
```go
func (db *DB) UpsertMemory(memory *Memory) error
```

Handles INSERT and UPDATE operations:
- Uses file hash to detect content changes
- Preserves existing ID on updates
- Triggers embedding regeneration when content changes

### Tag Management
```go
func (db *DB) UpsertTags(memoryID int64, tags map[string]interface{}) error
```

Tag synchronization:
1. Delete existing tags for the memory
2. Insert new tags from current frontmatter
3. Support both boolean and string tag values

## Chunk Operations

### Chunk Storage
```go
func (db *DB) UpsertChunks(memoryID int64, chunks []ChunkData) error
```

Chunk management:
1. Remove existing chunks for the memory
2. Batch insert new chunks with embeddings
3. Refresh vector similarity indexes

## Data Consistency

### Transaction Management
- Atomic Updates: Memory, tags, and chunks in single transaction
- Error Handling: Rollback on any failure during upsert
- Consistency: Foreign key constraints ensure referential integrity

### Change Detection
File hash comparison:
1. Calculate SHA-256 of file content
2. Retrieve stored hash from database
3. Skip processing if hashes match
4. Refresh database if content changed

### Cleanup Operations
```go
func (db *DB) DeleteMemory(name string) error
```

Cascading delete:
1. Remove all embedding chunks
2. Remove all tag associations  
3. Remove memory record
4. Automatic index cleanup

## Related Components
- [[simplemem-database-schema]] - Schema design
- [[database-vector-search]] - Vector search operations
- [[enhanced-store-sync-strategy]] - File-database sync