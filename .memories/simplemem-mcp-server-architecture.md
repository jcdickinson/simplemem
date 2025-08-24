---
title: SimpleMem MCP Server Architecture
description: Detailed analysis of the MCP server implementation and request handling
tags:
  architecture: true
  golang: true
  mcp: true
  protocol: true
  server: true
created: 2025-08-24T00:19:15.782028082-07:00
modified: 2025-08-24T01:04:47.560057942-07:00
---

# SimpleMem MCP Server Architecture

The MCP server (`internal/mcp/server.go`) implements the Model Context Protocol to provide memory management tools for Claude and other MCP clients. The server follows a clean architecture with embedded stores and tool-based request handling.

## Server Structure

### Core Server Type

```go
type Server struct {
    store         *memory.Store         // Basic file operations
    enhancedStore *memory.EnhancedStore // RAG-enabled operations
    mcpServer     *server.MCPServer     // MCP protocol handler
}
```

### Component Relationships

**Enhanced Store Wrapping**
- Basic `Store` handles file I/O, markdown parsing, frontmatter
- `EnhancedStore` wraps basic store with database sync and RAG capabilities
- MCP server uses enhanced store for all operations to ensure consistency

**MCP Protocol Integration**
- Uses `mark3labs/mcp-go` library for protocol implementation
- Handles JSON-RPC 2.0 communication with proper error management
- Supports tool capabilities and initial instructions

## Tool Registration

### Tool Architecture Pattern

Each MCP tool follows a consistent registration pattern:

```go
mcpServer.AddTool(
    mcp.NewTool("tool_name",
        mcp.WithDescription("Tool description"),
        mcp.WithString("param", mcp.Description("Parameter description"), mcp.Required()),
    ),
    s.handleToolName,
)
```

### Tool Implementation Strategy

**Parameter Extraction**
```go
name := request.GetString("name", "")
content := request.GetString("content", "")
```

**Enhanced Store Operations**
All tools use the enhanced store to ensure:
- Database synchronization
- RAG processing
- Consistent metadata handling

**Response Formatting**
```go
return &mcp.CallToolResult{
    Content: []mcp.Content{
        mcp.TextContent{Type: "text", Text: result},
    },
}, nil
```

## Core Tools Implementation

### Memory Management Tools

#### create_memory
- **Purpose**: Create new memory with automatic metadata
- **Process**: Enhanced store handles file creation + database sync + RAG processing
- **Features**: Auto-timestamps, frontmatter generation, embedding creation

#### read_memory  
- **Purpose**: Read memory with rich metadata presentation
- **Features**: Metadata summary, link extraction, formatting with emojis
- **Enhancement**: Could include related memory suggestions

#### update_memory
- **Purpose**: Update existing memory preserving metadata  
- **Process**: Enhanced store handles file update + database sync + RAG reprocessing
- **Features**: Timestamp preservation, change detection via file hash

#### delete_memory
- **Purpose**: Remove memory from both filesystem and database
- **Process**: Cascading delete (file → database → embeddings → indexes)
- **Safety**: Error handling ensures partial failures don't corrupt state

### Discovery Tools

#### list_memories
- **Purpose**: Enhanced memory listing with metadata preview
- **Features**: Title display, tag summary, modification dates
- **Format**: Rich text with emojis and consistent formatting
- **Performance**: Uses basic store for file listing, enhanced for metadata

#### search_memories
- **Purpose**: Semantic search with optional tag filtering
- **Algorithm**: Pure semantic search using vector embeddings
- **Features**: Tag filters (presence and value matching), boolean logic
- **Response**: Markdown formatted results with snippets and scores
- **Limit**: Fixed at 5 results for optimal relevance

#### get_backlinks  
- **Purpose**: Find memories linking to a specific memory
- **Features**: Explicit link detection + semantic similarity
- **Enhancement**: Optional query-based reranking for context-specific relevance
- **Format**: Markdown with context snippets and relationship types

## Request Processing Pipeline

### Standard Request Flow

1. **JSON-RPC Parsing**: MCP library handles protocol details
2. **Parameter Extraction**: Type-safe parameter parsing with defaults
3. **Enhanced Store Operation**: All operations go through enhanced store
4. **Database Sync**: Automatic sync for write operations
5. **RAG Processing**: Embedding generation/updates as needed
6. **Response Assembly**: Consistent text content response format

### Error Handling Strategy

**Validation Errors**
- Parameter validation before processing
- Required field checking via MCP library

**Store Errors** 
- File system errors (permissions, disk space, etc.)
- Database errors (connection, constraint violations, etc.)
- RAG processing errors (API failures, embedding issues, etc.)

**Recovery Patterns**
- Graceful degradation (semantic search falls back to basic search)
- Partial success handling (warn on database sync failure, succeed on file operation)
- User-friendly error messages with context

## Advanced Features

### Enhanced Search Implementation

#### Semantic Search Processing
```go
// Use semantic search with tag filtering (set to 5 docs as requested)
result, err := s.enhancedStore.SearchSemanticMarkdownWithTags(query, tags, requireAll, 5)
```

**Query Processing Steps:**
1. Parse query string and tag filters from request arguments
2. Convert tag map from `interface{}` to `string` format
3. Apply semantic search with vector similarity
4. Filter results by tag criteria (AND/OR logic)
5. Format results as markdown with snippets and relevance scores

#### Backlink Enhancement
```go
// Get enhanced backlinks with reranking (set to 5 docs as requested)  
result, err := s.enhancedStore.GetEnhancedBacklinks(name, query, 5)
```

**Backlink Discovery Process:**
1. Find explicit markdown links (`[text](target)` and `[[target]]`)
2. Find semantic backlinks using memory embedding as query
3. Optionally rerank by query relevance if query provided
4. Combine results with context snippets
5. Format as markdown with relationship type indicators

### Configuration and Lifecycle

#### Server Initialization
```go
func NewServer(dbPath string) (*Server, error) {
    // Load configuration
    // Create enhanced store with custom db path  
    // Initialize the enhanced store
    // Create MCP server with initial instructions support
    // Register all tools
}
```

**Initialization Steps:**
1. **Configuration Loading**: Environment variables and config files
2. **Enhanced Store Creation**: Database initialization, RAG processor setup
3. **Store Initialization**: File system sync, embedding processing
4. **MCP Server Creation**: Protocol setup with initial instructions
5. **Tool Registration**: All 7 tools registered with proper handlers

#### Embedded Initial Instructions

The server includes embedded initial instructions that teach agents to use SimpleMem effectively:
- Memory-first approach for knowledge persistence  
- Search-first workflow for building on existing knowledge
- Aggressive memory management patterns
- TODO tracking with memory integration

## Performance Characteristics

### Request Latency
- **File Operations**: ~1-5ms for typical memory sizes
- **Database Sync**: ~5-10ms with indexes
- **Semantic Search**: ~50-200ms depending on collection size  
- **Embedding Generation**: ~100-500ms via VoyageAI API

### Memory Usage
- **Server Process**: ~10-50MB baseline
- **Database**: ~6KB per memory + embeddings
- **Cache**: Configurable embedding cache size

### Scalability Considerations
- **Memory Limit**: Practical limit ~10,000 memories
- **Concurrent Requests**: Thread-safe via enhanced store synchronization
- **Database Performance**: Sub-second search for typical usage patterns

The MCP server provides a robust, feature-rich interface that bridges simple file-based memory storage with advanced RAG capabilities, ensuring both ease of use and powerful semantic search functionality.