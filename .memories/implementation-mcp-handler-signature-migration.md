---
title: 'Implementation: MCP Handler Signature Migration'
tags:
  completed: true
  function-signatures: true
  handlers: true
  implementation: true
  mcp: true
  migration: true
created: 2025-08-24T00:47:29.812757717-07:00
modified: 2025-08-24T00:49:41.741957467-07:00
---

# MCP Handler Function Signature Migration

## Status: ✅ COMPLETED
All handler functions successfully migrated from old to new API.

## Old Signature (ThinkInAIXYZ/go-mcp)
```go
func (s *Server) handleCreateMemory(_ context.Context, request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
    req := new(createMemoryReq)
    if err := protocol.VerifyAndUnmarshal(request.RawArguments, &req); err != nil {
        return nil, err
    }
    
    // Use req.Name, req.Content
    
    return &protocol.CallToolResult{
        Content: []protocol.Content{
            &protocol.TextContent{
                Type: "text",
                Text: "result",
            },
        },
    }, nil
}
```

## New Signature (mark3labs/mcp-go)
```go
func (s *Server) handleCreateMemory(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    // Direct parameter access from request
    name := request.GetString("name", "")
    content := request.GetString("content", "")
    
    return &mcp.CallToolResult{
        Content: []mcp.Content{
            mcp.TextContent{
                Type: "text", 
                Text: "result",
            },
        },
    }, nil
}
```

## Key Changes Applied
1. **Parameter type**: `*protocol.CallToolRequest` → `mcp.CallToolRequest` ✅
2. **No pointer**: New request is value, not pointer ✅
3. **Parameter access**: `VerifyAndUnmarshal()` → `GetString()`, `GetArguments()` ✅
4. **Return types**: `protocol.*` → `mcp.*` ✅
5. **Content types**: Pointer types → Value types ✅

## Migration Completed
- ✅ handleCreateMemory - Simple string params
- ✅ handleReadMemory - Single string param
- ✅ handleUpdateMemory - Two string params  
- ✅ handleDeleteMemory - Single string param
- ✅ handleListMemories - No params
- ✅ handleSearchMemories - Complex with tags/boolean
- ✅ handleGetBacklinks - Two string params

## Special Handling
For `handleSearchMemories`: Had to handle complex object parameter (tags) using `GetArguments()` and type assertions since tags is a map[string]string.

## Cleanup
- ✅ Removed all old request struct types
- ✅ Updated transport handling (Run method)
- ✅ Fixed Shutdown method for new API
- ✅ Build successful

## Related
- [[implementation-mcp-tool-registration-migration]]  
- [[mcp-library-migration-task]]