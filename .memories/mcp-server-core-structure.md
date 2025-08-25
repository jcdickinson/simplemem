---
title: MCP Server Core Structure
description: Core server type and component relationships in SimpleMem MCP implementation
tags:
  architecture: true
  golang: true
  mcp: true
  server: true
created: 2025-08-24T16:05:11.552709876-07:00
modified: 2025-08-24T16:05:11.552709876-07:00
---

# MCP Server Core Structure

The MCP server (`internal/mcp/server.go`) implements the Model Context Protocol to provide memory management tools for Claude and other MCP clients.

## Server Type Definition

```go
type Server struct {
    store         *memory.Store         // Basic file operations
    enhancedStore *memory.EnhancedStore // RAG-enabled operations
    mcpServer     *server.MCPServer     // MCP protocol handler
}
```

## Component Relationships

### Enhanced Store Wrapping
- Basic `Store` handles file I/O, markdown parsing, frontmatter
- `EnhancedStore` wraps basic store with database sync and RAG capabilities
- MCP server uses enhanced store for all operations to ensure consistency

### MCP Protocol Integration
- Uses `mark3labs/mcp-go` library for protocol implementation
- Handles JSON-RPC 2.0 communication with proper error management
- Supports tool capabilities and initial instructions

## Related Components
- [[mcp-tool-registration-pattern]] - Tool registration architecture
- [[mcp-request-processing-pipeline]] - Request flow and error handling
- [[simplemem-enhanced-store-core]] - Enhanced store implementation details