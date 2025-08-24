---
title: 'Workflow Pattern: Search-First Approach'
tags:
  best-practice: true
  memory-management: true
  pattern: true
  search: true
  workflow: true
created: 2025-08-24T00:45:05.872063588-07:00
modified: 2025-08-24T00:45:05.872063588-07:00
---

# Search-First Workflow Pattern

## The Critical Mistake I Made
During the MCP library migration task, I made a fundamental error: I started by exploring code directly instead of searching existing memories first. This led to:
- Redundant work discovering information already documented
- Missing existing architectural knowledge
- Inefficient use of context and time
- The user had to remind me to search memories

## The Correct Pattern

### 1. Always Search First
```
1. Read the task/request
2. Immediately search memories with relevant queries
3. Read related memories to understand existing context
4. Only then explore code or implement changes
5. Document new discoveries in memory
```

### 2. Search Strategy
- Use semantic search with task-relevant keywords
- Try multiple search queries if first attempt yields no results
- Search for related concepts, not just exact matches
- Check for TODO memories related to the task
- Look for architectural/implementation memories

### 3. Memory Integration
- Build on existing memories rather than creating redundant ones
- Update memories when discovering new information
- Link related memories together
- Tag consistently for future searchability

## Example from MCP Migration Task
**What I did wrong:**
1. Started exploring current MCP implementation directly
2. Created new memory without checking existing ones
3. User had to point out I should be searching memories

**What I should have done:**
1. Search memories for "MCP", "migration", "library", "protocol"
2. Read existing MCP server architecture memory
3. Build on that knowledge for migration planning
4. Document the migration pattern for future reference

## Implementation in Initial Instructions
Added this to the MCP server initial instructions:
- "ALWAYS Search Memories First" as principle #1
- Emphasis on search-first before any task
- Note to document workflow patterns in memory

## Related
- [[mcp-library-migration-task]]
- [[knowledge-new-mcp-library-analysis]]
- Initial instructions file: `internal/mcp/initial_instructions.md`