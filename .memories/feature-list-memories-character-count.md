---
title: 'Feature: List Memories Character Count'
description: Implementation of character count display in list_memories MCP tool
tags:
  completed: true
  feature: true
  implementation: true
  list-memories: true
  mcp: true
created: 2025-08-24T15:36:00-07:00
modified: 2025-08-24T15:37:50.607538155-07:00
---

# Feature: List Memories Character Count

## Status: ✅ COMPLETED

Successfully implemented character count display in the `list_memories` MCP tool response.

## Implementation Details

### Changes Made
**File**: `internal/mcp/server.go`
**Function**: `handleListMemories`

Added character count calculation and display:
```go
// Calculate content length
contentLength := len(memInfo.Content)

// ... existing formatting code ...

result += fmt.Sprintf(" (%d chars)", contentLength)
```

### Output Format
The character count now appears after the tags and before the modification date:
```
📄 **memory-name** - Title 🏷️[tags] (1234 chars) (modified: 2025-08-24)
```

## Testing Results

Tested with JSON-RPC call:
```bash
echo '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"list_memories","arguments":{}},"id":1}' | 
./simplemem --db /tmp/test.db
```

**Sample output**:
- **simplemem-database-layer** (7261 chars)
- **implementation-debug-testing-scripting-pattern** (5344 chars)  
- **test-ai-memory** (225 chars)

## Benefits
1. **Size Awareness**: Users can quickly see memory content size
2. **Storage Planning**: Helps understand memory storage usage
3. **Content Categorization**: Large vs small memories are easily distinguishable
4. **Debugging**: Useful for troubleshooting content issues

## User Experience
- Character counts provide immediate context about memory size
- Format is concise and doesn't clutter the display
- Consistent placement in the metadata line
- Works for all memories regardless of content type

The feature is production-ready and working as requested.