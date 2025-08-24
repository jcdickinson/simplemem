---
title: SimpleMem RAG Processor
description: Analysis of the RAG (Retrieval-Augmented Generation) processor that handles embeddings, semantic search, and similarity operations
tags:
  architecture: true
  embeddings: true
  golang: true
  rag: true
  semantic-search: true
  vector-processing: true
created: 2025-08-24T00:23:01.947606246-07:00
modified: 2025-08-24T01:01:42.342315808-07:00
---

# SimpleMem RAG Processor

The RAG Processor (`internal/rag/processor.go`) is the intelligence layer that handles embedding generation, semantic search, and advanced retrieval operations using vector similarity.

## Processor Architecture

### Core Processor Type

```go
type Processor struct {
    db         *db.DB              // Database layer for storage
    embeddings *embeddings.Client  // Embedding service client
}
```

### Component Responsibilities

#### RAG Processor Role
- **Embedding Generation**: Convert text to vector representations
- **Semantic Search**: Find similar content using vector operations
- **Content Chunking**: Split large documents for better embeddings
- **Similarity Ranking**: Score and rank search results
- **Backlink Analysis**: Combine explicit and semantic relationships

## Embedding Integration

### VoyageAI Integration

The processor integrates with VoyageAI for embedding generation:

#### Client Configuration
```go
func NewProcessor(database *db.DB, cfg *config.Config) (*Processor, error)
```
- **API Key Management**: Secure credential handling via configuration
- **Model Selection**: Support for different VoyageAI embedding models
- **Request Batching**: Efficient batch processing of multiple texts
- **Error Handling**: Graceful degradation when service unavailable

#### Embedding Models
- **Default Model**: `voyage-large-2` for high-quality embeddings
- **Dimension Size**: 1536 dimensions for vector storage
- **Context Length**: Support for large document chunks
- **Language Support**: Multi-language embedding capabilities

### Embedding Operations

#### Text Processing
```go
func (p *Processor) ProcessMemory(memory *db.Memory) error
```

1. **Content Chunking**: Split content into optimal sizes for embeddings
2. **Chunk Processing**: Generate embeddings for each chunk
3. **Storage**: Store embeddings with chunk metadata in database
4. **Indexing**: Update vector indexes for fast similarity search

#### Chunking Strategy
- **Semantic Chunking**: Split at natural boundaries (paragraphs, sections)
- **Overlap Handling**: Maintain context between chunks
- **Size Optimization**: Balance detail vs. performance
- **Metadata Preservation**: Track chunk position and context

## Search Operations

### Semantic Search

#### Basic Similarity Search
```go
func (p *Processor) SearchSimilarMemories(query string, limit int) ([]db.Memory, []float32, error)
```

1. **Query Embedding**: Generate embedding for search query
2. **Vector Search**: Find similar embeddings in database
3. **Ranking**: Sort results by cosine similarity score
4. **Aggregation**: Group chunks by memory and select best matches
5. **Result Assembly**: Return complete memory objects with scores

#### Tag-Filtered Search
```go
func (p *Processor) SearchSimilarMemoriesWithTags(query string, tagFilters []db.TagFilter, requireAll bool, limit int) ([]db.Memory, []float32, error)
```

- **Pre-filtering**: Apply tag filters before vector search
- **Performance**: Reduce vector operations by filtering candidates
- **Boolean Logic**: Support AND/OR combinations of tag conditions
- **Efficiency**: Use database indexes for tag filtering

### Backlink Analysis

#### Enhanced Backlinks
```go
func (p *Processor) GetEnhancedBacklinks(memoryName string, query string, limit int) ([]BacklinkResult, error)
```

The processor provides sophisticated backlink analysis:

1. **Explicit Links**: Find direct markdown and wiki-style links
2. **Semantic Links**: Find semantically similar memories
3. **Query Reranking**: Optionally rerank by relevance to query
4. **Unified Results**: Combine different link types with scores
5. **Context Extraction**: Provide snippets showing link context

#### Semantic Backlinks
```go  
func (p *Processor) GetSemanticBacklinks(name string, minSimilarity float32) ([]db.Memory, []float32, error)
```

- **Memory Embedding**: Use existing memory embedding as query
- **Similarity Threshold**: Filter results below minimum threshold
- **Self-Exclusion**: Remove the source memory from results
- **Bidirectional**: Find memories similar to the target memory

## Advanced RAG Features

### Content Processing Pipeline

#### Memory Processing Workflow
1. **Content Analysis**: Extract text content from markdown
2. **Chunk Generation**: Split content into processable segments  
3. **Embedding Generation**: Convert chunks to vector representations
4. **Link Extraction**: Identify references to other memories
5. **Database Storage**: Store all derived data with consistency
6. **Index Updates**: Maintain search indexes for performance

#### Error Recovery
- **Partial Processing**: Continue if some chunks fail
- **Retry Logic**: Attempt failed operations with backoff
- **Graceful Degradation**: Provide text search if embeddings fail
- **Status Tracking**: Monitor processing status for debugging

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

## Performance Optimizations

### Caching Strategies
- **Embedding Cache**: Avoid regenerating unchanged content embeddings
- **Query Cache**: Cache frequent similarity search results
- **Model Cache**: Cache embedding model responses
- **Database Cache**: Cache frequent database query results

### Batch Processing
- **Initial Setup**: Efficiently process existing memories
- **Bulk Operations**: Handle multiple memories in single requests
- **Queue Management**: Process updates asynchronously when possible
- **Resource Management**: Manage memory and API usage

### Monitoring and Metrics
- **Processing Times**: Track embedding generation performance
- **API Usage**: Monitor embedding service quota usage
- **Cache Hit Rates**: Optimize caching effectiveness
- **Search Performance**: Profile query execution times

The RAG Processor provides the intelligent layer that transforms SimpleMem from a simple file storage system into a sophisticated semantic memory system capable of understanding content relationships and providing relevant, context-aware search results.

---
📝 **Title:** SimpleMem RAG Processor
📄 **Description:** Analysis of the RAG (Retrieval-Augmented Generation) processor that handles embeddings, semantic search, and similarity operations
🏷️ **Tags:** semantic-search, vector-processing, architecture, embeddings, golang, rag
📅 **Created:** 2025-08-24 00:23:01
🔄 **Modified:** 2025-08-24 00:23:01