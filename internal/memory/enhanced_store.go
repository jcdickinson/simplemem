package memory

import (
	"crypto/sha256"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/jcdickinson/simplemem/internal/config"
	"github.com/jcdickinson/simplemem/internal/db"
	"github.com/jcdickinson/simplemem/internal/rag"
)

// EnhancedStore wraps the basic Store with RAG capabilities
type EnhancedStore struct {
	*Store      // Embed the basic store
	db          *db.DB
	ragProcessor *rag.Processor
	dbPath      string
}

// NewEnhancedStore creates a new enhanced store with RAG capabilities
func NewEnhancedStore(basePath string, cfg *config.Config) (*EnhancedStore, error) {
	return NewEnhancedStoreWithDBPath(basePath, cfg, filepath.Join(".cache", "simplemem.db"))
}

// NewEnhancedStoreWithDBPath creates a new enhanced store with a custom database path
func NewEnhancedStoreWithDBPath(basePath string, cfg *config.Config, dbPath string) (*EnhancedStore, error) {
	// Create basic store
	basicStore := NewStore(basePath)

	// Initialize database
	database, err := db.New(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Create RAG processor
	ragProcessor, err := rag.NewProcessor(database, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize RAG processor: %w", err)
	}

	return &EnhancedStore{
		Store:        basicStore,
		db:           database,
		ragProcessor: ragProcessor,
		dbPath:       dbPath,
	}, nil
}

// Initialize initializes both the file store and database
func (es *EnhancedStore) Initialize() error {
	if err := es.Store.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize file store: %w", err)
	}

	// Validate RAG configuration
	if err := es.ragProcessor.ValidateConfiguration(); err != nil {
		log.Printf("Warning: RAG configuration validation failed: %v", err)
		log.Printf("RAG features will be disabled")
	}

	// Sync existing files to database
	if err := es.syncFilesToDatabase(); err != nil {
		log.Printf("Warning: failed to sync files to database: %v", err)
	}

	// Process any pending memories
	if err := es.ragProcessor.ProcessAllPendingMemories(); err != nil {
		log.Printf("Warning: failed to process pending memories: %v", err)
	}

	return nil
}

// Create creates a new memory and processes it with RAG
func (es *EnhancedStore) Create(name, content string) error {
	// Create the file using the basic store
	if err := es.Store.Create(name, content); err != nil {
		return err
	}

	// Sync to database and process with RAG
	if err := es.syncMemoryToDatabase(name); err != nil {
		log.Printf("Warning: failed to sync memory to database: %v", err)
	}

	return nil
}

// Update updates a memory and reprocesses it with RAG
func (es *EnhancedStore) Update(name, content string) error {
	// Update the file using the basic store
	if err := es.Store.Update(name, content); err != nil {
		return err
	}

	// Sync to database and process with RAG
	if err := es.syncMemoryToDatabase(name); err != nil {
		log.Printf("Warning: failed to sync memory to database: %v", err)
	}

	return nil
}

// Delete deletes a memory from both file system and database
func (es *EnhancedStore) Delete(name string) error {
	// Delete from file system first
	if err := es.Store.Delete(name); err != nil {
		return err
	}

	// Clean up database entries (includes all related data)
	if err := es.db.DeleteMemory(name); err != nil {
		log.Printf("Warning: failed to delete memory from database: %v", err)
		// Don't fail the operation if database cleanup fails
	}

	return nil
}

// SearchSemantic performs semantic search using embeddings
func (es *EnhancedStore) SearchSemantic(query string, limit int) ([]MemoryInfo, []float32, error) {
	return es.SearchSemanticWithTags(query, nil, false, limit)
}

