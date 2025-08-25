---
title: RAG VoyageAI Integration
description: VoyageAI embedding service integration in SimpleMem
tags:
  embeddings: true
  integration: true
  rag: true
  voyageai: true
created: 2025-08-24T16:11:53.980978184-07:00
modified: 2025-08-24T16:11:53.980978184-07:00
---

# RAG VoyageAI Integration

SimpleMem integrates with VoyageAI for high-quality embedding generation.

## Client Configuration

```go
func NewProcessor(database *db.DB, cfg *config.Config) (*Processor, error)
```

- **API Key Management**: Secure credential handling via configuration
- **Model Selection**: Support for different VoyageAI embedding models
- **Request Batching**: Efficient batch processing of multiple texts
- **Error Handling**: Graceful degradation when service unavailable

## Embedding Models

- **Default Model**: `voyage-large-2` for high-quality embeddings
- **Dimension Size**: 1536 dimensions for vector storage
- **Context Length**: Support for large document chunks
- **Language Support**: Multi-language embedding capabilities

## Embedding Operations

### Text Processing
```go
func (p *Processor) ProcessMemory(memory *db.Memory) error
```

1. **Content Chunking**: Split content into optimal sizes for embeddings
2. **Chunk Processing**: Generate embeddings for each chunk
3. **Storage**: Store embeddings with chunk metadata in database
4. **Indexing**: Update vector indexes for fast similarity search

### Chunking Strategy
- **Semantic Chunking**: Split at natural boundaries (paragraphs, sections)
- **Overlap Handling**: Maintain context between chunks
- **Size Optimization**: Balance detail vs. performance
- **Metadata Preservation**: Track chunk position and context

## Related Components
- [[simplemem-rag-processor-core]] - Core RAG processor architecture
- [[rag-performance-optimizations]] - Caching and performance strategies
- [[rag-content-processing-pipeline]] - Full processing workflow