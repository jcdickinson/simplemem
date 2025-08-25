---
title: Database Vector Search
description: Vector search implementation using DuckDB's vector extension
tags:
  database: true
  duckdb: true
  search: true
  vectors: true
created: 2025-08-24T16:20:06.87597864-07:00
modified: 2025-08-24T16:20:06.87597864-07:00
---

# Database Vector Search

DuckDB's vector extension enables efficient semantic search in SimpleMem.

## Vector Similarity Search

```sql
SELECT m.*, c.content, 1 - array_cosine_similarity(c.embedding, $1) as similarity
FROM memories m 
JOIN chunks c ON m.id = c.memory_id 
ORDER BY similarity DESC 
LIMIT ?
```

## Search Process

1. **Vector Query**: Find chunks with embeddings similar to query
2. **Tag Filtering**: Apply metadata filters using SQL WHERE clauses
3. **Aggregation**: Group results by memory, selecting best chunk
4. **Ranking**: Sort by similarity score (cosine distance)
5. **Limit**: Return top N results

## Tag Filter Logic

```go
type TagFilter struct {
    Key        string
    Value      string  
    CheckValue bool    // If false, only check key presence
}
```

- Presence Check: `tags.key = ? AND tags.value IS NOT NULL`
- Value Check: `tags.key = ? AND tags.value = ?`
- Boolean Logic: AND/OR combinations via `requireAll` parameter

## Embedding Storage
- Vector Format: FLOAT[1536] array for VoyageAI embeddings
- Indexing: Automatic vector similarity indexes
- Compression: DuckDB's built-in vector compression

## Performance
- Vector Search: Sub-second for up to 100K chunks
- Tag Filtering: Millisecond response with indexing
- Combined Queries: Optimized execution plans

## Related Components
- [[simplemem-database-schema]] - Database schema design
- [[rag-semantic-search-operations]] - RAG search interface
- [[database-query-optimization]] - Query optimization strategies