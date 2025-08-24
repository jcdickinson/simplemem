# SimpleMem MCP Server Implementation Plan

## Project Overview
Building an MCP (Model Context Protocol) server that allows storing documents as markdown files in a `.memories` directory with full RAG (Retrieval-Augmented Generation) capabilities using DuckDB and VoyageAI for semantic search.

## ✅ Phase 1: Core MCP Server (COMPLETED)
- ✅ Basic project structure with `cmd/simplemem` and `internal/` directories
- ✅ MCP server implementation using `github.com/ThinkInAIXYZ/go-mcp`
- ✅ Memory storage with thread-safe file operations
- ✅ YAML frontmatter support with tags, timestamps, and metadata
- ✅ Enhanced markdown parsing with `github.com/gomarkdown/markdown`
- ✅ Inter-document linking (wiki-style `[[memory]]` and standard markdown links)
- ✅ 7 core MCP tools with intelligent augmentation:
  - `create_memory` - Create with auto frontmatter and timestamps
  - `read_memory` - Read with metadata summary and link extraction
  - `update_memory` - Update with timestamp preservation
  - `delete_memory` - Delete memory documents
  - `list_memories` - Enhanced listing with titles, tags, and dates
  - `search_memories` - Content search + tag search (`tag:tagname`, `tag:` for all tags)
  - `get_backlinks` - Find memories linking to a specific memory

## 🚀 Phase 2: RAG Implementation (CURRENT)

### 1. Database Integration
```bash
go get github.com/marcboeker/go-duckdb
```

#### DuckDB Schema Design
```sql
-- Memory documents table
CREATE TABLE memories (
    id INTEGER PRIMARY KEY,
    name VARCHAR UNIQUE,
    title VARCHAR,
    description TEXT,
    content TEXT,
    body TEXT,  -- content without frontmatter
    created TIMESTAMP,
    modified TIMESTAMP,
    file_hash VARCHAR,  -- for change detection
);

-- Tags table (normalized)
CREATE TABLE tags (
    id INTEGER PRIMARY KEY,
    memory_id INTEGER REFERENCES memories(id),
    tag_name VARCHAR,
    tag_value VARCHAR,  -- JSON for complex values
    INDEX (tag_name),
    INDEX (memory_id)
);

-- Links between memories
CREATE TABLE memory_links (
    id INTEGER PRIMARY KEY,
    from_memory_id INTEGER REFERENCES memories(id),
    to_memory_name VARCHAR,  -- target memory name
    link_text VARCHAR,
    link_type VARCHAR,  -- 'wiki' or 'markdown'
    INDEX (from_memory_id),
    INDEX (to_memory_name)
);

-- Vector embeddings
CREATE TABLE embeddings (
    id INTEGER PRIMARY KEY,
    memory_id INTEGER REFERENCES memories(id),
    chunk_text TEXT,
    chunk_index INTEGER,
    embedding FLOAT[1536],  -- VoyageAI dimension
    INDEX (memory_id)
);
```

### 2. VoyageAI Integration
```bash
go get github.com/voyage-ai/voyageai-go
```

#### Implementation Tasks
- Add `internal/db/duckdb.go` - Database operations and schema management
- Add `internal/embeddings/voyage.go` - VoyageAI client and embedding operations
- Add `internal/rag/` - RAG query processing and similarity search
- Update memory store to sync with database on changes
- Enhance search with semantic similarity

### 3. Enhanced RAG Features

#### New Search Capabilities
- **Semantic Search**: `search_memories` with similarity scoring
- **Hybrid Search**: Combine keyword + semantic search
- **Context-Aware**: Use backlinks and tags to improve relevance
- **Chunk-Based**: Split large memories into searchable chunks

#### Smart Memory Operations
- **Auto-tagging**: Suggest tags based on content analysis
- **Related Memories**: Find semantically similar memories
- **Content Summarization**: Generate descriptions from content
- **Link Suggestions**: Recommend memories to link to

### 4. File Structure (Updated)
```
simplemem/
├── cmd/
│   └── simplemem/
│       └── main.go          # Entry point
├── internal/
│   ├── mcp/
│   │   └── server.go        # MCP server (enhanced with RAG)
│   ├── memory/
│   │   ├── store.go         # File storage + DB sync
│   │   ├── frontmatter.go   # YAML frontmatter parsing
│   │   └── markdown.go      # Markdown + link handling
│   ├── db/
│   │   └── duckdb.go        # Database operations
│   ├── embeddings/
│   │   └── voyage.go        # VoyageAI integration
│   └── rag/
│       ├── search.go        # Semantic search
│       ├── chunks.go        # Content chunking
│       └── similarity.go    # Vector similarity
├── .memories/               # Checked-in memory files
├── .cache/
│   ├── simplemem.db         # DuckDB database
│   └── embeddings/          # Cached embedding data
├── justfile                 # Build commands
├── go.mod
└── go.sum
```

### 5. Configuration & Environment
- Add environment variables for VoyageAI API key
- Add configuration for embedding model selection
- Add caching strategies for embeddings
- Add batch processing for initial embedding generation

### 6. Enhanced MCP Tools (No New Tools Added)
Augment existing tools with RAG capabilities:
- `search_memories`: Add semantic search mode with `semantic:query`
- `read_memory`: Add related memories suggestions
- `create_memory`: Auto-generate suggested tags and links
- `list_memories`: Add similarity-based sorting options

## Technical Implementation Notes

### Database Strategy
- Use `.cache/simplemem.db` for DuckDB database (gitignored)
- Sync filesystem changes to database automatically
- Support offline operation (filesystem-first, DB optimization)
- Handle schema migrations gracefully

### Embedding Strategy
- Chunk large memories for better embedding quality
- Cache embeddings to avoid re-computation
- Support incremental updates when memories change
- Batch embed multiple memories for efficiency

### RAG Query Processing
1. Parse query (keyword vs semantic vs hybrid)
2. Generate embeddings for semantic queries
3. Execute similarity search in DuckDB
4. Combine with metadata filters (tags, dates, etc.)
5. Re-rank results using hybrid scoring
6. Return enriched results with context

## Dependencies
- ✅ `github.com/ThinkInAIXYZ/go-mcp` - MCP protocol
- ✅ `github.com/gomarkdown/markdown` - Markdown parsing
- ✅ `gopkg.in/yaml.v3` - YAML frontmatter
- 🔄 `github.com/marcboeker/go-duckdb` - Database integration
- 🔄 `github.com/voyage-ai/voyageai-go` - Embedding generation
- 🔄 Environment: `VOYAGE_API_KEY` - VoyageAI authentication

## Success Metrics
- Sub-100ms semantic search for <1000 memories
- Accurate similarity matching for related content
- Seamless offline/online operation
- Maintain existing 7-tool MCP interface
- Support for 10k+ memories with good performance