---
title: 'Implementation: MCP Initial Instructions Feature'
tags:
  completed: true
  embedded-content: true
  implementation: true
  initial-instructions: true
  mcp: true
  workflow: true
created: 2025-08-24T00:50:21.553779695-07:00
modified: 2025-08-24T00:50:21.553779695-07:00
---

# MCP Initial Instructions Implementation

## Status: ✅ COMPLETED

Successfully implemented initial instructions support using the new mark3labs/mcp-go library.

## Implementation Details

### 1. Instructions File Creation
Created `internal/mcp/initial_instructions.md` containing:
- Agent memory usage guidelines
- **Search-first workflow emphasis** (principle #1)
- Aggressive memory management patterns
- TODO memory tracking approach
- Technical implementation documentation guidelines

### 2. Go Embed Integration
```go
//go:embed initial_instructions.md
var initialInstructions string
```

### 3. Server Configuration
```go
mcpServer := server.NewMCPServer(
    "simplemem",
    "0.1.0", 
    server.WithInstructions(initialInstructions), // NEW FEATURE!
    server.WithToolCapabilities(true),
)
```

## Key Instruction Principles Added

### Search-First Workflow
```markdown
### 1. ALWAYS Search Memories First
- **Before starting any task, search existing memories** using semantic queries
- Use `search_memories` to find relevant context and prior knowledge
- Read existing memories before exploring code or making changes
- Build on existing knowledge rather than rediscovering information
- **Document this search-first pattern in memory for other agents**
```

### Technical Implementation Documentation
```markdown
### Technical Implementation Details
- **Document step-by-step migration processes** in separate memories
- **Create focused memories for API changes** when migrating libraries
- **Record specific code patterns** and their transformations
- **Keep implementation memories small and focused** on specific aspects
```

## Benefits Achieved

1. **Agent Guidance**: Future agents will receive these instructions on initialization
2. **Search-First Pattern**: Prevents redundant work by emphasizing memory search
3. **Documentation Culture**: Encourages creating focused implementation memories
4. **Workflow Consistency**: Standardizes approach across agent sessions

## Technical Architecture

- **Embedded Content**: Instructions are compiled into binary (no external file dependencies)
- **Markdown Format**: Easy to read and edit
- **Version Controlled**: Instructions evolve with the codebase
- **MCP Protocol**: Delivered via standard MCP initialization process

## User Feedback Addressed

The user emphasized the importance of:
1. ✅ Searching memories before starting tasks
2. ✅ Creating small focused memories for technical work
3. ✅ Documenting implementation patterns as you work

All of these principles are now embedded in the initial instructions.

## Related
- [[mcp-library-migration-task]]
- [[workflow-pattern-search-first-approach]]
- File: `internal/mcp/initial_instructions.md`