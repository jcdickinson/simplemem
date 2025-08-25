---
title: MCP Tool Registration Pattern
description: Consistent pattern for registering MCP tools in SimpleMem server
tags:
  architecture: true
  mcp: true
  patterns: true
  tools: true
created: 2025-08-24T16:07:56.970506481-07:00
modified: 2025-08-24T16:07:56.970506481-07:00
---

# MCP Tool Registration Pattern

Each MCP tool in SimpleMem follows a consistent registration pattern for clean architecture.

## Registration Pattern

```go
mcpServer.AddTool(
    mcp.NewTool("tool_name",
        mcp.WithDescription("Tool description"),
        mcp.WithString("param", mcp.Description("Parameter description"), mcp.Required()),
    ),
    s.handleToolName,
)
```

## Tool Implementation Strategy

### Parameter Extraction
```go
name := request.GetString("name", "")
content := request.GetString("content", "")
```

### Enhanced Store Operations
All tools use the enhanced store to ensure:
- Database synchronization
- RAG processing
- Consistent metadata handling

### Response Formatting
```go
return &mcp.CallToolResult{
    Content: []mcp.Content{
        mcp.TextContent{Type: "text", Text: result},
    },
}, nil
```

## Related Components
- [[mcp-server-core-structure]] - Server architecture overview
- [[mcp-memory-management-tools]] - Specific tool implementations
- [[mcp-discovery-tools]] - Search and listing tools