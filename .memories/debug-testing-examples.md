---
title: Debug Testing Examples
description: Specific examples of testing different MCP tools
tags:
  debug: true
  examples: true
  mcp: true
  testing: true
created: 2025-08-24T16:27:36.137729226-07:00
modified: 2025-08-24T16:27:36.137729226-07:00
---

# Debug Testing Examples

Specific examples for testing different MCP tools.

## Search Memories
```bash
just test-json '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"search_memories","arguments":{"query":"search term"}},"id":1}' /tmp/test.db
```

## Create Memory
```bash
just test-json '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"create_memory","arguments":{"name":"test-memory","content":"# Test\nContent here"}},"id":2}' /tmp/test.db
```

## List Memories
```bash
just test-json '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"list_memories","arguments":{}},"id":3}' /tmp/test.db
```

## Migration Testing

### Basic Functionality Test
```bash
just test-json '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"list_memories","arguments":{}},"id":1}' /tmp/migration-test.db
```

### Semantic Search Testing
```bash
just test-json '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"search_memories","arguments":{"query":"MCP migration"}},"id":2}' /tmp/migration-test.db
```

### Semantic Backlinks Fix Testing
```bash
# Use the comprehensive backlinks test
just test-backlinks /tmp/semantic-fix-test.db
```

## Related Components
- [[debug-json-rpc-testing-pattern]] - Core testing pattern
- [[debug-justfile-integration]] - Justfile automation
- [[debug-database-isolation]] - Database isolation strategy
- File: `debug.md` - Debugging insights
- File: `justfile` - Test automation