package mcp

import (
	"context"
	_ "embed"
	"fmt"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
	"github.com/jcdickinson/simplemem/internal/config"
	"github.com/jcdickinson/simplemem/internal/memory"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

//go:embed initial_instructions.md
var initialInstructions string

type Server struct {
	mcpServer     *server.MCPServer
	store         *memory.Store
	enhancedStore *memory.EnhancedStore
	config        *config.Config
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
		config:        cfg,
	}

	// Initialize the enhanced store (which also initializes the basic store)
	if err := s.enhancedStore.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize enhanced store: %w", err)
	}

	actualInitialInstructions := initialInstructions

	if cfg.MaxMemoryLength > 0 {
		actualInitialInstructions += fmt.Sprintf("\n\nIMPORTANT: Memories have a maximum length of %d characters (plain-text, with markdown formatting stripped) in order to encourage creating focused memories. If you need to explain something longer than that, consider splitting it into multiple memories.", cfg.MaxMemoryLength)
	}

	// Create MCP server with initial instructions support
	mcpServer := server.NewMCPServer(
		"simplemem",
		"0.1.0",
		server.WithInstructions(actualInitialInstructions),
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
			mcp.WithDescription("Create a new memory document. Requires metadata object for title, description, and tags. Content length is limited to 2000 characters."),
			mcp.WithString("name",
				mcp.Description("Name of the memory (without .md extension)"),
				mcp.Required(),
			),
			mcp.WithObject("metadata",
				mcp.Description("Metadata object with required title, description, and tags fields, plus any additional properties"),
				mcp.Required(),
			),
			mcp.WithString("content",
				mcp.Description("Content of the memory in markdown format"),
				mcp.Required(),
			),
		),
		s.handleCreateMemory,
	)

	// Update Memory tool
	mcpServer.AddTool(
		mcp.NewTool("update_memory",
			mcp.WithDescription("Update an existing memory document. Requires metadata object for title, description, and tags. Timestamps (created/modified) are automatically managed by the server. Content length is limited to 2000 characters."),
			mcp.WithString("name",
				mcp.Description("Name of the memory to update"),
				mcp.Required(),
			),
			mcp.WithObject("metadata",
				mcp.Description("Metadata object with required title, description, and tags fields, plus any additional properties"),
				mcp.Required(),
			),
			mcp.WithString("content",
				mcp.Description("New content for the memory"),
				mcp.Required(),
			),
		),
		s.handleUpdateMemory,
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

	// List Memories tool - temporarily removed to encourage semantic search usage
	// mcpServer.AddTool(
	// 	mcp.NewTool("list_memories",
	// 		mcp.WithDescription("List all memory documents with metadata preview including titles, tags, and modification dates"),
	// 	),
	// 	s.handleListMemories,
	// )

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

	// Change Tag tool
	mcpServer.AddTool(
		mcp.NewTool("change_tag",
			mcp.WithDescription("Change multiple tags on a memory document. Useful for setting TODO states and other metadata. Tags with null values will be removed. Example: {\"todo\": true, \"status\": \"in_progress\", \"priority\": \"high\", \"old_tag\": null}"),
			mcp.WithString("name",
				mcp.Description("Name of the memory to modify"),
				mcp.Required(),
			),
			mcp.WithObject("tags",
				mcp.Description("Object containing tag key-value pairs to set. Use null values to remove tags."),
				mcp.Required(),
			),
		),
		s.handleChangeTag,
	)
}

// textExtractor implements ast.NodeVisitor to extract plain text from markdown AST
type textExtractor struct {
	result strings.Builder
}

func (te *textExtractor) Visit(node ast.Node, entering bool) ast.WalkStatus {
	switch n := node.(type) {
	case *ast.Text:
		if entering {
			te.result.Write(n.Literal)
		}
	case *ast.Paragraph, *ast.Heading, *ast.BlockQuote, *ast.List, *ast.ListItem:
		if !entering {
			te.result.WriteString(" ")
		}
	case *ast.Softbreak, *ast.Hardbreak:
		if entering {
			te.result.WriteString(" ")
		}
	}
	return ast.GoToNext
}

// stripMarkdown converts markdown to plain text by parsing and then extracting text content
func stripMarkdown(content string) string {
	// Create a parser with common extensions
	p := parser.NewWithExtensions(parser.CommonExtensions)
	
	// Parse the markdown content
	doc := markdown.Parse([]byte(content), p)
	
	// Extract plain text from the AST
	extractor := &textExtractor{}
	ast.Walk(doc, extractor)
	
	return strings.TrimSpace(extractor.result.String())
}

