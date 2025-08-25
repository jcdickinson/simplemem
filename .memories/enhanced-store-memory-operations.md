---
title: Enhanced Store Memory Operations
description: CRUD operations in SimpleMem Enhanced Store
tags:
  crud: true
  enhanced-store: true
  operations: true
created: 2025-08-24T16:16:23.368859805-07:00
modified: 2025-08-24T16:16:23.368859805-07:00
---

# Enhanced Store Memory Operations

How the Enhanced Store handles memory CRUD operations with database sync and RAG processing.

## Create Operation
```go
func (es *EnhancedStore) Create(name, content string) error
```
1. Create file using basic store
2. Sync to database with metadata extraction
3. Process with RAG (generate embeddings, extract links)
4. Handle errors gracefully (file operations take precedence)

## Update Operation  
```go
func (es *EnhancedStore) Update(name, content string) error
```
1. Update file using basic store
2. Check content hash to detect changes
3. Re-sync to database if changed
4. Regenerate embeddings if content modified

## Delete Operation
```go
func (es *EnhancedStore) Delete(name string) error
```
1. Delete from file system first
2. Clean up all database entries (cascading delete)
3. Remove embeddings and related data
4. Continue even if database cleanup fails

## Error Handling Philosophy

### Graceful Degradation
- File Operations: Always succeed if possible
- Database Operations: Log warnings, don't fail primary operation
- RAG Processing: Continue without embeddings if service unavailable
- Configuration Issues: Disable RAG features, maintain basic functionality

### Error Logging Strategy
```go
log.Printf("Warning: failed to sync memory to database: %v", err)
```
- Non-critical errors logged as warnings
- Critical errors propagated to caller
- Detailed context provided for debugging

## Related Components
- [[simplemem-enhanced-store-core]] - Core architecture
- [[enhanced-store-sync-strategy]] - Synchronization details
- [[mcp-memory-management-tools]] - MCP tool interfaces