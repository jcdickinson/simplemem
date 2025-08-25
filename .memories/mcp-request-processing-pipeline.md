---
title: MCP Request Processing Pipeline
description: Request flow and error handling in SimpleMem MCP server
tags:
  architecture: true
  error-handling: true
  mcp: true
  pipeline: true
created: 2025-08-24T16:09:22.887621468-07:00
modified: 2025-08-24T16:09:22.887621468-07:00
---

# MCP Request Processing Pipeline

How SimpleMem processes MCP requests from receipt to response.

## Standard Request Flow

1. **JSON-RPC Parsing**: MCP library handles protocol details
2. **Parameter Extraction**: Type-safe parameter parsing with defaults
3. **Enhanced Store Operation**: All operations go through enhanced store
4. **Database Sync**: Automatic sync for write operations
5. **RAG Processing**: Embedding generation/updates as needed
6. **Response Assembly**: Consistent text content response format

## Error Handling Strategy

### Validation Errors
- Parameter validation before processing
- Required field checking via MCP library

### Store Errors 
- File system errors (permissions, disk space, etc.)
- Database errors (connection, constraint violations, etc.)
- RAG processing errors (API failures, embedding issues, etc.)

### Recovery Patterns
- Graceful degradation (semantic search falls back to basic search)
- Partial success handling (warn on database sync failure, succeed on file operation)
- User-friendly error messages with context

## Performance Characteristics

### Request Latency
- **File Operations**: ~1-5ms for typical memory sizes
- **Database Sync**: ~5-10ms with indexes
- **Semantic Search**: ~50-200ms depending on collection size  
- **Embedding Generation**: ~100-500ms via VoyageAI API

### Memory Usage
- **Server Process**: ~10-50MB baseline
- **Database**: ~6KB per memory + embeddings
- **Cache**: Configurable embedding cache size

## Related Components
- [[mcp-server-core-structure]] - Server architecture
- [[mcp-server-initialization]] - Server startup and configuration
- [[simplemem-enhanced-store-core]] - Store implementation