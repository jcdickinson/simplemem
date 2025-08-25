---
title: List Memories Tool Temporarily Removed
description: Disabled list_memories tool to encourage semantic search usage and prevent excessive token consumption
tags:
  completed: true
  mcp: true
  optimization: true
  token-efficiency: true
  tool-management: true
created: 2025-08-24T22:26:53.627953905-07:00
modified: 2025-08-24T22:26:53.627953905-07:00
---

# List Memories Tool Temporarily Removed

## Decision
Temporarily removed the `list_memories` tool registration from the MCP server to encourage more efficient memory discovery patterns.

## Rationale
- **Token Efficiency**: `list_memories` can consume excessive tokens when displaying all memories with full frontmatter
- **Encourage Better Patterns**: Push users toward semantic search for more targeted, relevant results
- **Prevent Information Overload**: Large memory collections become unwieldy with full listing

## Implementation
- Commented out the tool registration in `internal/mcp/server.go`
- Handler function remains intact for easy re-enablement if needed
- No breaking changes to existing functionality

## Alternative Approaches
Users should rely on:
- `search_memories` - Semantic search for targeted discovery
- Direct memory access via `read_memory` when name is known
- Tag-based filtering through `search_memories` with tag filters

## Re-enablement
Can be easily restored by uncommenting the tool registration:
```go
mcpServer.AddTool(
    mcp.NewTool("list_memories", ...),
    s.handleListMemories,
)
```

This change promotes more intelligent memory discovery while maintaining system performance.