// SearchSemanticWithTags performs semantic search using embeddings with tag filtering
func (es *EnhancedStore) SearchSemanticWithTags(query string, tagFilters map[string]string, requireAll bool, limit int) ([]MemoryInfo, []float32, error) {
	// Convert tag filters to db format
	var dbTagFilters []db.TagFilter
	for key, value := range tagFilters {
		filter := db.TagFilter{
			Key:        key,
			Value:      value,
			CheckValue: value != "", // If value is empty, only check presence
		}
		dbTagFilters = append(dbTagFilters, filter)
	}

	memories, similarities, err := es.ragProcessor.SearchSimilarMemoriesWithTags(query, dbTagFilters, requireAll, limit)
	if err != nil {
		return nil, nil, fmt.Errorf("semantic search failed: %w", err)
	}

	// Convert db.Memory to MemoryInfo
	var results []MemoryInfo
	for _, memory := range memories {
		// Parse frontmatter from content
		fm, body, err := ParseDocument(memory.Content)
		if err != nil {
			log.Printf("Warning: failed to parse document %s: %v", memory.Name, err)
			fm = &Frontmatter{}
			body = memory.Content
		}

		results = append(results, MemoryInfo{
			Name:        memory.Name,
			Content:     memory.Content,
			Body:        body,
			Frontmatter: fm,
		})
	}

	return results, similarities, nil
}

// GetSemanticBacklinks returns memories that are semantically similar to the given memory
func (es *EnhancedStore) GetSemanticBacklinks(name string, minSimilarity float32) ([]MemoryInfo, []float32, error) {
	memories, similarities, err := es.ragProcessor.GetSemanticBacklinks(name, minSimilarity)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get semantic backlinks: %w", err)
	}

	// Convert db.Memory to MemoryInfo
	var results []MemoryInfo
	for _, memory := range memories {
		// Parse frontmatter from content
		fm, body, err := ParseDocument(memory.Content)
		if err != nil {
			log.Printf("Warning: failed to parse document %s: %v", memory.Name, err)
			fm = &Frontmatter{}
			body = memory.Content
		}

		results = append(results, MemoryInfo{
			Name:        memory.Name,
			Content:     memory.Content,
			Body:        body,
			Frontmatter: fm,
		})
	}

	return results, similarities, nil
}


// matchesTagFilters checks if a memory matches the given tag filters
func (es *EnhancedStore) matchesTagFilters(memory *MemoryInfo, tagFilters map[string]string, requireAll bool) bool {
	if len(tagFilters) == 0 {
		return true
	}

	matchCount := 0
	for key, value := range tagFilters {
		if memoryValue, exists := memory.Frontmatter.Tags[key]; exists {
			if value == "" {
				// Just check for presence
				matchCount++
			} else {
				// Check for specific value
				if fmt.Sprintf("%v", memoryValue) == value {
					matchCount++
				}
			}
		}
	}

	if requireAll {
		return matchCount == len(tagFilters)
	}
	return matchCount > 0
}

// GetEnhancedBacklinks retrieves and reranks both explicit and semantic backlinks as markdown
func (es *EnhancedStore) GetEnhancedBacklinks(memoryName string, query string, limit int) (string, error) {
	backlinks, err := es.ragProcessor.GetEnhancedBacklinks(memoryName, query, limit)
	if err != nil {
		return "", fmt.Errorf("failed to get enhanced backlinks: %w", err)
	}

	return es.ragProcessor.FormatBacklinksAsMarkdown(memoryName, backlinks, query), nil
}

// SearchSemanticMarkdown performs semantic search and returns results as markdown
func (es *EnhancedStore) SearchSemanticMarkdown(query string, limit int) (string, error) {
	return es.SearchSemanticMarkdownWithTags(query, nil, false, limit)
}

