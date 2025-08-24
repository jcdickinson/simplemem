---
title: SimpleMem Timestamp Management Design
description: How SimpleMem handles created/modified timestamps in frontmatter
tags:
  architecture: true
  frontmatter: design
  timestamps: server-managed
created: 2025-08-24T01:29:07.630068798-07:00
modified: 2025-08-24T01:29:07.630068798-07:00
---

# SimpleMem Timestamp Management Design

## Philosophy
SimpleMem uses server-managed timestamps to ensure accuracy and prevent user manipulation of creation/modification times. This is especially important because Git doesn't preserve file modification times.

## Implementation

### Server-Side Timestamp Handling
1. **On Create**: Server automatically sets both `created` and `modified` timestamps to current time
2. **On Update**: Server preserves existing `created` timestamp and updates `modified` to current time
3. **User Input Ignored**: Any timestamps provided by users in frontmatter are automatically overwritten

### Storage Strategy
- **Stored on Disk**: Timestamps are written to frontmatter in the actual files
- **Git-Safe**: Since Git clobbers file mtimes, we store timestamps in document content
- **Reliable**: Server-managed timestamps ensure consistency across operations

## Code Implementation

### Store.Create Method
```go
// Parse user content (may contain user-provided timestamps)
fm, body, err := ParseDocument(content)

// Server overwrites with correct timestamps
fm.UpdateTimestamps(true) // Sets both created and modified

// Save to disk with server timestamps
```

### Store.Update Method  
```go
// Read existing file to get created timestamp
existingFm, _, _ := ParseDocument(string(existingData))

// Parse new content (may contain user timestamps)
fm, body, err := ParseDocument(content) 

// Preserve server-managed created time
if !existingFm.Created.IsZero() {
    fm.Created = existingFm.Created
}

// Server updates modified timestamp
fm.UpdateTimestamps(false) // Only updates modified time
```

## MCP Interface

### Tool Descriptions
- **create_memory**: "Timestamps are automatically managed by the server"
- **update_memory**: "Any timestamps in frontmatter will be overwritten by server-managed values"

### User Experience
Users can include timestamps in their frontmatter if they want, but they will be ignored:

**User Input:**
```yaml
---
name: my-memory
title: My Memory
created: 2020-01-01T00:00:00Z  # This will be ignored
modified: 2020-01-01T00:00:00Z # This will be ignored
---
```

**Stored Result:**
```yaml
---
title: My Memory
created: 2025-08-24T01:30:00Z  # Server-generated
modified: 2025-08-24T01:30:00Z # Server-generated
---
```

## Benefits

1. **Accuracy**: Timestamps reflect actual server-side operations
2. **Security**: Users cannot forge creation/modification times
3. **Consistency**: All memories have reliable timestamps
4. **Git-Compatible**: Timestamps survive Git operations
5. **User-Friendly**: Users don't need to manage timestamps manually

## Related
- [[mcp-frontmatter-improvements]] - Overall frontmatter API improvements
- [[simplemem-memory-frontmatter]] - Frontmatter system architecture