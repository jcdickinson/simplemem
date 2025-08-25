---
title: Debug Database Isolation
description: Database isolation strategy for testing MCP servers
tags:
  database: true
  debug: true
  isolation: true
  testing: true
created: 2025-08-24T16:27:04.419811943-07:00
modified: 2025-08-24T16:27:04.419811943-07:00
---

# Debug Database Isolation

Critical strategy for testing MCP servers with isolated databases.

## Why the --db Flag is Essential

### Avoids Database Locking
- Prevents "Conflicting lock" errors when main instance is running
- Allows testing while server is active

### Isolated Testing
- Each test can use its own database without interference
- No contamination between test scenarios

### Clean State
- Start with fresh database for reproducible test results
- Predictable test outcomes

### Parallel Testing
- Multiple test scenarios can run simultaneously
- Different tests can use different databases

### Debugging
- Preserve test databases for later inspection
- Compare database states across test runs

## Usage Pattern
```bash
# Always use custom database path for testing
./simplemem --db /tmp/test.db
```

## Debug Methodology

1. **Comprehensive Logging**: Add detailed logging at every pipeline step
2. **Direct Protocol Testing**: Use raw JSON-RPC calls for precise control
3. **Isolate Components**: Test each layer (API, database, embeddings) separately
4. **Real Data Testing**: Use actual data instead of synthetic examples
5. **Use Custom DB Paths**: Always use `--db` flag for testing to avoid conflicts

## Related Components
- [[debug-json-rpc-testing-pattern]] - Core testing pattern
- [[debug-justfile-integration]] - Justfile automation
- [[debug-testing-examples]] - Specific test examples