// SearchSemanticMarkdownWithTags performs semantic search with tag filtering and returns results as markdown
func (es *EnhancedStore) SearchSemanticMarkdownWithTags(query string, tagFilters map[string]string, requireAll bool, limit int) (string, error) {
	memories, similarities, err := es.SearchSemanticWithTags(query, tagFilters, requireAll, limit)
	if err != nil {
		return "", err
	}

	if len(memories) == 0 {
		searchDesc := fmt.Sprintf("'%s'", query)
		if len(tagFilters) > 0 {
			var tagDesc []string
			for key, value := range tagFilters {
				if value == "" {
					tagDesc = append(tagDesc, key)
				} else {
					tagDesc = append(tagDesc, fmt.Sprintf("%s:%s", key, value))
				}
			}
			connector := "any of"
			if requireAll {
				connector = "all of"
			}
			searchDesc += fmt.Sprintf(" with %s tags [%s]", connector, strings.Join(tagDesc, ", "))
		}
		return fmt.Sprintf("No memories found for semantic search: %s", searchDesc), nil
	}

	var md strings.Builder
	if len(tagFilters) > 0 {
		var tagDesc []string
		for key, value := range tagFilters {
			if value == "" {
				tagDesc = append(tagDesc, key)
			} else {
				tagDesc = append(tagDesc, fmt.Sprintf("%s:%s", key, value))
			}
		}
		connector := "any of"
		if requireAll {
			connector = "all of"
		}
		md.WriteString(fmt.Sprintf("# Semantic search results for '%s' with %s tags [%s]\n\n", query, connector, strings.Join(tagDesc, ", ")))
	} else {
		md.WriteString(fmt.Sprintf("# Semantic search results for '%s'\n\n", query))
	}

	for i, memory := range memories {
		md.WriteString(fmt.Sprintf("## %d. %s\n", i+1, memory.Name))
		
		if memory.Frontmatter.Title != "" && memory.Frontmatter.Title != memory.Name {
			md.WriteString(fmt.Sprintf("**Title:** %s\n\n", memory.Frontmatter.Title))
		}
		
		// Create snippet from body (first 300 chars)
		snippet := memory.Body
		if len(snippet) > 300 {
			snippet = snippet[:300] + "..."
		}
		md.WriteString(fmt.Sprintf("**Snippet:**\n%s\n\n", snippet))
		
		md.WriteString(fmt.Sprintf("**Similarity:** %.3f\n\n", similarities[i]))
		
		if memory.Frontmatter.Description != "" {
			md.WriteString(fmt.Sprintf("**Description:** %s\n\n", memory.Frontmatter.Description))
		}
		
		if len(memory.Frontmatter.Tags) > 0 {
			var tagsList []string
			for tag, value := range memory.Frontmatter.Tags {
				if value == true {
					tagsList = append(tagsList, tag)
				} else {
					tagsList = append(tagsList, fmt.Sprintf("%s: %v", tag, value))
				}
			}
			md.WriteString(fmt.Sprintf("**Tags:** %s\n\n", strings.Join(tagsList, ", ")))
		}
		
		md.WriteString("---\n\n")
	}

	return md.String(), nil
}


// Close closes database connections
func (es *EnhancedStore) Close() error {
	return es.db.Close()
}

// syncFilesToDatabase syncs all existing files to the database
func (es *EnhancedStore) syncFilesToDatabase() error {
	names, err := es.Store.List()
	if err != nil {
		return fmt.Errorf("failed to list memories: %w", err)
	}

	for _, name := range names {
		if err := es.syncMemoryToDatabase(name); err != nil {
			log.Printf("Warning: failed to sync memory %s: %v", name, err)
		}
	}

	return nil
}

// syncMemoryToDatabase syncs a single memory to the database
func (es *EnhancedStore) syncMemoryToDatabase(name string) error {
	memInfo, err := es.Store.ReadWithMetadata(name)
	if err != nil {
		return fmt.Errorf("failed to read memory: %w", err)
	}

	// Calculate file hash for change detection
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(memInfo.Content)))

	// Check if memory exists in database and if it has changed
	existing, err := es.db.GetMemory(name)
	if err != nil {
		return fmt.Errorf("failed to check existing memory: %w", err)
	}

	if existing != nil && existing.FileHash == hash {
		// Memory hasn't changed, no need to update
		return nil
	}

	// Create or update memory in database
	dbMemory := &db.Memory{
		Name:        name,
		Title:       memInfo.Frontmatter.Title,
		Description: memInfo.Frontmatter.Description,
		Content:     memInfo.Content,
		Body:        memInfo.Body,
		Created:     memInfo.Frontmatter.Created,
		Modified:    memInfo.Frontmatter.Modified,
		FileHash:    hash,
	}

	if existing != nil {
		dbMemory.ID = existing.ID
	}

	if err := es.db.UpsertMemory(dbMemory); err != nil {
		return fmt.Errorf("failed to upsert memory: %w", err)
	}

	// Sync tags to database
	if err := es.db.UpsertTags(dbMemory.ID, memInfo.Frontmatter.Tags); err != nil {
		log.Printf("Warning: failed to sync tags for memory %s: %v", name, err)
	}

	// Process with RAG if content changed
	if err := es.ragProcessor.ProcessMemory(dbMemory); err != nil {
		log.Printf("Warning: failed to process memory %s with RAG: %v", name, err)
	}

	return nil
}

// min returns the minimum of two float32 values
func min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}