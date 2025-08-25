---
title: Database Query Optimization
description: Query optimization patterns and performance characteristics
tags:
  database: true
  optimization: true
  performance: true
created: 2025-08-24T16:21:35.193978997-07:00
modified: 2025-08-24T16:21:35.193978997-07:00
---

# Database Query Optimization

Query patterns and performance optimizations in SimpleMem's database layer.

## Optimized Query Patterns

### Semantic Search Query
Uses CTEs for efficient chunk similarity calculation:
1. Calculate similarities for all chunks
2. Find best chunk per memory
3. Join with memory metadata
4. Sort by similarity and limit

### Tag-Filtered Search
Combines semantic search with tag filtering:
- Subquery for semantic results
- IN clause for tag matching
- GROUP BY with HAVING for requireAll logic

## Performance Characteristics

### Memory Usage
- Embedding Storage: ~6KB per chunk (1536 floats × 4 bytes)
- Index Overhead: ~20% additional storage for vector indexes
- Connection Pool: Single connection with statement caching

### Query Performance
- Vector Search: Sub-second for up to 100K chunks
- Tag Filtering: Millisecond response with indexing
- Combined Queries: Optimized execution plans

### Scalability
- Chunk Size: Optimal 200-500 tokens per chunk
- Memory Limits: Practical limit ~1M chunks per database
- Concurrent Access: Thread-safe via SQL connection pooling

## Related Components
- [[database-vector-search]] - Vector search implementation
- [[simplemem-database-schema]] - Index strategy
- [[rag-performance-optimizations]] - RAG-level optimizations