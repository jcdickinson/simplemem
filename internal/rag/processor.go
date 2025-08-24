package rag

import (
	"crypto/sha256"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jcdickinson/simplemem/internal/config"
	"github.com/jcdickinson/simplemem/internal/db"
	"github.com/jcdickinson/simplemem/internal/embeddings"
)

// Processor handles RAG operations for memories
type Processor struct {
	db              *db.DB
	voyageClient    *embeddings.VoyageClient
	batchEmbedder   *embeddings.BatchEmbedder
	chunkConfig     embeddings.ChunkConfig
	model           string
	rerankModel     string
	similarityThreshold float32
}

// NewProcessor creates a new RAG processor
func NewProcessor(database *db.DB, cfg *config.Config) (*Processor, error) {
	if cfg.VoyageAI.ApiKey.Value == "" {
		return nil, fmt.Errorf("VoyageAI API key is required")
	}

	voyageClient := embeddings.NewVoyageClient(cfg.VoyageAI.ApiKey.Value)
	batchEmbedder := embeddings.NewBatchEmbedder(voyageClient, 50, 200*time.Millisecond)

	// Use configured models or defaults
	model := cfg.VoyageAI.Model
	if model == "" {
		model = "voyage-3.5"
	}

	processor := &Processor{
		db:                  database,
		voyageClient:        voyageClient,
		batchEmbedder:       batchEmbedder,
		chunkConfig:         embeddings.DefaultChunkConfig(),
		model:               model,
		similarityThreshold: 0.5, // Minimum similarity for semantic backlinks (lowered from 0.7)
	}

	// Store rerank model for later use
	processor.rerankModel = cfg.VoyageAI.RerankModel
	if processor.rerankModel == "" {
		processor.rerankModel = "rerank-lite-1"
	}

	return processor, nil
}

// ProcessMemory handles the complete RAG workflow for a memory
func (p *Processor) ProcessMemory(memory *db.Memory) error {
	log.Printf("Processing memory: %s", memory.Name)

	// 1. Delete existing embeddings for this memory
	if err := p.db.DeleteEmbeddingsByMemoryID(memory.ID); err != nil {
		return fmt.Errorf("failed to delete existing embeddings: %w", err)
	}

	// 2. Chunk the content
	chunks := embeddings.ChunkMarkdown(memory.Body, p.chunkConfig)
	if len(chunks) == 0 {
		log.Printf("No chunks generated for memory: %s", memory.Name)
		return p.db.MarkMemoryProcessed(memory.ID)
	}

	log.Printf("Generated %d chunks for memory: %s", len(chunks), memory.Name)

	// 3. Generate embeddings for all chunks
	chunkEmbeddings, err := p.batchEmbedder.EmbedAllChunks(chunks, p.model)
	if err != nil {
		return fmt.Errorf("failed to generate embeddings: %w", err)
	}

	// 4. Store embeddings in database
	for i, chunk := range chunks {
		embedding := &db.Embedding{
			MemoryID:   memory.ID,
			ChunkText:  chunk.Text,
			ChunkIndex: chunk.Index,
			Embedding:  chunkEmbeddings[i],
		}

		if err := p.db.InsertEmbedding(embedding); err != nil {
			return fmt.Errorf("failed to insert embedding %d: %w", i, err)
		}
	}

	// 5. Find similar memories using the first chunk's embedding
	// (We use the first chunk as it's typically the most representative)
	if len(chunkEmbeddings) > 0 {
		if err := p.updateSemanticBacklinks(memory.ID, chunkEmbeddings[0]); err != nil {
			log.Printf("Failed to update semantic backlinks for %s: %v", memory.Name, err)
			// Don't fail the entire process if backlinks fail
		}
	}

	// 6. Mark memory as processed
	if err := p.db.MarkMemoryProcessed(memory.ID); err != nil {
		return fmt.Errorf("failed to mark memory as processed: %w", err)
	}

	log.Printf("Successfully processed memory: %s", memory.Name)
	return nil
}

