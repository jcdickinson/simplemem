---
title: SimpleMem Project Overview
description: Comprehensive overview of the SimpleMem MCP server project - a memory management system with RAG capabilities
tags:
  architecture: true
  golang: true
  mcp: true
  overview: true
created: 2025-08-24T00:18:44.504234473-07:00
modified: 2025-08-24T01:04:03.064737114-07:00
---

# SimpleMem Project Overview

SimpleMem is a Model Context Protocol (MCP) server that provides intelligent memory management with Retrieval-Augmented Generation (RAG) capabilities. It combines traditional file-based storage with semantic search powered by vector embeddings.

## Core Architecture

### MCP Server Foundation
Built on the Model Context Protocol, SimpleMem integrates seamlessly with Claude and other MCP-compatible clients. The server implements all required MCP capabilities:

- **Tool Registration**: Dynamic tool discovery and registration
- **Request Handling**: JSON-RPC 2.0 protocol implementation  
- **Error Management**: Structured error responses and recovery
- **Lifecycle Management**: Proper startup, shutdown, and signal handling

### Dual-Layer Storage Strategy

**Layer 1: File System (Primary)**
- Markdown files in `.memories/` directory
- YAML frontmatter for metadata (title, description, tags, timestamps)
- Markdown body content with wiki-style `[[links]]`
- Human-readable and git-friendly storage

**Layer 2: Database (Performance)**
- DuckDB with vector extension for similarity search
- Synchronized automatically with file changes
- Optimized for queries, aggregations, and vector operations
- Offline-capable with graceful degradation

## Key Features

### Advanced Search Capabilities
- **Semantic Search**: Vector-based similarity using VoyageAI embeddings
- **Tag Filtering**: Filter results by metadata tags
- **Backlink Discovery**: Find memories linking to specific documents

### RAG Features
- **Content Chunking**: Split large memories for better embeddings
- **Similarity Scoring**: Relevance ranking for search results
- **Context Awareness**: Use backlinks and tags to improve relevance
- **Incremental Processing**: Update embeddings when content changes

## MCP Tools

SimpleMem provides 7 core MCP tools:

1. **create_memory** - Create with auto frontmatter and timestamps
2. **read_memory** - Read with metadata summary and link extraction
3. **update_memory** - Update with timestamp preservation
4. **delete_memory** - Delete memory documents
5. **list_memories** - Enhanced listing with titles, tags, and dates
6. **search_memories** - Content search + tag search with semantic capabilities
7. **get_backlinks** - Find memories linking to a specific memory with semantic suggestions

## Data Storage

### File System
- **Primary Storage**: `.memories/` directory with markdown files
- **Cache**: `.cache/` directory for database and embeddings
- **Configuration**: `.mcp.json` for MCP server configuration

### Database Schema
- **memories**: Core document metadata and content
- **tags**: Normalized tag storage with key-value pairs
- **memory_links**: Inter-document relationships
- **embeddings**: Vector embeddings with chunking support

## Development Setup

### Dependencies
- Go 1.24+
- Just (build tool)
- Optional: Nix (for development environment)

### Key Commands
```bash
just run          # Start MCP server
just build        # Build binary  
just test         # Run tests
just fmt          # Format code
```

### Environment Variables
- `VOYAGE_API_KEY`: VoyageAI API key for embeddings
- Various configuration options in `internal/config/`

## Implementation Status

✅ **Phase 1 Complete**: Core MCP server with file-based storage
🚀 **Phase 2 Current**: RAG implementation with DuckDB and embeddings

The project successfully combines traditional file-based memory storage with modern RAG capabilities, providing both offline functionality and advanced semantic search when AI services are available.