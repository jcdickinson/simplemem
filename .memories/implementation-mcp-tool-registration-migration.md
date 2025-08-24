---
title: 'Implementation: MCP Tool Registration Migration Pattern'
tags:
  implementation: true
  mark3labs: true
  mcp: true
  migration: true
  tool-registration: true
created: 2025-08-24T00:47:15.429582012-07:00
modified: 2025-08-24T00:47:15.429582012-07:00
---

# MCP Tool Registration Migration Pattern

## Old Pattern (ThinkInAIXYZ/go-mcp)
```go
func (s *Server) registerTools(mcpServer *server.Server) error {
    // Create tool with struct definition
    createTool, err := protocol.NewTool(
        "create_memory",
        "Description here",
        createMemoryReq{}, // Struct defines parameters
    )
    if err != nil {
        return err
    }
    mcpServer.RegisterTool(createTool, s.handleCreateMemory)
    return nil
}
```

## New Pattern (mark3labs/mcp-go)  
```go
func (s *Server) registerTools(mcpServer *server.MCPServer) {
    // Direct tool definition with parameter builders
    mcpServer.AddTool(
        mcp.NewTool("create_memory",
            mcp.WithDescription("Description here"),
            mcp.WithString("name",
                mcp.Description("Parameter description"),
                mcp.Required(),
            ),
            mcp.WithString("content", 
                mcp.Description("Parameter description"),
                mcp.Required(),
            ),
        ),
        s.handleCreateMemory,
    )
    // No error handling needed - panics on invalid config
}
```

## Key Changes
1. **Function signature**: `*server.Server` → `*server.MCPServer`
2. **No error returns**: New API doesn't return errors from registration
3. **Parameter definition**: Struct-based → Builder pattern with `mcp.WithString()`, `mcp.WithObject()`, etc.
4. **Registration method**: `RegisterTool()` → `AddTool()`
5. **Tool creation**: `protocol.NewTool()` → `mcp.NewTool()`

## Parameter Types Available
- `mcp.WithString()` - String parameters
- `mcp.WithObject()` - Object/map parameters  
- `mcp.WithBoolean()` - Boolean parameters
- `mcp.WithNumber()` - Numeric parameters
- `mcp.Required()` - Mark parameter as required
- `mcp.Description()` - Add parameter description

## Status
✅ Successfully migrated all 7 tools:
- create_memory, read_memory, update_memory
- delete_memory, list_memories
- search_memories, get_backlinks

## Next Steps
Need to migrate handler function signatures from old to new API.

## Related
- [[mcp-library-migration-task]]
- [[knowledge-new-mcp-library-analysis]]