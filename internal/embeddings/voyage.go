package embeddings

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// VoyageClient handles communication with VoyageAI API
type VoyageClient struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// NewVoyageClient creates a new VoyageAI client
func NewVoyageClient(apiKey string) *VoyageClient {
	return &VoyageClient{
		apiKey:  apiKey,
		baseURL: "https://api.voyageai.com/v1",
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// EmbedRequest represents a request to the embeddings API
type EmbedRequest struct {
	Input []string `json:"input"`
	Model string   `json:"model"`
}

// EmbedResponse represents a response from the embeddings API
type EmbedResponse struct {
	Data []struct {
		Embedding []float32 `json:"embedding"`
		Index     int       `json:"index"`
	} `json:"data"`
	Model string `json:"model"`
	Usage struct {
		TotalTokens int `json:"total_tokens"`
	} `json:"usage"`
}

// EmbedTexts generates embeddings for a list of texts
func (c *VoyageClient) EmbedTexts(texts []string, model string) ([][]float32, error) {
	log.Printf("[VOYAGE AI] Starting embedding generation for %d texts using model: %s", len(texts), model)
	
	if len(texts) == 0 {
		return nil, fmt.Errorf("no texts provided")
	}

	if model == "" {
		model = "voyage-3.5" // Default model
	}
	
	// Log preview of texts being embedded
	for i, text := range texts {
		preview := text
		if len(preview) > 100 {
			preview = preview[:100] + "..."
		}
		log.Printf("[VOYAGE AI] Text %d preview: %s", i+1, preview)
	}

	// Create request payload
	reqData := EmbedRequest{
		Input: texts,
		Model: model,
	}

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", c.baseURL+"/embeddings", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	// Send request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var embedResp EmbedResponse
	if err := json.Unmarshal(body, &embedResp); err != nil {
		log.Printf("[VOYAGE AI] ERROR: Failed to parse response: %v", err)
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	log.Printf("[VOYAGE AI] Successfully parsed response - got %d embeddings, used %d tokens", 
		len(embedResp.Data), embedResp.Usage.TotalTokens)

	// Extract embeddings in the correct order
	embeddings := make([][]float32, len(texts))
	for _, item := range embedResp.Data {
		if item.Index >= len(embeddings) {
			return nil, fmt.Errorf("invalid embedding index: %d", item.Index)
		}
		embeddings[item.Index] = item.Embedding
		previewLen := 3
		if len(item.Embedding) < previewLen {
			previewLen = len(item.Embedding)
		}
		log.Printf("[VOYAGE AI] Embedding %d: length %d, preview: %v", 
			item.Index, len(item.Embedding), item.Embedding[:previewLen])
	}

	log.Printf("[VOYAGE AI] Successfully generated %d embeddings", len(embeddings))
	return embeddings, nil
}

// EmbedSingle generates an embedding for a single text
func (c *VoyageClient) EmbedSingle(text string, model string) ([]float32, error) {
	embeddings, err := c.EmbedTexts([]string{text}, model)
	if err != nil {
		return nil, err
	}
	
	if len(embeddings) == 0 {
		return nil, fmt.Errorf("no embeddings returned")
	}
	
	return embeddings[0], nil
}

// EmbedChunks generates embeddings for all chunks of text
func (c *VoyageClient) EmbedChunks(chunks []Chunk, model string) ([][]float32, error) {
	if len(chunks) == 0 {
		return nil, fmt.Errorf("no chunks provided")
	}

	// Extract text from chunks
	texts := make([]string, len(chunks))
	for i, chunk := range chunks {
		texts[i] = chunk.Text
	}

	return c.EmbedTexts(texts, model)
}

// BatchEmbedder handles batching of embedding requests to avoid rate limits
type BatchEmbedder struct {
	client    *VoyageClient
	batchSize int
	delay     time.Duration
}

// NewBatchEmbedder creates a new batch embedder
func NewBatchEmbedder(client *VoyageClient, batchSize int, delay time.Duration) *BatchEmbedder {
	if batchSize <= 0 {
		batchSize = 100 // Default batch size
	}
	if delay <= 0 {
		delay = 100 * time.Millisecond // Default delay
	}

	return &BatchEmbedder{
		client:    client,
		batchSize: batchSize,
		delay:     delay,
	}
}

// EmbedAllChunks processes chunks in batches with rate limiting
func (b *BatchEmbedder) EmbedAllChunks(chunks []Chunk, model string) ([][]float32, error) {
	if len(chunks) == 0 {
		return nil, fmt.Errorf("no chunks provided")
	}

	var allEmbeddings [][]float32

	// Process in batches
	for i := 0; i < len(chunks); i += b.batchSize {
		end := i + b.batchSize
		if end > len(chunks) {
			end = len(chunks)
		}

		batch := chunks[i:end]
		embeddings, err := b.client.EmbedChunks(batch, model)
		if err != nil {
			return nil, fmt.Errorf("failed to embed batch starting at %d: %w", i, err)
		}

		allEmbeddings = append(allEmbeddings, embeddings...)

		// Add delay between batches (except for the last batch)
		if end < len(chunks) {
			time.Sleep(b.delay)
		}
	}

	return allEmbeddings, nil
}

// RerankRequest represents a request to the rerank API
type RerankRequest struct {
	Query     string   `json:"query"`
	Documents []string `json:"documents"`
	Model     string   `json:"model"`
	TopK      int      `json:"top_k,omitempty"`
}

// RerankResponse represents a response from the rerank API
type RerankResponse struct {
	Data []struct {
		Document      string  `json:"document"`
		Index         int     `json:"index"`
		RelevanceScore float32 `json:"relevance_score"`
	} `json:"data"`
	Model string `json:"model"`
	Usage struct {
		TotalTokens int `json:"total_tokens"`
	} `json:"usage"`
}

// RerankResult represents a single rerank result
type RerankResult struct {
	Document       string
	OriginalIndex  int
	RelevanceScore float32
}

// RerankDocuments reranks documents based on query relevance
func (c *VoyageClient) RerankDocuments(query string, documents []string, model string, topK int) ([]RerankResult, error) {
	if len(documents) == 0 {
		return nil, fmt.Errorf("no documents provided")
	}

	if model == "" {
		model = "rerank-lite-1" // Default rerank model
	}

	// Create request payload
	reqData := RerankRequest{
		Query:     query,
		Documents: documents,
		Model:     model,
		TopK:      topK,
	}

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", c.baseURL+"/rerank", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	// Send request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var rerankResp RerankResponse
	if err := json.Unmarshal(body, &rerankResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Convert to results
	var results []RerankResult
	for _, item := range rerankResp.Data {
		results = append(results, RerankResult{
			Document:       item.Document,
			OriginalIndex:  item.Index,
			RelevanceScore: item.RelevanceScore,
		})
	}

	return results, nil
}

// ValidateAPIKey tests if the API key is valid
func (c *VoyageClient) ValidateAPIKey() error {
	_, err := c.EmbedSingle("test", "voyage-3.5")
	if err != nil {
		return fmt.Errorf("API key validation failed: %w", err)
	}
	return nil
}