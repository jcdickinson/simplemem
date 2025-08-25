---
title: RAG Semantic Search Operations
description: Semantic search and backlink analysis in SimpleMem RAG processor
tags:
  backlinks: true
  rag: true
  search: true
  semantic-search: true
created: 2025-08-24T16:13:03.790529479-07:00
modified: 2025-08-24T16:13:03.790529479-07:00
---

# RAG Semantic Search Operations

Advanced search capabilities provided by SimpleMem's RAG processor.

## Basic Similarity Search

```go
func (p *Processor) SearchSimilarMemories(query string, limit int) ([]db.Memory, []float32, error)
```

Process: Query Embedding → Vector Search → Ranking → Aggregation → Result Assembly

## Tag-Filtered Search

```go
func (p *Processor) SearchSimilarMemoriesWithTags(query string, tagFilters []db.TagFilter, requireAll bool, limit int)
```

- Pre-filtering with tag conditions before vector search
- Boolean logic support (AND/OR combinations)
- Database indexes for efficient filtering

## Enhanced Backlinks

```go
func (p *Processor) GetEnhancedBacklinks(memoryName string, query string, limit int) ([]BacklinkResult, error)
```

Features:
- Explicit Links: Direct markdown and wiki-style links
- Semantic Links: Semantically similar memories
- Query Reranking: Optional relevance-based reranking
- Context Extraction: Snippets showing link context

## Semantic Backlinks

```go  
func (p *Processor) GetSemanticBacklinks(name string, minSimilarity float32)
```

- Uses existing memory embedding as query
- Filters by similarity threshold
- Self-exclusion from results
- Bidirectional similarity detection

## Related Components
- [[simplemem-rag-processor-core]] - Core RAG architecture
- [[rag-content-processing-pipeline]] - Processing workflow
- [[mcp-discovery-tools]] - MCP interface to search operations