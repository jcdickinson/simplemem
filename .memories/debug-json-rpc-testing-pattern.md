---
title: Debug JSON-RPC Testing Pattern
description: Direct JSON-RPC testing pattern for MCP servers
tags:
  debug: true
  json-rpc: true
  mcp: true
  testing: true
created: 2025-08-24T16:26:07.839321314-07:00
modified: 2025-08-24T16:26:07.839321314-07:00
---

# Debug JSON-RPC Testing Pattern

Critical pattern for testing MCP servers using direct JSON-RPC calls via stdin.

## Basic Pattern
```bash
echo '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"search_memories","arguments":{"query":"medieval history"}},"id":1}' | ./simplemem
```

## JSON-RPC Request Structure
```json
{
  "jsonrpc": "2.0",
  "method": "tools/call", 
  "params": {
    "name": "TOOL_NAME",
    "arguments": {
      "param1": "value1",
      "param2": "value2"
    }
  },
  "id": 1
}
```

## Why This Pattern is Powerful

### Direct MCP Protocol Testing
- Tests actual MCP JSON-RPC 2.0 protocol
- Bypasses client libraries
- Allows precise control over request parameters
- Immediate feedback on protocol compliance

### Scriptable & Automatable  
- Can be used in shell scripts for automated testing
- Easy to parameterize for different test cases
- Can be combined with other command line tools
- Perfect for CI/CD testing pipelines

### Debug-Friendly
- Raw request/response visible
- Can test edge cases easily
- Immediate error messages from server
- No client-side abstractions hiding issues

## Related Components
- [[debug-justfile-integration]] - Justfile automation
- [[debug-database-isolation]] - Database isolation for testing
- [[debug-testing-examples]] - Specific test examples