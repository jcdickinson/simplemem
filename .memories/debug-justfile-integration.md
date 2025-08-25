---
title: Debug Justfile Integration
description: Justfile commands for automated MCP testing
tags:
  automation: true
  debug: true
  justfile: true
  testing: true
created: 2025-08-24T16:26:35.915180331-07:00
modified: 2025-08-24T16:26:35.915180331-07:00
---

# Debug Justfile Integration

Justfile integration for automated MCP testing workflows.

## Universal Test Command
```bash
# Test any JSON-RPC call with optional custom database
just test-json '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"list_memories","arguments":{}},"id":1}' /tmp/test.db
```

## Specialized Test Commands

### Semantic Backlinks Testing
```bash
# Comprehensive test that creates multiple memories to trigger semantic linking
just test-backlinks /tmp/backlinks-test.db
```

### Quick Clean Tests
```bash
# Quick test with fresh database
just test-clean /tmp/clean-test.db
```

## Justfile Advantages

### Parameterized Testing
- Easy to change database paths
- Reusable test patterns
- Clean command syntax

### Automated Workflows
- `test-backlinks` handles complete workflow
- Automatic database cleanup
- Consistent test environment

### Documentation as Code
- Test commands serve as documentation
- Easy to discover available tests with `just --list`
- Self-documenting test patterns

## Related Components
- [[debug-json-rpc-testing-pattern]] - Core testing pattern
- [[debug-database-isolation]] - Database isolation strategy
- [[debug-testing-examples]] - Specific test examples