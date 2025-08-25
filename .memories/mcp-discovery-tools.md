---
title: MCP Discovery Tools
description: Search, listing, and discovery tools in SimpleMem MCP server
tags:
  discovery: true
  mcp: true
  search: true
  tools: true
created: 2025-08-24T16:08:51.094912051-07:00
modified: 2025-08-24T16:08:51.094912051-07:00
---

# MCP Discovery Tools

Tools for finding and exploring memories in SimpleMem.

## list_memories
- **Purpose**: Enhanced memory listing with metadata preview
- **Features**: Title display, tag summary, modification dates
- **Format**: Rich text with emojis and consistent formatting
- **Performance**: Uses basic store for file listing, enhanced for metadata

## search_memories
- **Purpose**: Semantic search with optional tag filtering
- **Algorithm**: Pure semantic search using vector embeddings
- **Features**: Tag filters (presence and value matching), boolean logic
- **Response**: Markdown formatted results with snippets and scores
- **Limit**: Fixed at 5 results for optimal relevance

## get_backlinks  
- **Purpose**: Find memories linking to a specific memory
- **Features**: Explicit link detection + semantic similarity
- **Enhancement**: Optional query-based reranking for context-specific relevance
- **Format**: Markdown with context snippets and relationship types

## Implementation Details

### Semantic Search Processing
```go
// Use semantic search with tag filtering (set to 5 docs as requested)
result, err := s.enhancedStore.SearchSemanticMarkdownWithTags(query, tags, requireAll, 5)
```

### Backlink Enhancement
```go
// Get enhanced backlinks with reranking (set to 5 docs as requested)  
result, err := s.enhancedStore.GetEnhancedBacklinks(name, query, 5)
```

## Related Components
- [[mcp-memory-management-tools]] - CRUD operations
- [[simplemem-rag-processor-core]] - RAG processing details
- [[mcp-request-processing-pipeline]] - Request handling