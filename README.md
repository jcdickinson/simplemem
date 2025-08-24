# SimpleMem

A vibe-coded memory management system with RAG capabilities for Claude via the Model Context Protocol (MCP).

> ⚠️ **Warning**: This project is completely vibe-coded. It works, but don't look too closely at the implementation details. The code was written by an AI having a good time, not by someone following best practices.

## What is SimpleMem?

SimpleMem is an MCP server that provides persistent memory storage and retrieval for Claude and other MCP clients. It combines traditional file-based storage with modern RAG (Retrieval-Augmented Generation) capabilities, including semantic search and automatic relationship discovery.

Think of it as Claude's personal notebook that never forgets and can find connections between ideas automatically.

## Features

- 📝 **Persistent Memory Storage**: Store and retrieve memories with rich metadata
- 🔍 **Semantic Search**: Find memories using natural language queries
- 🔗 **Automatic Relationship Discovery**: Semantic backlinks connect related memories
- 🏷️ **Tag System**: Organize memories with flexible tagging
- 🎯 **Vector Embeddings**: Powered by Voyage AI for high-quality semantic understanding
- 📊 **DuckDB Backend**: Fast, efficient storage with vector similarity search
- 🔧 **MCP Protocol**: Seamless integration with Claude and other MCP clients

## Quick Start

### Prerequisites

- Go 1.21+
- A Voyage AI API key (set as `VOYAGE_API_KEY` environment variable)

### Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/simplemem
cd simplemem

# Install dependencies
just deps

# Build the binary
just build
```

### Usage

#### As an MCP Server

Add to your MCP client configuration:

```json
{
  "mcpServers": {
    "simplemem": {
      "command": "./simplemem",
      "args": ["--db", "path/to/your/database.db"]
    }
  }
}
```

#### Command Line Testing

```bash
# List all memories
just test-json '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"list_memories","arguments":{}},"id":1}'

# Create a memory
just test-json '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"create_memory","arguments":{"content":"---\nname: my-memory\ntitle: My First Memory\n---\n\n# Hello World\n\nThis is my first memory!"}},"id":1}'

# Search memories
just test-json '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"search_memories","arguments":{"query":"hello world"}},"id":1}'
```

#### Testing Semantic Backlinks

```bash
# Run comprehensive semantic backlinks test
just test-backlinks

# Verify backlinks are working
just test-json '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"get_backlinks","arguments":{"name":"my-memory"}},"id":1}'
```

## Available Tools (MCP)

- **`list_memories`**: List all stored memories with metadata
- **`create_memory`**: Create a new memory from markdown content  
- **`read_memory`**: Read a specific memory by name
- **`update_memory`**: Update existing memory content
- **`delete_memory`**: Remove a memory and all related data
- **`search_memories`**: Semantic search with optional tag filtering
- **`get_backlinks`**: Get memories related to a specific memory
- **`change_tag`**: Modify tags on memories

## Memory Format

Memories use YAML frontmatter + Markdown:

```markdown
---
name: my-memory-name
title: A Human-Readable Title
description: Optional description
tags:
  category: personal
  priority: high
  status: active
---

# Memory Content

Your memory content goes here in **Markdown** format.

- You can use lists
- **Bold text**
- Links to [[other-memories]]
- Whatever you need!
```

## Development

### Prerequisites

- [Just](https://github.com/casey/just) command runner
- [Jujutsu](https://github.com/martinvonz/jj) for version control (or Git)

### Common Commands

```bash
# Show available commands
just

# Run tests (when they exist)
just test

# Format code
just fmt

# Run with verbose logging
just run-verbose

# Test semantic backlinks functionality
just test-backlinks

# Clean up build artifacts
just clean
```

### Debug Testing

For rapid development iteration:

```bash
# Test any JSON-RPC call with custom database
just test-json '<json>' /tmp/test.db

# Quick minimal test
just test-clean /tmp/debug.db

# Build and test in one command  
just build-test
```

## Architecture

- **Frontend**: MCP JSON-RPC 2.0 protocol server
- **Storage**: File-based memory storage with DuckDB backend
- **Search**: Voyage AI embeddings + vector similarity
- **Relationships**: Automatic semantic backlink discovery
- **Memory**: YAML frontmatter + Markdown content

## Known Issues

- No proper test suite (it's vibe-coded, remember?)
- Error handling could be more robust
- Some edge cases might crash things
- Documentation is whatever this README covers
- Database migrations? What are those?

## Contributing

This is a vibe-coded project, so contributions should match the energy:

1. Make it work first
2. Make it work well second  
3. Make it pretty... maybe later?
4. Tests are aspirational
5. If it breaks, fix it with more vibes

## License

MIT License - Use it, modify it, vibe with it.

## Credits

- Built with love and caffeine
- Powered by [Voyage AI](https://voyageai.com) embeddings
- Uses [DuckDB](https://duckdb.org) for storage  
- MCP integration via [mark3labs/mcp-go](https://github.com/mark3labs/mcp-go)

---

*"It works on my machine, and that machine has good vibes."* ✨