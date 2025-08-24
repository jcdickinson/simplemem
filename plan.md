# SimpleMem MCP Server Implementation Plan

## Project Overview
Building an MCP (Model Context Protocol) server that allows storing documents as markdown files in a `.memories` directory. The server will support inter-document linking and will be prepared for future RAG (Retrieval-Augmented Generation) integration with DuckDB and VoyageAI.

## Current Status
- Created basic project structure with `cmd/simplemem` and `internal/` directories
- Started implementing MCP server (needs to be rewritten using go-mcp library)
- Module name: `github.com/jcdickinson/simplemem`

## Next Steps

### 1. Update Dependencies
```bash
go get github.com/mark3labs/go-mcp
```

### 2. Rewrite MCP Server
- Remove the custom MCP protocol implementation in `internal/mcp/`
- Use go-mcp library instead for proper MCP protocol handling

### 3. Implement Memory Storage (`internal/memory/store.go`)
- Store memories as `.md` files in `.memories` directory
- Support CRUD operations (Create, Read, Update, Delete)
- Handle markdown inter-document linking (e.g., `[[other-memory]]` style links)
- Ensure thread-safe file operations

### 4. Implement MCP Tools
The server should expose these tools via MCP:
- `create_memory` - Create a new memory document
- `read_memory` - Read a memory document
- `update_memory` - Update an existing memory
- `delete_memory` - Delete a memory
- `list_memories` - List all available memories
- `search_memories` - Basic text search across memories (will be enhanced with RAG later)

### 5. File Structure
```
simplemem/
├── cmd/
│   └── simplemem/
│       └── main.go          # Entry point
├── internal/
│   ├── mcp/
│   │   └── server.go        # MCP server using go-mcp
│   └── memory/
│       ├── store.go         # Memory storage logic
│       └── markdown.go      # Markdown handling & linking
├── .memories/               # Runtime directory for storing memories
├── go.mod
└── go.sum
```

### 6. Future Enhancements (Phase 2)
- Integrate DuckDB for efficient querying
- Add VoyageAI for semantic search/embeddings
- Implement RAG capabilities

## Technical Notes
- Use standard Go project layout
- Each memory is a standalone `.md` file
- Support markdown links between documents (e.g., `[Link Text](other-memory.md)` or `[[other-memory]]`)
- Ensure the `.memories` directory is created if it doesn't exist
- Consider adding metadata frontmatter to memories (YAML format)

## Dependencies Required
- `github.com/mark3labs/go-mcp` - MCP protocol implementation
- Future: DuckDB Go bindings
- Future: VoyageAI client library