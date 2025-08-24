package embeddings

import (
	"regexp"
	"strings"
	"unicode"
)

// ChunkConfig holds configuration for content chunking
type ChunkConfig struct {
	MaxChunkSize    int // Maximum characters per chunk
	OverlapSize     int // Characters to overlap between chunks
	MinChunkSize    int // Minimum characters per chunk (avoid tiny chunks)
}

// DefaultChunkConfig returns sensible defaults for chunking
func DefaultChunkConfig() ChunkConfig {
	return ChunkConfig{
		MaxChunkSize: 1000,  // ~200-300 tokens for most models
		OverlapSize:  100,   // ~20-30 tokens overlap
		MinChunkSize: 100,   // Avoid chunks smaller than this
	}
}

// Chunk represents a piece of content with metadata
type Chunk struct {
	Text  string
	Index int
	Start int // Character position in original text
	End   int // Character position in original text
}

// ChunkText splits text into semantically meaningful chunks
func ChunkText(text string, config ChunkConfig) []Chunk {
	if len(text) <= config.MaxChunkSize {
		return []Chunk{{
			Text:  text,
			Index: 0,
			Start: 0,
			End:   len(text),
		}}
	}

	// Clean and normalize the text
	text = normalizeText(text)
	
	// Split by natural boundaries (paragraphs, sentences, etc.)
	boundaries := findTextBoundaries(text)
	
	return createChunks(text, boundaries, config)
}

// normalizeText cleans up whitespace and formatting
func normalizeText(text string) string {
	// Replace multiple whitespace with single spaces
	re := regexp.MustCompile(`\s+`)
	text = re.ReplaceAllString(text, " ")
	
	// Trim leading/trailing whitespace
	return strings.TrimSpace(text)
}

// findTextBoundaries identifies good split points in the text
func findTextBoundaries(text string) []int {
	var boundaries []int
	
	// Paragraph boundaries (double newlines)
	paragraphRe := regexp.MustCompile(`\n\s*\n`)
	for _, match := range paragraphRe.FindAllStringIndex(text, -1) {
		boundaries = append(boundaries, match[1])
	}
	
	// Sentence boundaries
	sentenceRe := regexp.MustCompile(`[.!?]\s+[A-Z]`)
	for _, match := range sentenceRe.FindAllStringIndex(text, -1) {
		// Add boundary after the punctuation, before the space
		boundaries = append(boundaries, match[0]+1)
	}
	
	// Header boundaries (markdown-style)
	headerRe := regexp.MustCompile(`(?m)^#+\s+`)
	for _, match := range headerRe.FindAllStringIndex(text, -1) {
		boundaries = append(boundaries, match[0])
	}
	
	// List item boundaries
	listRe := regexp.MustCompile(`(?m)^[\s]*[-*+]\s+`)
	for _, match := range listRe.FindAllStringIndex(text, -1) {
		boundaries = append(boundaries, match[0])
	}
	
	// Code block boundaries
	codeRe := regexp.MustCompile("```")
	for _, match := range codeRe.FindAllStringIndex(text, -1) {
		boundaries = append(boundaries, match[0])
	}
	
	// Remove duplicates and sort
	boundaries = removeDuplicatesAndSort(boundaries)
	
	// Always include start and end
	if len(boundaries) == 0 || boundaries[0] != 0 {
		boundaries = append([]int{0}, boundaries...)
	}
	if boundaries[len(boundaries)-1] != len(text) {
		boundaries = append(boundaries, len(text))
	}
	
	return boundaries
}

