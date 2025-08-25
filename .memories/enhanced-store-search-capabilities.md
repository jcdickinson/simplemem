---
title: Enhanced Store Search Capabilities
description: Search implementations in SimpleMem Enhanced Store
tags:
  enhanced-store: true
  hybrid-search: true
  search: true
  semantic-search: true
created: 2025-08-24T16:16:59.010673706-07:00
modified: 2025-08-24T18:48:01.562710098-07:00
---

# Enhanced Store Search Capabilities

Advanced search features provided by the Enhanced Store.

## Semantic Search

### Basic Semantic Search
```go
func (es *EnhancedStore) SearchSemantic(query string, limit int)
```
Generates embeddings for search query, performs vector similarity search, returns memories with scores.

### Tag-Filtered Semantic Search
```go
func (es *EnhancedStore) SearchSemanticWithTags(query string, tagFilters map[string]string, requireAll bool, limit int)
```
- Apply tag filters before semantic search
- Support presence-only filters (empty value)
- Support exact value matching
- Boolean logic: ALL tags vs ANY tags

## Hybrid Search

### Algorithm Overview
1. **Keyword Search**: Traditional text matching
2. **Semantic Search**: Vector similarity
3. **Result Combination**: Merge and deduplicate
4. **Score Boosting**: Keyword matches get 1.2x boost
5. **Tag Filtering**: Apply filters to both sets
6. **Final Ranking**: Sort by combined relevance

### Scoring Strategy
- **Semantic Only**: Original similarity score (0.0-1.0)
- **Keyword Only**: Fixed high score (0.9)
- **Both**: Boosted similarity (similarity * 1.2, max 1.0)

## Backlink Analysis

### Enhanced Backlinks
```go  
func (es *EnhancedStore) GetEnhancedBacklinks(memoryName string, query string, limit int)
```
Combines explicit links with semantic similarity, optional query reranking, returns formatted markdown.

## Related Components
- [[simplemem-enhanced-store-core]] - Core architecture
- [[rag-semantic-search-operations]] - RAG search details
- [[mcp-discovery-tools]] - MCP tool interfaces