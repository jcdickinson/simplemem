---
title: 'MCP API Refactor: Frontmatter to Metadata Object'
description: Successfully refactored SimpleMem MCP server to move frontmatter from content field to separate metadata object in tool requests
tags:
  api: true
  completed: true
  mcp: true
  priority: high
  refactor: true
created: 2025-08-24T22:15:47.142596546-07:00
modified: 2025-08-24T22:18:05.684433142-07:00
tested: true
version: "1.0"
---

# MCP API Refactor: Frontmatter to Metadata Object

Successfully refactored SimpleMem MCP server to move frontmatter from content field to separate metadata object in tool requests. **TESTED AND WORKING!** ✅

## Changes Made

### Tool Schema Updates
- **create_memory**: Now requires separate `metadata` object parameter
- **update_memory**: Now requires separate `metadata` object parameter  
- Removed support for name-in-frontmatter feature
- Made `metadata` parameter required to encourage proper usage

### Metadata Structure
```json
{
  "name": "memory-name",
  "metadata": {
    "title": "Required Title",        // REQUIRED
    "description": "Required desc",  // REQUIRED  
    "tags": {                       // REQUIRED
      "key": "value",
      "bool_tag": true
    },
    "custom_field": "any_value"     // Any additional properties allowed
  },
  "content": "Clean markdown content without frontmatter"
}
```

### Implementation Details
- Validates required fields: title (string), description (string), tags (object)
- Supports arbitrary additional metadata properties in `fm.Metadata`
- Constructs proper frontmatter internally using `memory.FormatDocument()`
- Maintains backward compatibility for on-disk storage and responses
- Frontmatter still appears in file storage and tool responses

### Benefits
1. **Clean separation**: Content is pure markdown, metadata is structured
2. **Type safety**: Required fields are validated at API level
3. **Extensibility**: Arbitrary metadata properties supported
4. **Consistency**: Enforces proper metadata usage across all memories

### Testing Results
- ✅ Build successful
- ✅ Create memory with new API works
- ✅ Update memory with new API works
- ✅ Custom metadata fields properly supported
- ✅ Required field validation working

The refactor maintains full compatibility with existing file storage and response formats while providing a cleaner, more structured API for creating and updating memories.