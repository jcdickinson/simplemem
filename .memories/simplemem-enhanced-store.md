---
title: SimpleMem Enhanced Store Implementation
description: Deep dive into the RAG-enabled Enhanced Store that wraps basic file operations with semantic search capabilities
tags:
  architecture: true
  embeddings: true
  golang: true
  rag: true
  storage: true
created: 2025-08-24T00:20:07.623783532-07:00
modified: 2025-08-24T00:20:07.623783532-07:00
---

# SimpleMem Enhanced Store Implementation

The Enhanced Store (`internal/memory/enhanced_store.go`) is the core component that bridges traditional file-based memory storage with advanced RAG (Retrieval-Augmented Generation) capabilities.

## Architecture Overview

### Composition Pattern

```go
type EnhancedStore struct {
    *Store                    // Embedded basic store
    db          *db.DB        // Database layer
    ragProcessor *rag.Processor // RAG processing
    dbPath      string        // Database file path
}
```

The Enhanced Store uses composition to:
- **Embed Basic Store**: Inherit all file system operations
- **Add Database Layer**: DuckDB for structured storage and vector operations
- **Integrate RAG Processor**: Semantic search and similarity matching

## Core Functionality

### Dual-Storage Architecture

#### File System (Primary)
- **Location**: `.memories/` directory
- **Format**: Markdown files with YAML frontmatter
- **Purpose**: Human-readable, version-controllable storage
- **Thread-Safety**: Handled by embedded Store

#### Database (Secondary)
- **Location**: `.cache/simplemem.db` 
- **Format**: DuckDB with vector extensions
- **Purpose**: Structured queries, vector similarity, indexing
- **Sync Strategy**: Automatic sync on file changes using SHA256 hashing

### Initialization Process

```go
func NewEnhancedStore(basePath string, cfg *config.Config) (*EnhancedStore, error)
```

1. **Create Basic Store**: Initialize file system operations
2. **Initialize Database**: Set up DuckDB with vector extensions
3. **Create RAG Processor**: Initialize embedding and search capabilities
4. **Sync Existing Files**: Import existing memories into database
5. **Process Embeddings**: Generate vectors for existing content

### Memory Operations

#### Create Operation
```go
func (es *EnhancedStore) Create(name, content string) error
```
1. Create file using basic store
2. Sync to database with metadata extraction
3. Process with RAG (generate embeddings, extract links)
4. Handle errors gracefully (file operations take precedence)

#### Update Operation  
```go
func (es *EnhancedStore) Update(name, content string) error
```
1. Update file using basic store
2. Check content hash to detect changes
3. Re-sync to database if changed
4. Regenerate embeddings if content modified

#### Delete Operation
```go
func (es *EnhancedStore) Delete(name string) error
```
1. Delete from file system first
2. Clean up all database entries (cascading delete)
3. Remove embeddings and related data
4. Continue even if database cleanup fails

## Search Capabilities

### Semantic Search

#### Basic Semantic Search
```go
func (es *EnhancedStore) SearchSemantic(query string, limit int) ([]MemoryInfo, []float32, error)
```
- Generate embeddings for search query
- Perform vector similarity search in database
- Return memories with similarity scores

#### Tag-Filtered Semantic Search
```go
func (es *EnhancedStore) SearchSemanticWithTags(query string, tagFilters map[string]string, requireAll bool, limit int) ([]MemoryInfo, []float32, error)
```
- Apply tag filters before semantic search
- Support presence-only filters (empty value)
- Support exact value matching
- Boolean logic: ALL tags vs ANY tags

### Hybrid Search

#### Algorithm Overview
```go
func (es *EnhancedStore) SearchHybridWithTags(query string, tagFilters map[string]string, requireAll bool, limit int) ([]MemoryInfo, []float32, error)
```

1. **Keyword Search**: Traditional text matching using basic store
2. **Semantic Search**: Vector similarity with increased limit for coverage
3. **Result Combination**: Merge and deduplicate results
4. **Score Boosting**: Keyword matches get 1.2x similarity boost
5. **Tag Filtering**: Apply filters to both result sets
6. **Final Ranking**: Sort by combined relevance score

#### Scoring Strategy
- **Semantic Only**: Original similarity score (0.0-1.0)
- **Keyword Only**: Fixed high score (0.9)
- **Both**: Boosted similarity score (similarity * 1.2, max 1.0)

### Backlink Analysis

#### Enhanced Backlinks
```go  
func (es *EnhancedStore) GetEnhancedBacklinks(memoryName string, query string, limit int) (string, error)
```
- Combines explicit links (markdown/wiki-style) with semantic similarity
- Optionally re-ranks by query relevance
- Returns formatted markdown with context snippets

#### Semantic Backlinks
```go
func (es *EnhancedStore) GetSemanticBacklinks(name string, minSimilarity float32) ([]MemoryInfo, []float32, error)
```
- Finds memories semantically similar to target memory
- Uses configurable similarity threshold
- Excludes the source memory from results

## Synchronization Strategy

### File-to-Database Sync

#### Change Detection
```go
// Calculate SHA256 hash for content change detection
hash := fmt.Sprintf("%x", sha256.Sum256([]byte(memInfo.Content)))
```

#### Sync Process
1. **Read File**: Load content and parse frontmatter
2. **Calculate Hash**: Generate SHA256 of content
3. **Check Existing**: Query database for existing record
4. **Compare Hashes**: Skip if content unchanged
5. **Upsert Memory**: Update or insert memory record
6. **Sync Tags**: Update normalized tag table
7. **Process RAG**: Generate embeddings and extract links

### Batch Operations

#### Initial Sync
```go
func (es *EnhancedStore) syncFilesToDatabase() error
```
- Processes all existing files on startup
- Handles errors gracefully (warns but continues)
- Used during Enhanced Store initialization

## Error Handling Philosophy

### Graceful Degradation
- **File Operations**: Always succeed if possible
- **Database Operations**: Log warnings, don't fail primary operation
- **RAG Processing**: Continue without embeddings if service unavailable
- **Configuration Issues**: Disable RAG features, maintain basic functionality

### Error Logging Strategy
```go
log.Printf("Warning: failed to sync memory to database: %v", err)
```
- Non-critical errors logged as warnings
- Critical errors propagated to caller
- Detailed context provided for debugging

## Performance Optimizations

### Caching
- **Content Hash Checking**: Avoid unnecessary processing
- **Database Query Optimization**: Indexed searches
- **Embedding Caching**: Avoid regeneration of unchanged content

### Concurrent Operations
- **Thread-Safe**: All operations safe for concurrent use
- **Background Processing**: RAG operations don't block file operations
- **Batch Processing**: Efficient initial setup and bulk operations

The Enhanced Store successfully provides a seamless upgrade path from basic file storage to advanced RAG capabilities while maintaining backward compatibility and operational reliability.