func (s *Server) validateMemoryLength(content string) error {
	// If max_memory_length is <= 0, disable length check
	if s.config.MaxMemoryLength <= 0 {
		return nil
	}

	// Strip markdown to get plain text length
	plainText := stripMarkdown(content)
	plainTextLength := len(plainText)

	if plainTextLength > s.config.MaxMemoryLength {
		return fmt.Errorf("memory content exceeds maximum length of %d characters (%d plain-text characters provided)", s.config.MaxMemoryLength, plainTextLength)
	}

	return nil
}

func (s *Server) handleCreateMemory(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name := request.GetString("name", "")
	content := request.GetString("content", "")

	// Get metadata object from request
	args := request.GetArguments()
	metadataArg, ok := args["metadata"]
	if !ok {
		return nil, fmt.Errorf("metadata parameter is required")
	}

	metadataMap, ok := metadataArg.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("metadata must be an object")
	}

	// Validate that we have a name
	if name == "" {
		return nil, fmt.Errorf("memory name is required")
	}

	// Validate required fields
	title, titleExists := metadataMap["title"]
	if !titleExists {
		return nil, fmt.Errorf("metadata.title is required")
	}
	titleStr, ok := title.(string)
	if !ok {
		return nil, fmt.Errorf("metadata.title must be a string")
	}

	description, descExists := metadataMap["description"]
	if !descExists {
		return nil, fmt.Errorf("metadata.description is required")
	}
	descStr, ok := description.(string)
	if !ok {
		return nil, fmt.Errorf("metadata.description must be a string")
	}

	tags, tagsExists := metadataMap["tags"]
	if !tagsExists {
		return nil, fmt.Errorf("metadata.tags is required")
	}
	tagsMap, ok := tags.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("metadata.tags must be an object")
	}

	// Build frontmatter from metadata object
	fm := &memory.Frontmatter{
		Title:       titleStr,
		Description: descStr,
		Tags:        tagsMap,
		Metadata:    make(map[string]interface{}),
	}

	// Add any additional metadata properties (excluding the standard ones)
	for key, value := range metadataMap {
		if key != "title" && key != "description" && key != "tags" {
			fm.Metadata[key] = value
		}
	}

	// Create document content with frontmatter
	finalContent, err := memory.FormatDocument(fm, content)
	if err != nil {
		return nil, fmt.Errorf("failed to format document: %w", err)
	}

	// Validate memory length
	if err := s.validateMemoryLength(finalContent); err != nil {
		return nil, err
	}

	if err := s.enhancedStore.Create(name, finalContent); err != nil {
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

	// Get metadata object from request
	args := request.GetArguments()
	metadataArg, ok := args["metadata"]
	if !ok {
		return nil, fmt.Errorf("metadata parameter is required")
	}

	metadataMap, ok := metadataArg.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("metadata must be an object")
	}

	// Validate that we have a name
	if name == "" {
		return nil, fmt.Errorf("memory name is required")
	}

	// Validate required fields
	title, titleExists := metadataMap["title"]
	if !titleExists {
		return nil, fmt.Errorf("metadata.title is required")
	}
	titleStr, ok := title.(string)
	if !ok {
		return nil, fmt.Errorf("metadata.title must be a string")
	}

	description, descExists := metadataMap["description"]
	if !descExists {
		return nil, fmt.Errorf("metadata.description is required")
	}
	descStr, ok := description.(string)
	if !ok {
		return nil, fmt.Errorf("metadata.description must be a string")
	}

	tags, tagsExists := metadataMap["tags"]
	if !tagsExists {
		return nil, fmt.Errorf("metadata.tags is required")
	}
	tagsMap, ok := tags.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("metadata.tags must be an object")
	}

	// Build frontmatter from metadata object
	fm := &memory.Frontmatter{
		Title:       titleStr,
		Description: descStr,
		Tags:        tagsMap,
		Metadata:    make(map[string]interface{}),
	}

	// Add any additional metadata properties (excluding the standard ones)
	for key, value := range metadataMap {
		if key != "title" && key != "description" && key != "tags" {
			fm.Metadata[key] = value
		}
	}

	// Create document content with frontmatter
	finalContent, err := memory.FormatDocument(fm, content)
	if err != nil {
		return nil, fmt.Errorf("failed to format document: %w", err)
	}

	// Validate memory length
	if err := s.validateMemoryLength(finalContent); err != nil {
		return nil, err
	}

	if err := s.enhancedStore.Update(name, finalContent); err != nil {
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

		// Calculate plain-text content length (markdown stripped)
		plainText := stripMarkdown(memInfo.Content)
		contentLength := len(plainText)

		result += fmt.Sprintf("📄 **%s**", memory)
		if memInfo.Frontmatter.Title != "" {
			result += fmt.Sprintf(" - %s", memInfo.Frontmatter.Title)
		}

		result += fmt.Sprintf(" (%d chars)", contentLength)

		// Warn if memory exceeds configured maximum length
		if s.config.MaxMemoryLength > 0 && contentLength > s.config.MaxMemoryLength {
			result += fmt.Sprintf(" ⚠️ OVER LIMIT (%d/%d)", contentLength, s.config.MaxMemoryLength)
		}

		result += "\n"

		// Show all frontmatter fields
		if memInfo.Frontmatter.Description != "" {
			result += fmt.Sprintf("  📄 **Description:** %s\n", memInfo.Frontmatter.Description)
		}

		if len(memInfo.Frontmatter.Tags) > 0 {
			result += "  🏷️ **Tags:**\n"
			for tag, value := range memInfo.Frontmatter.Tags {
				if value == true {
					result += fmt.Sprintf("    - %s\n", tag)
				} else {
					result += fmt.Sprintf("    - %s: %v\n", tag, value)
				}
			}
		}

		if !memInfo.Frontmatter.Created.IsZero() {
			result += fmt.Sprintf("  📅 **Created:** %s\n", memInfo.Frontmatter.Created.Format("2006-01-02 15:04:05"))
		}

		if !memInfo.Frontmatter.Modified.IsZero() {
			result += fmt.Sprintf("  🔄 **Modified:** %s\n", memInfo.Frontmatter.Modified.Format("2006-01-02 15:04:05"))
		}

		// Show non-required fields as YAML
		nonRequiredFields := make(map[string]interface{})
		
		if len(memInfo.Frontmatter.Links) > 0 {
			nonRequiredFields["links"] = memInfo.Frontmatter.Links
		}
		
		// Add any additional metadata
		for key, value := range memInfo.Frontmatter.Metadata {
			nonRequiredFields[key] = value
		}
		
		if len(nonRequiredFields) > 0 {
			result += "  ```yaml\n"
			for key, value := range nonRequiredFields {
				result += fmt.Sprintf("  %s: %v\n", key, value)
			}
			result += "  ```\n"
		}

		result += "\n"
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

func (s *Server) handleChangeTag(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name := request.GetString("name", "")

	// Get tags from arguments
	args := request.GetArguments()
	tagsArg, ok := args["tags"]
	if !ok {
		return nil, fmt.Errorf("tags parameter is required")
	}

	tagsMap, ok := tagsArg.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("tags must be an object")
	}

	if name == "" {
		return nil, fmt.Errorf("memory name is required")
	}
	if len(tagsMap) == 0 {
		return nil, fmt.Errorf("at least one tag must be specified")
	}

	// Read the current memory
	memInfo, err := s.store.ReadWithMetadata(name)
	if err != nil {
		return nil, fmt.Errorf("failed to read memory '%s': %w", name, err)
	}

	// Parse the current document to get frontmatter and body
	fm, body, err := memory.ParseDocument(memInfo.Content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse memory document: %w", err)
	}

	// Initialize tags map if it doesn't exist
	if fm.Tags == nil {
		fm.Tags = make(map[string]interface{})
	}

	// Track changes for response message
	var changes []string

	// Process each tag
	for tagKey, tagValue := range tagsMap {
		oldValue := fm.Tags[tagKey]

		if tagValue == nil {
			// Remove the tag if value is null
			if oldValue != nil {
				delete(fm.Tags, tagKey)
				changes = append(changes, fmt.Sprintf("'%s' removed (was: %v)", tagKey, oldValue))
			} else {
				changes = append(changes, fmt.Sprintf("'%s' already absent", tagKey))
			}
		} else {
			// Set or update the tag
			fm.Tags[tagKey] = tagValue
			if oldValue != nil {
				changes = append(changes, fmt.Sprintf("'%s' changed from %v to %v", tagKey, oldValue, tagValue))
			} else {
				changes = append(changes, fmt.Sprintf("'%s' set to %v", tagKey, tagValue))
			}
		}
	}

	// Format the updated document
	updatedContent, err := memory.FormatDocument(fm, body)
	if err != nil {
		return nil, fmt.Errorf("failed to format updated document: %w", err)
	}

	// Update the memory using the enhanced store
	if err := s.enhancedStore.Update(name, updatedContent); err != nil {
		return nil, fmt.Errorf("failed to update memory: %w", err)
	}

	// Build response message
	message := fmt.Sprintf("Updated tags in memory '%s':\n- %s", name, strings.Join(changes, "\n- "))

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: message,
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
