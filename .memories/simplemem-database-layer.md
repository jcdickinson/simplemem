---
title: SimpleMem Database Layer
description: Analysis of the DuckDB integration and schema design for structured storage and vector operations
tags:
  architecture: true
  database: true
  duckdb: true
  golang: true
  storage: true
  vectors: true
created: 2025-08-24T00:22:16.675733747-07:00
modified: 2025-08-24T01:02:29.476549454-07:00
---

# SimpleMem Database Layer

The database layer (`internal/db/`) provides structured storage and vector operations using DuckDB as the backend database, enabling both relational queries and semantic search capabilities.

## Database Architecture

### Core Database Type

```go
type DB struct {
    conn *sql.DB  // Database connection handle
}
```

### Schema Design

#### Primary Tables

**memories table**
- `id`: Primary key (auto-increment)
- `name`: Unique memory identifier/filename
- `title`: Optional display title from frontmatter
- `description`: Optional description from frontmatter  
- `content`: Full markdown content including frontmatter
- `body`: Content without frontmatter (for display/search)
- `file_hash`: SHA-256 hash for change detection
- `created`: Creation timestamp from frontmatter
- `modified`: Last modification timestamp from frontmatter

**tags table**
- `memory_id`: Foreign key to memories table
- `key`: Tag key/name
- `value`: Tag value (can be boolean or string)

**chunks table**
- `id`: Primary key (auto-increment)
- `memory_id`: Foreign key to memories table
- `content`: Chunk text content
- `chunk_index`: Position within the memory (0-based)
- `embedding`: Vector embedding (FLOAT[1536])

#### Relationships
- One memory can have many tags (one-to-many)
- One memory can have many chunks (one-to-many)
- Each chunk belongs to exactly one memory

## Vector Search Integration

### DuckDB Vector Extension

The database leverages DuckDB's vector extension for embedding storage and similarity search:

```sql
-- Vector similarity search using cosine distance
SELECT m.*, c.content, 1 - array_cosine_similarity(c.embedding, $1) as similarity
FROM memories m 
JOIN chunks c ON m.id = c.memory_id 
ORDER BY similarity DESC 
LIMIT ?
```

### Embedding Operations

#### Similarity Search
```go
func (db *DB) SearchSimilarMemoriesWithTags(queryEmbedding []float32, tagFilters []TagFilter, requireAll bool, limit int) ([]Memory, []float32, error)
```

The search process:
1. **Vector Query**: Find chunks with embeddings similar to query
2. **Tag Filtering**: Apply metadata filters using SQL WHERE clauses
3. **Aggregation**: Group results by memory, selecting best chunk per memory
4. **Ranking**: Sort by similarity score (cosine distance)
5. **Limit**: Return top N results

#### Tag Filter Logic
```go
type TagFilter struct {
    Key        string
    Value      string  
    CheckValue bool    // If false, only check key presence
}
```

- **Presence Check**: `tags.key = ? AND tags.value IS NOT NULL`
- **Value Check**: `tags.key = ? AND tags.value = ?`
- **Boolean Logic**: Support AND/OR combinations via `requireAll` parameter

## Database Operations

### Memory Management

#### Memory Storage
```go
func (db *DB) UpsertMemory(memory *Memory) error
```

Handles both INSERT and UPDATE operations:
- Uses file hash to detect content changes
- Preserves existing ID on updates
- Triggers embedding regeneration when content changes

#### Tag Management
```go
func (db *DB) UpsertTags(memoryID int64, tags map[string]interface{}) error
```

Tag synchronization process:
1. **Delete Existing**: Remove all tags for the memory
2. **Insert New**: Add tags from current frontmatter
3. **Type Handling**: Support both boolean and string tag values

### Chunk Operations

#### Chunk Storage
```go
func (db *DB) UpsertChunks(memoryID int64, chunks []ChunkData) error
```

Chunk management workflow:
1. **Cleanup**: Remove existing chunks for the memory
2. **Batch Insert**: Add new chunks with embeddings
3. **Index Update**: Refresh vector similarity indexes

#### Embedding Storage
- **Vector Format**: FLOAT[1536] array for VoyageAI embeddings
- **Indexing**: Automatic vector similarity indexes for fast search
- **Compression**: DuckDB's built-in vector compression

## Query Optimization

### Index Strategy

**Primary Indexes**
- `memories.name`: Unique index for fast lookup by filename
- `memories.file_hash`: Index for change detection
- `tags.memory_id`: Foreign key index for tag filtering
- `chunks.memory_id`: Foreign key index for chunk retrieval

**Vector Indexes**
- Automatic vector similarity indexes on `chunks.embedding`
- Optimized for cosine similarity operations
- Background index maintenance

### Query Patterns

#### Semantic Search Query
```sql
WITH chunk_similarities AS (
  SELECT c.memory_id, c.content, 
         1 - array_cosine_similarity(c.embedding, ?) as similarity
  FROM chunks c
  WHERE c.embedding IS NOT NULL
), best_chunks AS (
  SELECT memory_id, MAX(similarity) as best_similarity
  FROM chunk_similarities  
  GROUP BY memory_id
)
SELECT m.*, cs.similarity
FROM memories m
JOIN best_chunks bc ON m.id = bc.memory_id
JOIN chunk_similarities cs ON cs.memory_id = m.id AND cs.similarity = bc.best_similarity
ORDER BY similarity DESC
LIMIT ?
```

#### Tag-Filtered Search
```sql
-- Semantic search with tag filtering (requireAll = true)
SELECT DISTINCT m.*, similarity
FROM (previous query) m
WHERE m.id IN (
  SELECT memory_id 
  FROM tags 
  WHERE (key = ? AND value = ?) OR (key = ? AND value = ?)
  GROUP BY memory_id
  HAVING COUNT(DISTINCT key) = ?  -- Number of required tags
)
```

## Performance Characteristics

### Memory Usage
- **Embedding Storage**: ~6KB per chunk (1536 floats × 4 bytes)
- **Index Overhead**: ~20% additional storage for vector indexes
- **Connection Pool**: Single connection with statement caching

### Query Performance
- **Vector Search**: Sub-second for collections up to 100K chunks
- **Tag Filtering**: Millisecond response with proper indexing
- **Combined Queries**: Optimized execution plans for hybrid operations

### Scalability Considerations
- **Chunk Size**: Optimal 200-500 tokens per chunk
- **Memory Limits**: Practical limit ~1M chunks per database
- **Concurrent Access**: Thread-safe operations via SQL connection pooling

## Data Consistency

### Transaction Management
- **Atomic Updates**: Memory, tags, and chunks updated in single transaction
- **Error Handling**: Rollback on any failure during upsert operations
- **Consistency**: Foreign key constraints ensure referential integrity

### Change Detection
```go
func (db *DB) GetMemory(name string) (*Memory, error)
```

File hash comparison workflow:
1. **Current Hash**: Calculate SHA-256 of file content
2. **Stored Hash**: Retrieve hash from database
3. **Compare**: Skip processing if hashes match
4. **Update**: Refresh database if content changed

### Cleanup Operations
```go
func (db *DB) DeleteMemory(name string) error
```

Cascading delete process:
1. **Chunks**: Remove all embedding chunks
2. **Tags**: Remove all tag associations  
3. **Memory**: Remove memory record
4. **Indexes**: Automatic index cleanup

The database layer provides a solid foundation for both traditional relational operations and modern vector similarity search, enabling SimpleMem's semantic approach to memory retrieval and analysis.