package mcp

import (
	"context"
	_ "embed"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/jcdickinson/simplemem/internal/config"
	"github.com/jcdickinson/simplemem/internal/memory"
)

//go:embed initial_instructions.md
var initialInstructions string

type Server struct {
	mcpServer     *server.MCPServer
	store         *memory.Store
	enhancedStore *memory.EnhancedStore
}


func NewServer(dbPath string) (*Server, error) {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Create enhanced store with custom db path
	enhancedStore, err := memory.NewEnhancedStoreWithDBPath(".memories", cfg, dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create enhanced store: %w", err)
	}

	s := &Server{
		store:         memory.NewStore(".memories"),
		enhancedStore: enhancedStore,
	}

	// Initialize the enhanced store (which also initializes the basic store)
	if err := s.enhancedStore.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize enhanced store: %w", err)
	}

	// Create MCP server with initial instructions support
	mcpServer := server.NewMCPServer(
		"simplemem",
		"0.1.0",
		server.WithInstructions(initialInstructions),
		server.WithToolCapabilities(true),
	)

	// Register all tools
	s.registerTools(mcpServer)

	s.mcpServer = mcpServer
	return s, nil
}

func (s *Server) registerTools(mcpServer *server.MCPServer) {
	// Create Memory tool
	mcpServer.AddTool(
		mcp.NewTool("create_memory",
			mcp.WithDescription("Create a new memory document. Supports YAML frontmatter for metadata including title, description, tags, and timestamps. Name can be specified either as a parameter or in the frontmatter 'name' field."),
			mcp.WithString("name",
				mcp.Description("Name of the memory (without .md extension). Optional if specified in frontmatter."),
			),
			mcp.WithString("content",
				mcp.Description("Content of the memory in markdown format"),
				mcp.Required(),
			),
		),
		s.handleCreateMemory,
	)

	// Read Memory tool
	mcpServer.AddTool(
		mcp.NewTool("read_memory",
			mcp.WithDescription("Read a memory document with full metadata including tags, timestamps, and links"),
			mcp.WithString("name",
				mcp.Description("Name of the memory to read"),
				mcp.Required(),
			),
		),
		s.handleReadMemory,
	)

	// Update Memory tool
	mcpServer.AddTool(
		mcp.NewTool("update_memory",
			mcp.WithDescription("Update an existing memory document. Name can be specified either as a parameter or in the frontmatter 'name' field."),
			mcp.WithString("name",
				mcp.Description("Name of the memory to update. Optional if specified in frontmatter."),
			),
			mcp.WithString("content",
				mcp.Description("New content for the memory"),
				mcp.Required(),
			),
		),
		s.handleUpdateMemory,
	)

	// Delete Memory tool
	mcpServer.AddTool(
		mcp.NewTool("delete_memory",
			mcp.WithDescription("Delete a memory document"),
			mcp.WithString("name",
				mcp.Description("Name of the memory to delete"),
				mcp.Required(),
			),
		),
		s.handleDeleteMemory,
	)

	// List Memories tool
	mcpServer.AddTool(
		mcp.NewTool("list_memories",
			mcp.WithDescription("List all memory documents with metadata preview including titles, tags, and modification dates"),
		),
		s.handleListMemories,
	)

	// Search Memories tool
	mcpServer.AddTool(
		mcp.NewTool("search_memories",
			mcp.WithDescription("Semantically search memories using natural language queries with optional tag filtering. Returns ranked results with snippets and relevance scores. Tags can filter by presence (empty value) or specific values."),
			mcp.WithString("query",
				mcp.Description("Semantic search query to find related memories"),
				mcp.Required(),
			),
			mcp.WithObject("tags",
				mcp.Description("Optional tag filters - key:value pairs. Use empty string as value to check for tag presence only"),
			),
			mcp.WithBoolean("require_all",
				mcp.Description("If true, memory must have ALL specified tags. If false, memory needs ANY of the tags (default: false)"),
			),
		),
		s.handleSearchMemories,
	)

	// Get Backlinks tool
	mcpServer.AddTool(
		mcp.NewTool("get_backlinks",
			mcp.WithDescription("Get memories related to a specific memory through explicit links and semantic similarity. Optionally rerank by query relevance."),
			mcp.WithString("name",
				mcp.Description("Name of the memory to find backlinks for"),
				mcp.Required(),
			),
			mcp.WithString("query",
				mcp.Description("Optional query to rerank backlinks by relevance"),
			),
		),
		s.handleGetBacklinks,
	)
}

func (s *Server) handleCreateMemory(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name := request.GetString("name", "")
	content := request.GetString("content", "")

	// Try to extract name from frontmatter if not provided in request
	if name == "" {
		fm, _, err := memory.ParseDocument(content)
		if err == nil && fm.Name != "" {
			name = fm.Name
		}
	}

	// Validate that we have a name
	if name == "" {
		return nil, fmt.Errorf("memory name must be provided either as parameter or in frontmatter 'name' field")
	}

	if err := s.enhancedStore.Create(name, content); err != nil {
		return nil, err
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("Memory '%s' created successfully", name),
			},
		},
	}, nil
}

