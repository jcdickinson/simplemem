package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/marcboeker/go-duckdb"
)

type DB struct {
	conn *sql.DB
}

// New creates a new DuckDB connection and initializes the schema
func New(dbPath string) (*DB, error) {
	// Ensure cache directory exists
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	conn, err := sql.Open("duckdb", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db := &DB{conn: conn}
	if err := db.initSchema(); err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return db, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}

// initSchema creates all necessary tables
func (db *DB) initSchema() error {
	queries := []string{
		// Enable vector extension
		`INSTALL vss;`,
		`LOAD vss;`,
		
		// Create sequences for auto-increment IDs
		`CREATE SEQUENCE IF NOT EXISTS seq_memory_id START 1;`,
		`CREATE SEQUENCE IF NOT EXISTS seq_tag_id START 1;`,
		`CREATE SEQUENCE IF NOT EXISTS seq_link_id START 1;`,
		`CREATE SEQUENCE IF NOT EXISTS seq_embedding_id START 1;`,
		`CREATE SEQUENCE IF NOT EXISTS seq_backlink_id START 1;`,
		
		// Memory documents table with tracking column
		`CREATE TABLE IF NOT EXISTS memories (
			id INTEGER PRIMARY KEY,
			name VARCHAR UNIQUE,
			title VARCHAR,
			description TEXT,
			content TEXT,
			body TEXT,
			created TIMESTAMP,
			modified TIMESTAMP,
			last_processed TIMESTAMP,
			file_hash VARCHAR
		)`,

		// Create indexes separately for memories table
		`CREATE INDEX IF NOT EXISTS idx_memories_name ON memories (name)`,

		// Tags table (normalized)
		`CREATE TABLE IF NOT EXISTS tags (
			id INTEGER PRIMARY KEY,
			memory_id INTEGER REFERENCES memories(id),
			tag_name VARCHAR,
			tag_value VARCHAR
		)`,

		// Create indexes separately for tags table
		`CREATE INDEX IF NOT EXISTS idx_tags_name ON tags (tag_name)`,
		`CREATE INDEX IF NOT EXISTS idx_tags_memory_id ON tags (memory_id)`,

		// Links between memories
		`CREATE TABLE IF NOT EXISTS memory_links (
			id INTEGER PRIMARY KEY,
			from_memory_id INTEGER REFERENCES memories(id),
			to_memory_name VARCHAR,
			link_text VARCHAR,
			link_type VARCHAR
		)`,

		// Create indexes separately for memory_links table
		`CREATE INDEX IF NOT EXISTS idx_memory_links_from ON memory_links (from_memory_id)`,
		`CREATE INDEX IF NOT EXISTS idx_memory_links_to ON memory_links (to_memory_name)`,

		// Vector embeddings with chunking support
		`CREATE TABLE IF NOT EXISTS embeddings (
			id INTEGER PRIMARY KEY,
			memory_id INTEGER REFERENCES memories(id),
			chunk_text TEXT,
			chunk_index INTEGER,
			embedding FLOAT[1024],
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,

		// Create index separately for embeddings table
		`CREATE INDEX IF NOT EXISTS idx_embeddings_memory_id ON embeddings (memory_id)`,

		// Bidirectional semantic backlinks
		`CREATE TABLE IF NOT EXISTS semantic_backlinks (
			id INTEGER PRIMARY KEY,
			memory_a_id INTEGER REFERENCES memories(id),
			memory_b_id INTEGER REFERENCES memories(id),
			similarity_score FLOAT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(memory_a_id, memory_b_id)
		)`,

		// Create indexes separately for semantic_backlinks table
		`CREATE INDEX IF NOT EXISTS idx_semantic_backlinks_a ON semantic_backlinks (memory_a_id)`,
		`CREATE INDEX IF NOT EXISTS idx_semantic_backlinks_b ON semantic_backlinks (memory_b_id)`,
		`CREATE INDEX IF NOT EXISTS idx_semantic_backlinks_score ON semantic_backlinks (similarity_score)`,
	}

	for _, query := range queries {
		if _, err := db.conn.Exec(query); err != nil {
			return fmt.Errorf("failed to execute schema query: %s: %w", query, err)
		}
	}

	return nil
}

// Memory represents a memory document in the database
type Memory struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Content       string    `json:"content"`
	Body          string    `json:"body"`
	Created       time.Time `json:"created"`
	Modified      time.Time `json:"modified"`
	LastProcessed *time.Time `json:"last_processed,omitempty"`
	FileHash      string    `json:"file_hash"`
}

// Embedding represents a vector embedding for a memory chunk
type Embedding struct {
	ID        int       `json:"id"`
	MemoryID  int       `json:"memory_id"`
	ChunkText string    `json:"chunk_text"`
	ChunkIndex int      `json:"chunk_index"`
	Embedding []float32 `json:"embedding"`
	CreatedAt time.Time `json:"created_at"`
}

// SemanticBacklink represents a bidirectional semantic relationship
type SemanticBacklink struct {
	ID              int       `json:"id"`
	MemoryAID       int       `json:"memory_a_id"`
	MemoryBID       int       `json:"memory_b_id"`
	SimilarityScore float32   `json:"similarity_score"`
	CreatedAt       time.Time `json:"created_at"`
}

// UpsertMemory inserts or updates a memory document
func (db *DB) UpsertMemory(memory *Memory) error {
	log.Printf("[DB UPSERT] Upserting memory: %s (created: %v, modified: %v)", 
		memory.Name, memory.Created, memory.Modified)
	
	// First try to get existing memory to see if it's an update
	existing, err := db.GetMemory(memory.Name)
	if err != nil {
		log.Printf("[DB UPSERT] ERROR: Failed to check existing memory: %v", err)
		return fmt.Errorf("failed to check existing memory: %w", err)
	}
	
	if existing != nil {
		// Update existing memory
		log.Printf("[DB UPSERT] Updating existing memory (ID: %d)", existing.ID)
		memory.ID = existing.ID
		
		query := `
			UPDATE memories SET 
				title = ?, description = ?, content = ?, body = ?, 
				modified = ?, file_hash = ?
			WHERE id = ?`
		
		_, err = db.conn.Exec(query, memory.Title, memory.Description, 
			memory.Content, memory.Body, memory.Modified, memory.FileHash, memory.ID)
		if err != nil {
			log.Printf("[DB UPSERT] ERROR: Failed to update memory %s: %v", memory.Name, err)
			return fmt.Errorf("failed to update memory: %w", err)
		}
	} else {
		// Insert new memory using sequence
		log.Printf("[DB UPSERT] Inserting new memory")
		
		query := `
			INSERT INTO memories (id, name, title, description, content, body, created, modified, file_hash)
			VALUES (nextval('seq_memory_id'), ?, ?, ?, ?, ?, ?, ?, ?)`
		
		_, err = db.conn.Exec(query, memory.Name, memory.Title, memory.Description, 
			memory.Content, memory.Body, memory.Created, memory.Modified, memory.FileHash)
		if err != nil {
			log.Printf("[DB UPSERT] ERROR: Failed to insert memory %s: %v", memory.Name, err)
			return fmt.Errorf("failed to insert memory: %w", err)
		}
		
		// Get the ID that was just inserted
		err = db.conn.QueryRow("SELECT currval('seq_memory_id')").Scan(&memory.ID)
		if err != nil {
			log.Printf("[DB UPSERT] ERROR: Failed to get inserted memory ID: %v", err)
			return fmt.Errorf("failed to get inserted memory ID: %w", err)
		}
	}

	log.Printf("[DB UPSERT] Successfully upserted memory: %s (ID: %d)", memory.Name, memory.ID)
	return nil
}

// GetMemory retrieves a memory by name
func (db *DB) GetMemory(name string) (*Memory, error) {
	query := `SELECT id, name, title, description, content, body, created, modified, last_processed, file_hash 
		FROM memories WHERE name = ?`
	
	memory := &Memory{}
	err := db.conn.QueryRow(query, name).Scan(&memory.ID, &memory.Name, &memory.Title, 
		&memory.Description, &memory.Content, &memory.Body, &memory.Created, 
		&memory.Modified, &memory.LastProcessed, &memory.FileHash)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get memory: %w", err)
	}

	return memory, nil
}

// GetMemoriesNeedingProcessing returns memories that have been modified since last processing
func (db *DB) GetMemoriesNeedingProcessing() ([]Memory, error) {
	query := `SELECT id, name, title, description, content, body, created, modified, last_processed, file_hash 
		FROM memories WHERE last_processed IS NULL OR modified > last_processed`
	
	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get memories needing processing: %w", err)
	}
	defer rows.Close()

	var memories []Memory
	for rows.Next() {
		var memory Memory
		err := rows.Scan(&memory.ID, &memory.Name, &memory.Title, &memory.Description, 
			&memory.Content, &memory.Body, &memory.Created, &memory.Modified, 
			&memory.LastProcessed, &memory.FileHash)
		if err != nil {
			return nil, fmt.Errorf("failed to scan memory: %w", err)
		}
		memories = append(memories, memory)
	}

	return memories, nil
}

// MarkMemoryProcessed updates the last_processed timestamp for a memory
func (db *DB) MarkMemoryProcessed(memoryID int) error {
	query := `UPDATE memories SET last_processed = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := db.conn.Exec(query, memoryID)
	if err != nil {
		return fmt.Errorf("failed to mark memory as processed: %w", err)
	}
	return nil
}

// InsertEmbedding stores a vector embedding for a memory chunk
func (db *DB) InsertEmbedding(embedding *Embedding) error {
	log.Printf("[DB EMBEDDING] Inserting embedding for memory %d, chunk %d (vector size: %d)", 
		embedding.MemoryID, embedding.ChunkIndex, len(embedding.Embedding))
	
	// Convert []float32 to a format DuckDB can handle
	// For DuckDB with VSS extension, we need to convert to proper array format
	embeddingStr := fmt.Sprintf("[%s]", strings.Join(func() []string {
		strs := make([]string, len(embedding.Embedding))
		for i, v := range embedding.Embedding {
			strs[i] = fmt.Sprintf("%g", v)
		}
		return strs
	}(), ","))
	
	query := `INSERT INTO embeddings (id, memory_id, chunk_text, chunk_index, embedding)
		VALUES (nextval('seq_embedding_id'), ?, ?, ?, ?::FLOAT[1024])`
	
	_, err := db.conn.Exec(query, embedding.MemoryID, embedding.ChunkText, 
		embedding.ChunkIndex, embeddingStr)
	if err != nil {
		log.Printf("[DB EMBEDDING] ERROR: Failed to insert embedding: %v", err)
		return fmt.Errorf("failed to insert embedding: %w", err)
	}

	log.Printf("[DB EMBEDDING] Successfully inserted embedding for memory %d", embedding.MemoryID)
	return nil
}

// DeleteEmbeddingsByMemoryID removes all embeddings for a memory
func (db *DB) DeleteEmbeddingsByMemoryID(memoryID int) error {
	query := `DELETE FROM embeddings WHERE memory_id = ?`
	_, err := db.conn.Exec(query, memoryID)
	if err != nil {
		return fmt.Errorf("failed to delete embeddings: %w", err)
	}
	return nil
}

// DeleteMemory removes a memory and all related data
func (db *DB) DeleteMemory(name string) error {
	// Get memory ID first
	memory, err := db.GetMemory(name)
	if err != nil {
		return fmt.Errorf("failed to get memory for deletion: %w", err)
	}
	if memory == nil {
		return fmt.Errorf("memory not found: %s", name)
	}

	// Delete related data manually (since we can't use CASCADE)
	queries := []string{
		`DELETE FROM embeddings WHERE memory_id = ?`,
		`DELETE FROM tags WHERE memory_id = ?`,
		`DELETE FROM memory_links WHERE from_memory_id = ?`,
		`DELETE FROM semantic_backlinks WHERE memory_a_id = ? OR memory_b_id = ?`,
		`DELETE FROM memories WHERE id = ?`,
	}

	for i, query := range queries {
		if i == 3 { // semantic_backlinks query needs both parameters
			_, err := db.conn.Exec(query, memory.ID, memory.ID)
			if err != nil {
				return fmt.Errorf("failed to delete semantic backlinks: %w", err)
			}
		} else {
			_, err := db.conn.Exec(query, memory.ID)
			if err != nil {
				return fmt.Errorf("failed to delete related data (step %d): %w", i, err)
			}
		}
	}

	return nil
}

// UpsertTags updates tags for a memory, replacing all existing tags
func (db *DB) UpsertTags(memoryID int, tags map[string]interface{}) error {
	// Delete existing tags for this memory
	_, err := db.conn.Exec(`DELETE FROM tags WHERE memory_id = ?`, memoryID)
	if err != nil {
		return fmt.Errorf("failed to delete existing tags: %w", err)
	}

	// Insert new tags
	if len(tags) > 0 {
		for tagName, tagValue := range tags {
			// Convert tag value to string
			tagValueStr := fmt.Sprintf("%v", tagValue)
			
			_, err := db.conn.Exec(
				`INSERT INTO tags (id, memory_id, tag_name, tag_value) VALUES (nextval('seq_tag_id'), ?, ?, ?)`,
				memoryID, tagName, tagValueStr,
			)
			if err != nil {
				return fmt.Errorf("failed to insert tag %s: %w", tagName, err)
			}
		}
	}

	return nil
}

// FindSimilarMemories finds memories similar to the given embedding vector
func (db *DB) FindSimilarMemories(embedding []float32, threshold float32, limit int, excludeMemoryID int) ([]struct {
	Memory     Memory
	Similarity float32
}, error) {
	log.Printf("[DB VECTOR SEARCH] Starting vector search - embedding_len: %d, threshold: %.3f, limit: %d, exclude_id: %d", 
		len(embedding), threshold, limit, excludeMemoryID)
	
	// Check if we have any embeddings at all
	var embeddingCount int
	err := db.conn.QueryRow("SELECT COUNT(*) FROM embeddings").Scan(&embeddingCount)
	if err != nil {
		log.Printf("[DB VECTOR SEARCH] ERROR: Failed to count embeddings: %v", err)
	} else {
		log.Printf("[DB VECTOR SEARCH] Total embeddings in database: %d", embeddingCount)
	}
	
	// Convert embedding to DuckDB array format
	embeddingStr := fmt.Sprintf("[%s]", strings.Join(func() []string {
		strs := make([]string, len(embedding))
		for i, v := range embedding {
			strs[i] = fmt.Sprintf("%g", v)
		}
		return strs
	}(), ","))

	// Embed the array directly in the query since DuckDB can't handle array parameters properly
	query := fmt.Sprintf(`
		SELECT DISTINCT m.id, m.name, m.title, m.description, m.content, m.body, 
		       m.created, m.modified, m.last_processed, m.file_hash,
		       1 - (e.embedding <=> %s) as similarity
		FROM memories m
		JOIN embeddings e ON m.id = e.memory_id
		WHERE m.id != ? AND (1 - (e.embedding <=> %s)) > ?
		ORDER BY similarity DESC
		LIMIT ?`, embeddingStr, embeddingStr)
	
	log.Printf("[DB VECTOR SEARCH] Executing query with cosine similarity calculation")
	rows, err := db.conn.Query(query, excludeMemoryID, threshold, limit)
	if err != nil {
		log.Printf("[DB VECTOR SEARCH] ERROR: Query failed: %v", err)
		return nil, fmt.Errorf("failed to find similar memories: %w", err)
	}
	defer rows.Close()

	var results []struct {
		Memory     Memory
		Similarity float32
	}

	resultCount := 0
	for rows.Next() {
		var result struct {
			Memory     Memory
			Similarity float32
		}
		
		err := rows.Scan(&result.Memory.ID, &result.Memory.Name, &result.Memory.Title,
			&result.Memory.Description, &result.Memory.Content, &result.Memory.Body,
			&result.Memory.Created, &result.Memory.Modified, &result.Memory.LastProcessed,
			&result.Memory.FileHash, &result.Similarity)
		if err != nil {
			log.Printf("[DB VECTOR SEARCH] ERROR: Failed to scan row: %v", err)
			return nil, fmt.Errorf("failed to scan similar memory: %w", err)
		}
		
		resultCount++
		log.Printf("[DB VECTOR SEARCH] Found result %d: '%s' (ID: %d, similarity: %.4f)", 
			resultCount, result.Memory.Name, result.Memory.ID, result.Similarity)
		
		results = append(results, result)
	}

	log.Printf("[DB VECTOR SEARCH] Query completed - returned %d results", len(results))
	return results, nil
}

// UpsertSemanticBacklink creates or updates a bidirectional semantic relationship
func (db *DB) UpsertSemanticBacklink(memoryAID, memoryBID int, similarity float32) error {
	// Ensure consistent ordering (smaller ID first) for bidirectional relationship
	if memoryAID > memoryBID {
		memoryAID, memoryBID = memoryBID, memoryAID
	}

	query := `
		INSERT INTO semantic_backlinks (id, memory_a_id, memory_b_id, similarity_score)
		VALUES (nextval('seq_backlink_id'), ?, ?, ?)
		ON CONFLICT (memory_a_id, memory_b_id) DO UPDATE SET
			similarity_score = excluded.similarity_score`

	_, err := db.conn.Exec(query, memoryAID, memoryBID, similarity)
	if err != nil {
		return fmt.Errorf("failed to upsert semantic backlink: %w", err)
	}

	return nil
}

// GetMemoryByID retrieves a memory by its ID
func (db *DB) GetMemoryByID(memoryID int) (*Memory, error) {
	query := `SELECT id, name, title, description, content, body, created, modified, last_processed, file_hash 
		FROM memories WHERE id = ?`
	
	memory := &Memory{}
	err := db.conn.QueryRow(query, memoryID).Scan(&memory.ID, &memory.Name, &memory.Title, 
		&memory.Description, &memory.Content, &memory.Body, &memory.Created, 
		&memory.Modified, &memory.LastProcessed, &memory.FileHash)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get memory by ID: %w", err)
	}

	return memory, nil
}

// TagFilter represents a tag filter for searches
type TagFilter struct {
	Key        string
	Value      string
	CheckValue bool // If false, only check for tag presence
}

// FindSimilarMemoriesWithTags finds memories similar to the given embedding vector, filtered by tags
func (db *DB) FindSimilarMemoriesWithTags(embedding []float32, threshold float32, limit int, excludeMemoryID int, tagFilters []TagFilter, requireAll bool) ([]struct {
	Memory     Memory
	Similarity float32
}, error) {
	// Build the base query with tag filtering
	var tagConditions []string
	var tagParams []interface{}
	
	for i, filter := range tagFilters {
		if filter.CheckValue && filter.Value != "" {
			tagConditions = append(tagConditions, fmt.Sprintf("EXISTS (SELECT 1 FROM tags t%d WHERE t%d.memory_id = m.id AND t%d.tag_name = ? AND t%d.tag_value = ?)", i, i, i, i))
			tagParams = append(tagParams, filter.Key, filter.Value)
		} else {
			tagConditions = append(tagConditions, fmt.Sprintf("EXISTS (SELECT 1 FROM tags t%d WHERE t%d.memory_id = m.id AND t%d.tag_name = ?)", i, i, i))
			tagParams = append(tagParams, filter.Key)
		}
	}

	var tagWhereClause string
	if len(tagConditions) > 0 {
		connector := " OR "
		if requireAll {
			connector = " AND "
		}
		tagWhereClause = " AND (" + strings.Join(tagConditions, connector) + ")"
	}

	// Convert embedding to DuckDB array format
	embeddingStr := fmt.Sprintf("[%s]", strings.Join(func() []string {
		strs := make([]string, len(embedding))
		for i, v := range embedding {
			strs[i] = fmt.Sprintf("%g", v)
		}
		return strs
	}(), ","))

	// Embed the array directly in the query since DuckDB can't handle array parameters properly
	query := fmt.Sprintf(`
		SELECT DISTINCT m.id, m.name, m.title, m.description, m.content, m.body, 
		       m.created, m.modified, m.last_processed, m.file_hash,
		       1 - (e.embedding <=> %s) as similarity
		FROM memories m
		JOIN embeddings e ON m.id = e.memory_id
		WHERE m.id != ? AND (1 - (e.embedding <=> %s)) > ?%s
		ORDER BY similarity DESC
		LIMIT ?`, embeddingStr, embeddingStr, tagWhereClause)
	
	params := []interface{}{excludeMemoryID, threshold}
	params = append(params, tagParams...)
	params = append(params, limit)

	rows, err := db.conn.Query(query, params...)
	if err != nil {
		return nil, fmt.Errorf("failed to find similar memories with tags: %w", err)
	}
	defer rows.Close()

	var results []struct {
		Memory     Memory
		Similarity float32
	}

	for rows.Next() {
		var result struct {
			Memory     Memory
			Similarity float32
		}
		
		err := rows.Scan(&result.Memory.ID, &result.Memory.Name, &result.Memory.Title,
			&result.Memory.Description, &result.Memory.Content, &result.Memory.Body,
			&result.Memory.Created, &result.Memory.Modified, &result.Memory.LastProcessed,
			&result.Memory.FileHash, &result.Similarity)
		if err != nil {
			return nil, fmt.Errorf("failed to scan similar memory with tags: %w", err)
		}
		
		results = append(results, result)
	}

	return results, nil
}

// GetMemoriesByTags retrieves memories filtered by tags (for non-semantic searches)
func (db *DB) GetMemoriesByTags(tagFilters []TagFilter, requireAll bool, limit int) ([]Memory, error) {
	// Build tag filtering conditions
	var tagConditions []string
	var params []interface{}
	
	for i, filter := range tagFilters {
		if filter.CheckValue && filter.Value != "" {
			tagConditions = append(tagConditions, fmt.Sprintf("EXISTS (SELECT 1 FROM tags t%d WHERE t%d.memory_id = m.id AND t%d.tag_name = ? AND t%d.tag_value = ?)", i, i, i, i))
			params = append(params, filter.Key, filter.Value)
		} else {
			tagConditions = append(tagConditions, fmt.Sprintf("EXISTS (SELECT 1 FROM tags t%d WHERE t%d.memory_id = m.id AND t%d.tag_name = ?)", i, i, i))
			params = append(params, filter.Key)
		}
	}

	var whereClause string
	if len(tagConditions) > 0 {
		connector := " OR "
		if requireAll {
			connector = " AND "
		}
		whereClause = " WHERE " + strings.Join(tagConditions, connector)
	}

	query := fmt.Sprintf(`
		SELECT DISTINCT m.id, m.name, m.title, m.description, m.content, m.body, 
		       m.created, m.modified, m.last_processed, m.file_hash
		FROM memories m%s
		ORDER BY m.modified DESC
		LIMIT ?`, whereClause)
	
	params = append(params, limit)

	rows, err := db.conn.Query(query, params...)
	if err != nil {
		return nil, fmt.Errorf("failed to get memories by tags: %w", err)
	}
	defer rows.Close()

	var memories []Memory
	for rows.Next() {
		var memory Memory
		err := rows.Scan(&memory.ID, &memory.Name, &memory.Title, &memory.Description, 
			&memory.Content, &memory.Body, &memory.Created, &memory.Modified, 
			&memory.LastProcessed, &memory.FileHash)
		if err != nil {
			return nil, fmt.Errorf("failed to scan memory by tags: %w", err)
		}
		memories = append(memories, memory)
	}

	return memories, nil
}

// GetSemanticBacklinks retrieves semantic backlinks for a memory
func (db *DB) GetSemanticBacklinks(memoryID int, minSimilarity float32) ([]SemanticBacklink, error) {
	query := `
		SELECT id, memory_a_id, memory_b_id, similarity_score, created_at
		FROM semantic_backlinks
		WHERE (memory_a_id = ? OR memory_b_id = ?) AND similarity_score >= ?
		ORDER BY similarity_score DESC`

	rows, err := db.conn.Query(query, memoryID, memoryID, minSimilarity)
	if err != nil {
		return nil, fmt.Errorf("failed to get semantic backlinks: %w", err)
	}
	defer rows.Close()

	var backlinks []SemanticBacklink
	for rows.Next() {
		var backlink SemanticBacklink
		err := rows.Scan(&backlink.ID, &backlink.MemoryAID, &backlink.MemoryBID,
			&backlink.SimilarityScore, &backlink.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan semantic backlink: %w", err)
		}
		backlinks = append(backlinks, backlink)
	}

	return backlinks, nil
}