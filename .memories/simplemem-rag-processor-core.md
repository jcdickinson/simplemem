---
title: SimpleMem RAG Processor Core
description: Core RAG processor architecture and components
tags:
  architecture: true
  golang: true
  rag: true
created: 2025-08-24T16:11:24.46969765-07:00
modified: 2025-08-24T16:11:24.46969765-07:00
---

# SimpleMem RAG Processor Core

The RAG Processor (`internal/rag/processor.go`) is the intelligence layer that handles embedding generation, semantic search, and advanced retrieval operations using vector similarity.

## Core Processor Type

```go
type Processor struct {
    db         *db.DB              // Database layer for storage
    embeddings *embeddings.Client  // Embedding service client
}
```

## Component Responsibilities

### RAG Processor Role
- **Embedding Generation**: Convert text to vector representations
- **Semantic Search**: Find similar content using vector operations
- **Content Chunking**: Split large documents for better embeddings
- **Similarity Ranking**: Score and rank search results
- **Backlink Analysis**: Combine explicit and semantic relationships

## Configuration and Validation

### Configuration Management
```go
func (p *Processor) ValidateConfiguration() error
```

- **API Key Validation**: Test embedding service connectivity
- **Model Verification**: Confirm embedding model availability
- **Database Schema**: Verify vector extension support
- **Performance Limits**: Check system resource constraints

### Processing Options
- **Batch Size**: Configure embedding request batching
- **Similarity Thresholds**: Set minimum similarity scores
- **Chunk Sizes**: Optimize content chunking parameters
- **Cache Settings**: Configure embedding result caching

## Related Components
- [[rag-voyageai-integration]] - VoyageAI embedding service integration
- [[rag-semantic-search-operations]] - Search implementation details
- [[rag-content-processing-pipeline]] - Processing workflow