// createChunks builds chunks from text and boundaries
func createChunks(text string, boundaries []int, config ChunkConfig) []Chunk {
	var chunks []Chunk
	chunkIndex := 0
	
	i := 0
	for i < len(boundaries)-1 {
		chunkStart := boundaries[i]
		chunkEnd := chunkStart
		
		// Find the optimal end position within size limits
		for j := i + 1; j < len(boundaries); j++ {
			if boundaries[j] - chunkStart <= config.MaxChunkSize {
				chunkEnd = boundaries[j]
			} else {
				break
			}
		}
		
		// If we couldn't fit even one boundary, split at max size
		if chunkEnd == chunkStart {
			chunkEnd = min(chunkStart + config.MaxChunkSize, len(text))
		}
		
		chunkText := strings.TrimSpace(text[chunkStart:chunkEnd])
		
		// Skip chunks that are too small (unless it's the last chunk)
		if len(chunkText) >= config.MinChunkSize || chunkEnd == len(text) {
			chunks = append(chunks, Chunk{
				Text:  chunkText,
				Index: chunkIndex,
				Start: chunkStart,
				End:   chunkEnd,
			})
			chunkIndex++
		}
		
		// Move to next chunk with overlap
		nextStart := chunkEnd - config.OverlapSize
		if nextStart <= chunkStart {
			nextStart = chunkEnd
		}
		
		// Find the boundary closest to our desired start position
		for i < len(boundaries)-1 && boundaries[i] < nextStart {
			i++
		}
		
		// If we've reached the end, break
		if boundaries[i] >= len(text) {
			break
		}
	}
	
	return chunks
}

// removeDuplicatesAndSort removes duplicate positions and sorts the slice
func removeDuplicatesAndSort(boundaries []int) []int {
	if len(boundaries) <= 1 {
		return boundaries
	}
	
	// Use a map to remove duplicates
	seen := make(map[int]bool)
	var result []int
	
	for _, boundary := range boundaries {
		if !seen[boundary] {
			seen[boundary] = true
			result = append(result, boundary)
		}
	}
	
	// Simple insertion sort (boundaries array is usually small)
	for i := 1; i < len(result); i++ {
		key := result[i]
		j := i - 1
		for j >= 0 && result[j] > key {
			result[j+1] = result[j]
			j--
		}
		result[j+1] = key
	}
	
	return result
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ChunkMarkdown chunks markdown content while preserving structure
func ChunkMarkdown(content string, config ChunkConfig) []Chunk {
	// For markdown, we want to be more careful about preserving structure
	// This is a specialized version that understands markdown better
	
	chunks := ChunkText(content, config)
	
	// Post-process chunks to ensure markdown integrity
	for i := range chunks {
		chunks[i].Text = ensureMarkdownIntegrity(chunks[i].Text)
	}
	
	return chunks
}

// ensureMarkdownIntegrity cleans up chunk boundaries to avoid broken markdown
func ensureMarkdownIntegrity(text string) string {
	text = strings.TrimSpace(text)
	
	// If chunk starts mid-word, try to find a better break point
	if len(text) > 0 && !unicode.IsSpace(rune(text[0])) && !unicode.IsPunct(rune(text[0])) {
		// Find first space and start from there
		if spaceIdx := strings.IndexFunc(text, unicode.IsSpace); spaceIdx != -1 {
			text = strings.TrimSpace(text[spaceIdx:])
		}
	}
	
	// If chunk ends mid-word, try to find a better break point
	if len(text) > 0 && !unicode.IsSpace(rune(text[len(text)-1])) && !unicode.IsPunct(rune(text[len(text)-1])) {
		// Find last space and end there
		if spaceIdx := strings.LastIndexFunc(text, unicode.IsSpace); spaceIdx != -1 {
			text = strings.TrimSpace(text[:spaceIdx])
		}
	}
	
	// Ensure we don't have broken markdown links or formatting
	text = fixBrokenMarkdown(text)
	
	return text
}

// fixBrokenMarkdown attempts to fix common markdown breaks at chunk boundaries
func fixBrokenMarkdown(text string) string {
	// Remove incomplete markdown links [text](
	linkRe := regexp.MustCompile(`\[([^\]]*)\]\($`)
	text = linkRe.ReplaceAllString(text, "$1")
	
	// Remove incomplete markdown links [text
	incompleteLinkRe := regexp.MustCompile(`\[[^\]]*$`)
	text = incompleteLinkRe.ReplaceAllString(text, "")
	
	// Remove incomplete emphasis **text or *text
	emphasisRe := regexp.MustCompile(`\*{1,2}[^*]*$`)
	text = emphasisRe.ReplaceAllString(text, "")
	
	// Remove incomplete code blocks ```
	codeBlockRe := regexp.MustCompile("^```[^`]*$")
	text = codeBlockRe.ReplaceAllString(text, "")
	
	return strings.TrimSpace(text)
}