// updateSemanticBacklinks finds similar memories and creates bidirectional links
func (p *Processor) updateSemanticBacklinks(memoryID int, representativeEmbedding []float32) error {
	// Find similar memories
	similarMemories, err := p.db.FindSimilarMemories(
		representativeEmbedding,
		p.similarityThreshold,
		20, // Limit to top 20 similar memories
		memoryID,
	)
	if err != nil {
		return fmt.Errorf("failed to find similar memories: %w", err)
	}

	log.Printf("Found %d similar memories for memory ID %d", len(similarMemories), memoryID)

	// Create semantic backlinks
	for _, similar := range similarMemories {
		if err := p.db.UpsertSemanticBacklink(memoryID, similar.Memory.ID, similar.Similarity); err != nil {
			log.Printf("Failed to create semantic backlink between %d and %d: %v", 
				memoryID, similar.Memory.ID, err)
			// Continue with other backlinks even if one fails
		} else {
			log.Printf("Created semantic backlink between %d and %d (similarity: %.3f)", 
				memoryID, similar.Memory.ID, similar.Similarity)
		}
	}

	return nil
}

// ProcessAllPendingMemories processes all memories that need embeddings
func (p *Processor) ProcessAllPendingMemories() error {
	memories, err := p.db.GetMemoriesNeedingProcessing()
	if err != nil {
		return fmt.Errorf("failed to get memories needing processing: %w", err)
	}

	if len(memories) == 0 {
		log.Println("No memories need processing")
		return nil
	}

	log.Printf("Processing %d memories that need embeddings", len(memories))

	for _, memory := range memories {
		if err := p.ProcessMemory(&memory); err != nil {
			log.Printf("Failed to process memory %s: %v", memory.Name, err)
			// Continue with other memories even if one fails
		}
	}

	return nil
}

// SearchSimilarMemories performs semantic search using embeddings
func (p *Processor) SearchSimilarMemories(query string, limit int) ([]db.Memory, []float32, error) {
	return p.SearchSimilarMemoriesWithTags(query, nil, false, limit)
}

// SearchSimilarMemoriesWithTags performs semantic search using embeddings with tag filtering
func (p *Processor) SearchSimilarMemoriesWithTags(query string, tagFilters []db.TagFilter, requireAll bool, limit int) ([]db.Memory, []float32, error) {
	log.Printf("[SEMANTIC SEARCH] Starting search - Query: '%s', TagFilters: %+v, RequireAll: %v, Limit: %d", query, tagFilters, requireAll, limit)
	
	// If only tag filtering (no semantic search), use direct tag search
	if query == "" && len(tagFilters) > 0 {
		log.Printf("[SEMANTIC SEARCH] Empty query with tag filters - using direct tag search")
		memories, err := p.db.GetMemoriesByTags(tagFilters, requireAll, limit)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get memories by tags: %w", err)
		}
		
		log.Printf("[SEMANTIC SEARCH] Tag-only search returned %d memories", len(memories))
		
		// Return with neutral similarity scores
		similarities := make([]float32, len(memories))
		for i := range similarities {
			similarities[i] = 1.0
		}
		
		return memories, similarities, nil
	}

	log.Printf("[SEMANTIC SEARCH] Generating embedding for query using model: %s", p.model)
	// Generate embedding for the search query
	queryEmbedding, err := p.voyageClient.EmbedSingle(query, p.model)
	if err != nil {
		log.Printf("[SEMANTIC SEARCH] ERROR: Failed to generate query embedding: %v", err)
		return nil, nil, fmt.Errorf("failed to generate query embedding: %w", err)
	}
	
	log.Printf("[SEMANTIC SEARCH] Successfully generated embedding vector of length: %d", len(queryEmbedding))
	if len(queryEmbedding) > 0 {
		previewLen := 5
		if len(queryEmbedding) < previewLen {
			previewLen = len(queryEmbedding)
		}
		log.Printf("[SEMANTIC SEARCH] Embedding preview - first %d values: %v", previewLen, queryEmbedding[:previewLen])
	}

	// Find similar memories with tag filtering
	var similarMemories []struct {
		Memory     db.Memory
		Similarity float32
	}
	
	threshold := float32(0.1)
	log.Printf("[SEMANTIC SEARCH] Searching with threshold: %.3f, limit: %d", threshold, limit)
	
	if len(tagFilters) > 0 {
		log.Printf("[SEMANTIC SEARCH] Using tagged search with %d tag filters", len(tagFilters))
		similarMemories, err = p.db.FindSimilarMemoriesWithTags(
			queryEmbedding,
			threshold, // Much lower threshold for search results (was 0.3)
			limit,
			-1, // Don't exclude any memories
			tagFilters,
			requireAll,
		)
	} else {
		log.Printf("[SEMANTIC SEARCH] Using unfiltered semantic search")
		similarMemories, err = p.db.FindSimilarMemories(
			queryEmbedding,
			threshold, // Much lower threshold for search results (was 0.3)
			limit,
			-1, // Don't exclude any memories
		)
	}
	
	if err != nil {
		log.Printf("[SEMANTIC SEARCH] ERROR: Database search failed: %v", err)
		return nil, nil, fmt.Errorf("failed to find similar memories: %w", err)
	}
	
	log.Printf("[SEMANTIC SEARCH] Database returned %d similar memories", len(similarMemories))

	// Extract results
	memories := make([]db.Memory, len(similarMemories))
	similarities := make([]float32, len(similarMemories))
	for i, result := range similarMemories {
		memories[i] = result.Memory
		similarities[i] = result.Similarity
		log.Printf("[SEMANTIC SEARCH] Result %d: '%s' (similarity: %.4f)", i+1, result.Memory.Name, result.Similarity)
	}

	log.Printf("[SEMANTIC SEARCH] Returning %d memories with similarities", len(memories))
	return memories, similarities, nil
}

