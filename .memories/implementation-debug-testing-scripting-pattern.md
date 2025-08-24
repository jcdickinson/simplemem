---
title: 'Implementation: Debug Testing & Scripting Pattern'
tags:
  debug: true
  implementation: true
  json-rpc: true
  mcp: true
  scripting: true
  testing: true
created: 2025-08-24T00:53:26.113392818-07:00
modified: 2025-08-24T01:13:14.660757595-07:00
---

# Debug Testing & Scripting Pattern from debug.md

## Key Scripting Pattern Discovered

From `debug.md`, the critical pattern for testing MCP servers is using **direct JSON-RPC calls via stdin**:

```bash
echo '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"search_memories","arguments":{"query":"medieval history"}},"id":1}' | ./simplemem
```

## Why This Pattern is Powerful

### 1. Direct MCP Protocol Testing
- Tests the actual MCP JSON-RPC 2.0 protocol
- Bypasses client libraries and tests raw server implementation
- Allows precise control over request parameters
- Immediate feedback on protocol compliance

### 2. Scriptable & Automatable  
- Can be used in shell scripts for automated testing
- Easy to parameterize for different test cases
- Can be combined with other command line tools
- Perfect for CI/CD testing pipelines

### 3. Debug-Friendly
- Raw request/response visible
- Can test edge cases easily
- Immediate error messages from server
- No client-side abstractions hiding issues

## JSON-RPC Request Structure for MCP Tools

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

## Testing Different Tools

### Search Memories
```bash
echo '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"search_memories","arguments":{"query":"search term"}},"id":1}' | ./simplemem
```

### Create Memory
```bash
echo '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"create_memory","arguments":{"name":"test-memory","content":"# Test\nContent here"}},"id":2}' | ./simplemem
```

### List Memories
```bash
echo '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"list_memories","arguments":{}},"id":3}' | ./simplemem
```

## Critical: Testing with Custom Database Path

**IMPORTANT**: Use the `--db` flag to avoid database locking conflicts when testing:

```bash
# Test basic functionality with custom DB path
echo '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"list_memories","arguments":{}},"id":1}' | ./simplemem --db /tmp/test.db

# Test semantic search with isolated database
echo '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"search_memories","arguments":{"query":"search term"}},"id":1}' | ./simplemem --db /tmp/semantic-test.db

# Test with specific test scenarios
echo '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"search_memories","arguments":{"query":"project overview"}},"id":1}' | ./simplemem --db /tmp/fix-test.db
```

### Why the --db Flag is Essential
- **Avoids Database Locking**: Prevents "Conflicting lock" errors when main instance is running
- **Isolated Testing**: Each test can use its own database without interference
- **Clean State**: Start with fresh database for reproducible test results
- **Parallel Testing**: Multiple test scenarios can run simultaneously
- **Debugging**: Preserve test databases for later inspection

## Debug Methodology from debug.md

1. **Comprehensive Logging**: Add detailed logging at every pipeline step
2. **Direct Protocol Testing**: Use raw JSON-RPC calls for precise control
3. **Isolate Components**: Test each layer (API, database, embeddings) separately
4. **Real Data Testing**: Use actual data instead of synthetic examples
5. **Use Custom DB Paths**: Always use `--db` flag for testing to avoid conflicts

## Application to Current Migration Testing

We can now test our migrated MCP server using:
```bash
# Test basic functionality
echo '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"list_memories","arguments":{}},"id":1}' | ./simplemem --db /tmp/migration-test.db

# Test with existing memories
echo '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"search_memories","arguments":{"query":"MCP migration"}},"id":2}' | ./simplemem --db /tmp/migration-test.db

# Test semantic search fix
echo '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"search_memories","arguments":{"query":"starting hints getting started"}},"id":3}' | ./simplemem --db /tmp/semantic-fix-test.db
```

## Related
- [[implementation-mcp-library-migration-task]]
- [[todo-fix-semantic-backlinks-sql-constraint]] - Outstanding SQL issue
- File: `debug.md` (debugging insights)
- Pattern: Direct JSON-RPC testing for MCP protocol validation

---
📝 **Title:** Implementation: Debug Testing & Scripting Pattern
🏷️ **Tags:** debug, implementation, json-rpc, mcp, scripting, testing
📅 **Created:** 2025-08-24 00:53:26
🔄 **Modified:** 2025-08-24 01:11:30

🔗 **Links found in this memory:**
- [implementation-mcp-library-migration-task](implementation-mcp-library-migration-task.md) (wiki link)
- [todo-fix-semantic-backlinks-sql-constraint](todo-fix-semantic-backlinks-sql-constraint.md) (wiki link)