func (s *Server) handleReadMemory(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name := request.GetString("name", "")

	memInfo, err := s.store.ReadWithMetadata(name)
	if err != nil {
		return nil, err
	}

	// Build response with metadata
	response := memInfo.Content

	// Add metadata information
	metaInfo := ""
	if memInfo.Frontmatter != nil {
		if memInfo.Frontmatter.Title != "" {
			metaInfo += fmt.Sprintf("\n📝 **Title:** %s", memInfo.Frontmatter.Title)
		}
		if memInfo.Frontmatter.Description != "" {
			metaInfo += fmt.Sprintf("\n📄 **Description:** %s", memInfo.Frontmatter.Description)
		}
		if len(memInfo.Frontmatter.Tags) > 0 {
			metaInfo += "\n🏷️ **Tags:** "
			var tagsList []string
			for tag, value := range memInfo.Frontmatter.Tags {
				if value == true {
					tagsList = append(tagsList, tag)
				} else {
					tagsList = append(tagsList, fmt.Sprintf("%s: %v", tag, value))
				}
			}
			metaInfo += strings.Join(tagsList, ", ")
		}
		if !memInfo.Frontmatter.Created.IsZero() {
			metaInfo += fmt.Sprintf("\n📅 **Created:** %s", memInfo.Frontmatter.Created.Format("2006-01-02 15:04:05"))
		}
		if !memInfo.Frontmatter.Modified.IsZero() {
			metaInfo += fmt.Sprintf("\n🔄 **Modified:** %s", memInfo.Frontmatter.Modified.Format("2006-01-02 15:04:05"))
		}
	}

	// Extract links from the content body
	links := memory.ExtractLinks(memInfo.Body)
	if len(links) > 0 {
		metaInfo += "\n\n🔗 **Links found in this memory:**\n"
		for _, link := range links {
			metaInfo += fmt.Sprintf("- [%s](%s) (%s link)\n", link.Text, link.Target, link.Type)
		}
	}

	if metaInfo != "" {
		response += "\n\n---" + metaInfo
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: response,
			},
		},
	}, nil
}

func (s *Server) handleUpdateMemory(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name := request.GetString("name", "")
	content := request.GetString("content", "")

	// Try to extract name from frontmatter if not provided in request
	if name == "" {
		fm, _, err := memory.ParseDocument(content)
		if err == nil && fm.Name != "" {
			name = fm.Name
		}
	}

	// Validate that we have a name
	if name == "" {
		return nil, fmt.Errorf("memory name must be provided either as parameter or in frontmatter 'name' field")
	}

	if err := s.enhancedStore.Update(name, content); err != nil {
		return nil, err
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("Memory '%s' updated successfully", name),
			},
		},
	}, nil
}

func (s *Server) handleDeleteMemory(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name := request.GetString("name", "")

	if err := s.enhancedStore.Delete(name); err != nil {
		return nil, err
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("Memory '%s' deleted successfully", name),
			},
		},
	}, nil
}

func (s *Server) handleListMemories(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	memories, err := s.store.List()
	if err != nil {
		return nil, err
	}

	if len(memories) == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "No memories found",
				},
			},
		}, nil
	}

	result := "Available memories:\n\n"
	for _, memory := range memories {
		memInfo, err := s.store.ReadWithMetadata(memory)
		if err != nil {
			result += fmt.Sprintf("📄 %s (error reading metadata)\n", memory)
			continue
		}

		result += fmt.Sprintf("📄 **%s**", memory)
		if memInfo.Frontmatter.Title != "" {
			result += fmt.Sprintf(" - %s", memInfo.Frontmatter.Title)
		}
		
		if len(memInfo.Frontmatter.Tags) > 0 {
			var tagsList []string
			for tag, value := range memInfo.Frontmatter.Tags {
				if value == true {
					tagsList = append(tagsList, tag)
				} else {
					tagsList = append(tagsList, fmt.Sprintf("%s:%v", tag, value))
				}
			}
			result += fmt.Sprintf(" 🏷️[%s]", strings.Join(tagsList, ", "))
		}

		if !memInfo.Frontmatter.Modified.IsZero() {
			result += fmt.Sprintf(" (modified: %s)", memInfo.Frontmatter.Modified.Format("2006-01-02"))
		}

		if memInfo.Frontmatter.Description != "" {
			result += fmt.Sprintf("\n  %s", memInfo.Frontmatter.Description)
		}

		result += "\n\n"
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: result,
			},
		},
	}, nil
}

func (s *Server) handleSearchMemories(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query := request.GetString("query", "")
	
	// Get tags and require_all from arguments
	args := request.GetArguments()
	var tags map[string]string
	var requireAll bool
	
	if tagsArg, ok := args["tags"]; ok {
		if tagsMap, ok := tagsArg.(map[string]interface{}); ok {
			tags = make(map[string]string)
			for k, v := range tagsMap {
				tags[k] = fmt.Sprintf("%v", v)
			}
		}
	}
	
	if requireAllArg, ok := args["require_all"]; ok {
		requireAll, _ = requireAllArg.(bool)
	}

	// Use semantic search with tag filtering (set to 5 docs as requested)
	result, err := s.enhancedStore.SearchSemanticMarkdownWithTags(query, tags, requireAll, 5)
	if err != nil {
		return nil, fmt.Errorf("failed to perform semantic search: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: result,
			},
		},
	}, nil
}

func (s *Server) handleGetBacklinks(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name := request.GetString("name", "")
	query := request.GetString("query", "")

	// Use query for reranking if provided, otherwise use memory name as query
	if query == "" {
		query = name // Use memory name as default query for reranking
	}

	// Get enhanced backlinks with reranking (set to 5 docs as requested)
	result, err := s.enhancedStore.GetEnhancedBacklinks(name, query, 5)
	if err != nil {
		return nil, fmt.Errorf("failed to get enhanced backlinks: %w", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: result,
			},
		},
	}, nil
}

func (s *Server) Run() error {
	return server.ServeStdio(s.mcpServer)
}

func (s *Server) Shutdown(ctx context.Context) error {
	// Close enhanced store
	if s.enhancedStore != nil {
		if err := s.enhancedStore.Close(); err != nil {
			// Log error but don't fail shutdown
			fmt.Printf("Warning: failed to close enhanced store: %v\n", err)
		}
	}
	// New MCP library handles shutdown automatically
	return nil
}