// GetSemanticBacklinks retrieves memories semantically related to the given memory
func (p *Processor) GetSemanticBacklinks(memoryName string, minSimilarity float32) ([]db.Memory, []float32, error) {
	// Get the memory first
	memory, err := p.db.GetMemory(memoryName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get memory: %w", err)
	}
	if memory == nil {
		return nil, nil, fmt.Errorf("memory not found: %s", memoryName)
	}

	// Get semantic backlinks
	backlinks, err := p.db.GetSemanticBacklinks(memory.ID, minSimilarity)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get semantic backlinks: %w", err)
	}

	// Resolve memory IDs to memory objects
	var memories []db.Memory
	var similarities []float32

	for _, backlink := range backlinks {
		// Determine which memory ID to fetch (not the current one)
		var targetID int
		if backlink.MemoryAID == memory.ID {
			targetID = backlink.MemoryBID
		} else {
			targetID = backlink.MemoryAID
		}

		// This is a simplified approach - in a real implementation, you might want
		// to batch fetch these memories for better performance
		targetMemory, err := p.getMemoryByID(targetID)
		if err != nil {
			log.Printf("Failed to get memory for ID %d: %v", targetID, err)
			continue
		}

		memories = append(memories, *targetMemory)
		similarities = append(similarities, backlink.SimilarityScore)
	}

	return memories, similarities, nil
}

// getMemoryByID is a helper function to get a memory by ID
func (p *Processor) getMemoryByID(memoryID int) (*db.Memory, error) {
	return p.db.GetMemoryByID(memoryID)
}

// calculateHash generates a hash for memory content to detect changes
func calculateHash(content string) string {
	hash := sha256.Sum256([]byte(content))
	return fmt.Sprintf("%x", hash)
}

// BacklinkResult represents a combined backlink result
type BacklinkResult struct {
	Memory         db.Memory
	Snippet        string
	LinkType       string // "explicit", "semantic"
	RelevanceScore float32
	SourceType     string // "wiki", "markdown", "embedding"
}

// GetEnhancedBacklinks retrieves and reranks both explicit and semantic backlinks
func (p *Processor) GetEnhancedBacklinks(memoryName string, query string, limit int) ([]BacklinkResult, error) {
	// Get the target memory
	memory, err := p.db.GetMemory(memoryName)
	if err != nil {
		return nil, fmt.Errorf("failed to get memory: %w", err)
	}
	if memory == nil {
		return nil, fmt.Errorf("memory not found: %s", memoryName)
	}

	var allResults []BacklinkResult

	// 1. Get explicit backlinks (wiki-style and markdown links)
	explicitBacklinks, err := p.getExplicitBacklinks(memoryName)
	if err != nil {
		log.Printf("Warning: failed to get explicit backlinks: %v", err)
	} else {
		allResults = append(allResults, explicitBacklinks...)
	}

	// 2. Get semantic backlinks
	semanticBacklinks, err := p.getSemanticBacklinksAsResults(memory.ID)
	if err != nil {
		log.Printf("Warning: failed to get semantic backlinks: %v", err)
	} else {
		allResults = append(allResults, semanticBacklinks...)
	}

	// 3. If we have a query and results, rerank them
	if query != "" && len(allResults) > 0 {
		allResults, err = p.rerankBacklinks(query, allResults, limit)
		if err != nil {
			log.Printf("Warning: reranking failed, returning original order: %v", err)
		}
	}

	// 4. Limit results if no reranking was performed
	if limit > 0 && len(allResults) > limit {
		allResults = allResults[:limit]
	}

	return allResults, nil
}

