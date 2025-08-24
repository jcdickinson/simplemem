package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/ThinkInAIXYZ/go-mcp/protocol"
	"github.com/ThinkInAIXYZ/go-mcp/server"
	"github.com/ThinkInAIXYZ/go-mcp/transport"
	"github.com/jcdickinson/simplemem/internal/memory"
)

type Server struct {
	mcpServer *server.Server
	store     *memory.Store
}

type createMemoryReq struct {
	Name    string `json:"name" description:"Name of the memory (without .md extension)"`
	Content string `json:"content" description:"Content of the memory in markdown format"`
}

type readMemoryReq struct {
	Name string `json:"name" description:"Name of the memory to read"`
}

type updateMemoryReq struct {
	Name    string `json:"name" description:"Name of the memory to update"`
	Content string `json:"content" description:"New content for the memory"`
}

type deleteMemoryReq struct {
	Name string `json:"name" description:"Name of the memory to delete"`
}

type searchMemoriesReq struct {
	Query string `json:"query" description:"Search query - can search content or use 'tag:tagname' to search by tag"`
}

type getBacklinksReq struct {
	Name string `json:"name" description:"Name of the memory to find backlinks for"`
}

func NewServer() (*Server, error) {
	s := &Server{
		store: memory.NewStore(".memories"),
	}

	// Initialize the store
	if err := s.store.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize store: %w", err)
	}

	// Create MCP server with stdio transport
	mcpServer, err := server.NewServer(
		transport.NewStdioServerTransport(),
		server.WithServerInfo(protocol.Implementation{
			Name:    "simplemem",
			Version: "0.1.0",
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create MCP server: %w", err)
	}

	// Register all tools
	if err := s.registerTools(mcpServer); err != nil {
		return nil, fmt.Errorf("failed to register tools: %w", err)
	}

	s.mcpServer = mcpServer
	return s, nil
}

func (s *Server) registerTools(mcpServer *server.Server) error {
	// Create Memory tool
	createTool, err := protocol.NewTool(
		"create_memory",
		"Create a new memory document. Supports YAML frontmatter for metadata including title, description, tags, and timestamps",
		createMemoryReq{},
	)
	if err != nil {
		return err
	}
	mcpServer.RegisterTool(createTool, s.handleCreateMemory)

	// Read Memory tool
	readTool, err := protocol.NewTool(
		"read_memory",
		"Read a memory document with full metadata including tags, timestamps, and links",
		readMemoryReq{},
	)
	if err != nil {
		return err
	}
	mcpServer.RegisterTool(readTool, s.handleReadMemory)

	// Update Memory tool
	updateTool, err := protocol.NewTool(
		"update_memory",
		"Update an existing memory document",
		updateMemoryReq{},
	)
	if err != nil {
		return err
	}
	mcpServer.RegisterTool(updateTool, s.handleUpdateMemory)

	// Delete Memory tool
	deleteTool, err := protocol.NewTool(
		"delete_memory",
		"Delete a memory document",
		deleteMemoryReq{},
	)
	if err != nil {
		return err
	}
	mcpServer.RegisterTool(deleteTool, s.handleDeleteMemory)

	// List Memories tool
	listTool, err := protocol.NewTool(
		"list_memories",
		"List all memory documents with metadata preview including titles, tags, and modification dates",
		struct{}{},
	)
	if err != nil {
		return err
	}
	mcpServer.RegisterTool(listTool, s.handleListMemories)

	// Search Memories tool
	searchTool, err := protocol.NewTool(
		"search_memories",
		"Search memories by content or tags. Use 'tag:tagname' to search by tag, or 'tag:' to list all tags",
		searchMemoriesReq{},
	)
	if err != nil {
		return err
	}
	mcpServer.RegisterTool(searchTool, s.handleSearchMemories)

	// Get Backlinks tool
	backlinksTool, err := protocol.NewTool(
		"get_backlinks",
		"Get all memories that link to a specific memory",
		getBacklinksReq{},
	)
	if err != nil {
		return err
	}
	mcpServer.RegisterTool(backlinksTool, s.handleGetBacklinks)

	return nil
}

func (s *Server) handleCreateMemory(_ context.Context, request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	req := new(createMemoryReq)
	if err := protocol.VerifyAndUnmarshal(request.RawArguments, &req); err != nil {
		return nil, err
	}

	if err := s.store.Create(req.Name, req.Content); err != nil {
		return nil, err
	}

	return &protocol.CallToolResult{
		Content: []protocol.Content{
			&protocol.TextContent{
				Type: "text",
				Text: fmt.Sprintf("Memory '%s' created successfully", req.Name),
			},
		},
	}, nil
}

func (s *Server) handleReadMemory(_ context.Context, request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	req := new(readMemoryReq)
	if err := protocol.VerifyAndUnmarshal(request.RawArguments, &req); err != nil {
		return nil, err
	}

	memInfo, err := s.store.ReadWithMetadata(req.Name)
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

	return &protocol.CallToolResult{
		Content: []protocol.Content{
			&protocol.TextContent{
				Type: "text",
				Text: response,
			},
		},
	}, nil
}

func (s *Server) handleUpdateMemory(_ context.Context, request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	req := new(updateMemoryReq)
	if err := protocol.VerifyAndUnmarshal(request.RawArguments, &req); err != nil {
		return nil, err
	}

	if err := s.store.Update(req.Name, req.Content); err != nil {
		return nil, err
	}

	return &protocol.CallToolResult{
		Content: []protocol.Content{
			&protocol.TextContent{
				Type: "text",
				Text: fmt.Sprintf("Memory '%s' updated successfully", req.Name),
			},
		},
	}, nil
}

func (s *Server) handleDeleteMemory(_ context.Context, request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	req := new(deleteMemoryReq)
	if err := protocol.VerifyAndUnmarshal(request.RawArguments, &req); err != nil {
		return nil, err
	}

	if err := s.store.Delete(req.Name); err != nil {
		return nil, err
	}

	return &protocol.CallToolResult{
		Content: []protocol.Content{
			&protocol.TextContent{
				Type: "text",
				Text: fmt.Sprintf("Memory '%s' deleted successfully", req.Name),
			},
		},
	}, nil
}

func (s *Server) handleListMemories(_ context.Context, request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	memories, err := s.store.List()
	if err != nil {
		return nil, err
	}

	if len(memories) == 0 {
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				&protocol.TextContent{
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

	return &protocol.CallToolResult{
		Content: []protocol.Content{
			&protocol.TextContent{
				Type: "text",
				Text: result,
			},
		},
	}, nil
}

func (s *Server) handleSearchMemories(_ context.Context, request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	req := new(searchMemoriesReq)
	if err := protocol.VerifyAndUnmarshal(request.RawArguments, &req); err != nil {
		return nil, err
	}

	// Check if this is a tag search
	if strings.HasPrefix(req.Query, "tag:") {
		tagName := strings.TrimPrefix(req.Query, "tag:")
		if tagName == "" {
			// Show all tags
			allTags, err := s.store.GetAllTags()
			if err != nil {
				return nil, err
			}

			if len(allTags) == 0 {
				return &protocol.CallToolResult{
					Content: []protocol.Content{
						&protocol.TextContent{
							Type: "text",
							Text: "No tags found across memories",
						},
					},
				}, nil
			}

			result := "All tags:\n\n"
			for tag, memories := range allTags {
				result += fmt.Sprintf("🏷️ **%s** (%d memories): %s\n", tag, len(memories), strings.Join(memories, ", "))
			}

			return &protocol.CallToolResult{
				Content: []protocol.Content{
					&protocol.TextContent{
						Type: "text",
						Text: result,
					},
				},
			}, nil
		}

		// Search by specific tag
		memories, err := s.store.SearchByTag(tagName)
		if err != nil {
			return nil, err
		}

		if len(memories) == 0 {
			return &protocol.CallToolResult{
				Content: []protocol.Content{
					&protocol.TextContent{
						Type: "text",
						Text: fmt.Sprintf("No memories found with tag '%s'", tagName),
					},
				},
			}, nil
		}

		result := fmt.Sprintf("Memories with tag '%s':\n\n", tagName)
		for _, memory := range memories {
			result += fmt.Sprintf("📄 %s\n", memory)
		}

		return &protocol.CallToolResult{
			Content: []protocol.Content{
				&protocol.TextContent{
					Type: "text",
					Text: result,
				},
			},
		}, nil
	}

	// Regular content search
	results, err := s.store.Search(req.Query)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				&protocol.TextContent{
					Type: "text",
					Text: fmt.Sprintf("No memories found matching '%s'", req.Query),
				},
			},
		}, nil
	}

	result := fmt.Sprintf("Search results for '%s':\n\n", req.Query)
	for memory, matches := range results {
		result += fmt.Sprintf("📄 %s:\n", memory)
		for _, match := range matches {
			result += fmt.Sprintf("  %s\n", match)
		}
		result += "\n"
	}

	return &protocol.CallToolResult{
		Content: []protocol.Content{
			&protocol.TextContent{
				Type: "text",
				Text: result,
			},
		},
	}, nil
}

func (s *Server) handleGetBacklinks(_ context.Context, request *protocol.CallToolRequest) (*protocol.CallToolResult, error) {
	req := new(getBacklinksReq)
	if err := protocol.VerifyAndUnmarshal(request.RawArguments, &req); err != nil {
		return nil, err
	}

	backlinks, err := memory.GetBacklinks(s.store, req.Name)
	if err != nil {
		return nil, err
	}

	if len(backlinks) == 0 {
		return &protocol.CallToolResult{
			Content: []protocol.Content{
				&protocol.TextContent{
					Type: "text",
					Text: fmt.Sprintf("No backlinks found for '%s'", req.Name),
				},
			},
		}, nil
	}

	result := fmt.Sprintf("Memories linking to '%s':\n", req.Name)
	for _, backlink := range backlinks {
		result += fmt.Sprintf("- %s\n", backlink)
	}

	return &protocol.CallToolResult{
		Content: []protocol.Content{
			&protocol.TextContent{
				Type: "text",
				Text: result,
			},
		},
	}, nil
}

func (s *Server) Run() error {
	return s.mcpServer.Run()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.mcpServer.Shutdown(ctx)
}