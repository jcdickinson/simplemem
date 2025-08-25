---
title: MCP Memory Management Tools
description: Core CRUD tools for memory operations in SimpleMem MCP server
tags:
  crud: true
  mcp: true
  memory-management: true
  tools: true
created: 2025-08-24T16:08:20.142109115-07:00
modified: 2025-08-24T16:08:20.142109115-07:00
---

# MCP Memory Management Tools

Core tools that handle CRUD operations for memories in SimpleMem.

## create_memory
- **Purpose**: Create new memory with automatic metadata
- **Process**: Enhanced store handles file creation + database sync + RAG processing
- **Features**: Auto-timestamps, frontmatter generation, embedding creation

## read_memory  
- **Purpose**: Read memory with rich metadata presentation
- **Features**: Metadata summary, link extraction, formatting with emojis
- **Enhancement**: Could include related memory suggestions

## update_memory
- **Purpose**: Update existing memory preserving metadata  
- **Process**: Enhanced store handles file update + database sync + RAG reprocessing
- **Features**: Timestamp preservation, change detection via file hash

## delete_memory
- **Purpose**: Remove memory from both filesystem and database
- **Process**: Cascading delete (file → database → embeddings → indexes)
- **Safety**: Error handling ensures partial failures don't corrupt state

## Related Components
- [[mcp-discovery-tools]] - Search and listing capabilities
- [[mcp-tool-registration-pattern]] - How tools are registered
- [[simplemem-enhanced-store-core]] - Underlying store operations