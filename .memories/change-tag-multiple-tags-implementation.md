---
title: Change Tag Multiple Tags Implementation
description: Implementation of enhanced change_tag tool supporting multiple tag operations
tags:
  change_tag: true
  completed: true
  implementation: true
  mcp: true
  multiple-tags: true
created: 2025-08-24T01:36:25.466758471-07:00
modified: 2025-08-24T01:36:25.466758471-07:00
---

# Change Tag Multiple Tags Implementation

## Status: ✅ COMPLETED

Successfully enhanced the `change_tag` MCP tool to support setting multiple tags in a single operation.

## Implementation Details

### API Changes
- Changed from individual `tag_key`/`tag_value` parameters to a `tags` object parameter
- Supports setting multiple tags atomically: `{"todo": true, "status": "in_progress", "priority": "high"}`
- Use `null` values to remove tags: `{"old_tag": null}`

### Files Modified
1. **internal/mcp/server.go**:
   - Updated tool registration to use `mcp.WithObject("tags")` parameter
   - Rewrote `handleChangeTag` method to process multiple tags
   - Enhanced response messages to show all changes

2. **internal/mcp/initial_instructions.md**:
   - Updated TODO tracking section with new API examples
   - Added multiple tag usage patterns

### Key Features
- **Atomic Updates**: All tag changes happen in a single operation
- **Mixed Operations**: Can set new tags and remove existing ones simultaneously  
- **Detailed Feedback**: Response shows all changes made
- **Null Handling**: `null` values properly remove tags from frontmatter
- **Validation**: Ensures required parameters and proper object structure

### Usage Examples
```
change_tag name="my-todo" tags={"status": "in_progress"}
change_tag name="task-123" tags={"todo": true, "status": "completed", "priority": "high"}
change_tag name="old-task" tags={"archived": true, "status": null, "todo": null}
```

This enhancement makes TODO state management much more efficient and supports complex tag operations in a single call.