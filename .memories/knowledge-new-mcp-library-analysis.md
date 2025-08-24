---
title: 'Knowledge: New MCP Library Analysis (mark3labs/mcp-go)'
tags:
  area: implementation
  component: mcp
  knowledge: true
  library: mark3labs-mcp-go
  migration: true
created: 2025-08-24T00:42:45.539118196-07:00
modified: 2025-08-24T00:42:45.539118196-07:00
---

# New MCP Library Structure (mark3labs/mcp-go v0.38.0)

## Key Packages
- `github.com/mark3labs/mcp-go/server` - Server functionality
- `github.com/mark3labs/mcp-go/mcp` - Core MCP types and protocols
- `github.com/mark3labs/mcp-go/client` - Client functionality

## Major API Differences from ThinkInAIXYZ/go-mcp

### Server Creation
```go
// NEW API
server := server.NewMCPServer(name, version string, opts ...ServerOption)

// vs OLD API  
server := server.NewServer(transport, server.WithServerInfo(...))
```

### Key Features Available
- **Initial Instructions Support**: `WithInstructions(instructions string)` option!
- **Multiple Transport Types**: Stdio, HTTP, SSE
- **Enhanced Capabilities**: Tools, Resources, Prompts with capability flags
- **Middleware Support**: `WithToolHandlerMiddleware`
- **Hooks System**: Before/After request hooks
- **Recovery & Logging**: Built-in options

### Tool Registration Pattern
Based on docs, likely uses method-based registration:
```go
server.AddTool("tool_name", handler_func)
```

### Handler Function Signature
```go
type ToolHandlerFunc func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
```

### Transport Handling
```go
// Stdio transport
server.ServeStdio(mcpServer, opts...)
```

## Migration Strategy
1. Update server creation to use `NewMCPServer`
2. Add `WithInstructions` option for initial instructions
3. Convert tool registration from `RegisterTool` to new pattern
4. Update handler signatures if needed
5. Replace transport initialization
6. Test all functionality

## Initial Instructions Implementation
The `WithInstructions(instructions string)` option allows setting initial instructions that will be sent to clients on initialization - exactly what we need!

## Related
- [[mcp-library-migration-task]]
- [[knowledge-current-mcp-implementation]]
- [[simplemem-mcp-server-architecture]]