// getExplicitBacklinks finds memories that explicitly link to the target memory
func (p *Processor) getExplicitBacklinks(targetMemoryName string) ([]BacklinkResult, error) {
	// This would need to be implemented by scanning all memories for links
	// For now, return empty slice as we'd need to integrate with the memory store
	return []BacklinkResult{}, nil
}

// getSemanticBacklinksAsResults converts semantic backlinks to BacklinkResult format
func (p *Processor) getSemanticBacklinksAsResults(memoryID int) ([]BacklinkResult, error) {
	backlinks, err := p.db.GetSemanticBacklinks(memoryID, 0.1) // Much lower threshold for more results (was 0.3)
	if err != nil {
		return nil, err
	}

	var results []BacklinkResult
	for _, backlink := range backlinks {
		// Determine which memory ID to fetch (not the current one)
		var targetID int
		if backlink.MemoryAID == memoryID {
			targetID = backlink.MemoryBID
		} else {
			targetID = backlink.MemoryAID
		}

		targetMemory, err := p.db.GetMemoryByID(targetID)
		if err != nil {
			log.Printf("Failed to get memory for ID %d: %v", targetID, err)
			continue
		}
		if targetMemory == nil {
			continue
		}

		// Create snippet from the memory content (first 200 chars)
		snippet := targetMemory.Body
		if len(snippet) > 200 {
			snippet = snippet[:200] + "..."
		}

		results = append(results, BacklinkResult{
			Memory:         *targetMemory,
			Snippet:        snippet,
			LinkType:       "semantic",
			RelevanceScore: backlink.SimilarityScore,
			SourceType:     "embedding",
		})
	}

	return results, nil
}

// rerankBacklinks uses VoyageAI reranker to reorder backlinks by relevance to query
func (p *Processor) rerankBacklinks(query string, backlinks []BacklinkResult, topK int) ([]BacklinkResult, error) {
	if len(backlinks) == 0 {
		return backlinks, nil
	}

	// Prepare documents for reranking (use snippets as documents)
	documents := make([]string, len(backlinks))
	for i, backlink := range backlinks {
		// Use title + snippet for better reranking
		doc := backlink.Memory.Title
		if doc != "" {
			doc += "\n\n"
		}
		doc += backlink.Snippet
		documents[i] = doc
	}

	// Rerank using VoyageAI
	rerankResults, err := p.voyageClient.RerankDocuments(query, documents, p.rerankModel, topK)
	if err != nil {
		return nil, fmt.Errorf("failed to rerank: %w", err)
	}

	// Reorder backlinks based on rerank results
	var reorderedBacklinks []BacklinkResult
	for _, result := range rerankResults {
		if result.OriginalIndex < len(backlinks) {
			backlink := backlinks[result.OriginalIndex]
			// Update relevance score with rerank score
			backlink.RelevanceScore = result.RelevanceScore
			reorderedBacklinks = append(reorderedBacklinks, backlink)
		}
	}

	return reorderedBacklinks, nil
}

// FormatBacklinksAsMarkdown formats backlink results as clean markdown
func (p *Processor) FormatBacklinksAsMarkdown(memoryName string, backlinks []BacklinkResult, query string) string {
	if len(backlinks) == 0 {
		return fmt.Sprintf("No backlinks found for '%s'", memoryName)
	}

	var md strings.Builder
	if query != "" {
		md.WriteString(fmt.Sprintf("# Related memories for '%s' (query: %s)\n\n", memoryName, query))
	} else {
		md.WriteString(fmt.Sprintf("# Related memories for '%s'\n\n", memoryName))
	}

	for i, backlink := range backlinks {
		md.WriteString(fmt.Sprintf("## %d. %s\n", i+1, backlink.Memory.Name))
		
		if backlink.Memory.Title != "" && backlink.Memory.Title != backlink.Memory.Name {
			md.WriteString(fmt.Sprintf("**Title:** %s\n\n", backlink.Memory.Title))
		}
		
		md.WriteString(fmt.Sprintf("**Snippet:**\n%s\n\n", backlink.Snippet))
		
		md.WriteString(fmt.Sprintf("**Link Type:** %s | **Relevance:** %.3f\n\n", 
			backlink.LinkType, backlink.RelevanceScore))
		
		if backlink.Memory.Description != "" {
			md.WriteString(fmt.Sprintf("**Description:** %s\n\n", backlink.Memory.Description))
		}
		
		md.WriteString("---\n\n")
	}

	return md.String()
}

// ValidateConfiguration checks if the processor is properly configured
func (p *Processor) ValidateConfiguration() error {
	return p.voyageClient.ValidateAPIKey()
}