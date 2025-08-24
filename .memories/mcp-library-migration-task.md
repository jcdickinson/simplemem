---
title: '[TODO] Migrate MCP Library to mark3labs/mcp-go'
tags:
  area: refactor
  component: mcp
  priority: high
  status: completed
  todo: true
created: 2025-08-24T00:39:32.950834948-07:00
modified: 2025-08-24T00:50:02.114623991-07:00
---

# Task Description
Migrate from current MCP library (ThinkInAIXYZ/go-mcp v0.2.21) to the new mark3labs/mcp-go library to enable initial instructions support.

## Status: ✅ COMPLETED SUCCESSFULLY

## Context
- Current library: `github.com/ThinkInAIXYZ/go-mcp v0.2.21` (auto-updated by go mod tidy)
- Target libraries:
  - `github.com/mark3labs/mcp-go/mcp` ✅
  - `github.com/mark3labs/mcp-go/server` ✅
- Need to implement initial instructions functionality ✅
- User wants to maintain simplemem memories updated throughout migration ✅

## Migration Steps Completed
1. ✅ Update go.mod dependencies
2. ✅ Update imports in server.go
3. ✅ Update server creation with WithInstructions
4. ✅ Convert all tool registrations (7 tools)
5. ✅ Update handler function signatures (7 handlers)
6. ✅ Fix request/response handling
7. ✅ Update transport handling
8. ✅ Test compilation - PASSED

## Key Achievements

### Initial Instructions Implemented ✅
- Created `internal/mcp/initial_instructions.md` with embedded instructions
- Used `//go:embed` to include instructions in binary
- Added `server.WithInstructions()` option to server creation
- Instructions emphasize **search-first workflow** for agents

### Complete API Migration ✅
- All 7 MCP tools successfully migrated
- All 7 handler functions updated
- Tool registration uses new builder pattern
- Parameter handling uses `GetString()`, `GetArguments()` methods
- Response handling uses value types instead of pointers

### Clean Architecture ✅
- Removed all old struct request types
- Updated function signatures throughout
- Maintained all existing functionality
- Build passes without errors

## New Features Enabled
- **Initial Instructions Support**: Agents will receive instructions on initialization
- **Search-First Emphasis**: Instructions guide agents to search memories before exploring
- **Embedded Instructions**: Instructions are built into the binary, not external file dependency

## Technical Implementation Memories Created
- [[implementation-mcp-tool-registration-migration]] - Tool registration patterns
- [[implementation-mcp-handler-signature-migration]] - Handler function updates  
- [[workflow-pattern-search-first-approach]] - Search-first workflow documentation

## Acceptance Criteria
- ✅ Update go.mod to use new MCP library
- ✅ Migrate server.go to new library API  
- ✅ Implement initial instructions support
- ✅ Maintain all existing functionality
- ✅ Update memories with migration details
- ⏳ Test migrated implementation (next step)

## Related
- [[simplemem-mcp-server-architecture]]
- [[knowledge-project-dependencies]]
- [[knowledge-current-mcp-implementation]]
- [[knowledge-new-mcp-library-analysis]]