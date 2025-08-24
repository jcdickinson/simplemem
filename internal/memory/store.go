package memory

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Store struct {
	basePath string
	mu       sync.RWMutex
}

func NewStore(basePath string) *Store {
	return &Store{
		basePath: basePath,
	}
}

func (s *Store) Initialize() error {
	return os.MkdirAll(s.basePath, 0755)
}

func (s *Store) Create(name, content string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !strings.HasSuffix(name, ".md") {
		name = name + ".md"
	}

	path := filepath.Join(s.basePath, name)
	
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("memory %s already exists", name)
	}

	// Parse frontmatter if present, or create new one
	fm, body, err := ParseDocument(content)
	if err != nil {
		// If parsing fails, treat entire content as body with new frontmatter
		fm = &Frontmatter{}
		body = content
	}

	// Update timestamps for new document
	fm.UpdateTimestamps(true)

	// Format the complete document
	finalContent, err := FormatDocument(fm, body)
	if err != nil {
		return fmt.Errorf("failed to format document: %w", err)
	}

	return os.WriteFile(path, []byte(finalContent), 0644)
}

func (s *Store) Read(name string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !strings.HasSuffix(name, ".md") {
		name = name + ".md"
	}

	path := filepath.Join(s.basePath, name)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("memory %s not found", name)
		}
		return "", err
	}

	return string(data), nil
}

func (s *Store) Update(name, content string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !strings.HasSuffix(name, ".md") {
		name = name + ".md"
	}

	path := filepath.Join(s.basePath, name)
	
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("memory %s not found", name)
	}

	// Read existing file to preserve created timestamp
	existingData, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read existing file: %w", err)
	}

	existingFm, _, _ := ParseDocument(string(existingData))

	// Parse new content
	fm, body, err := ParseDocument(content)
	if err != nil {
		fm = &Frontmatter{}
		body = content
	}

	// Preserve created timestamp if it exists
	if !existingFm.Created.IsZero() {
		fm.Created = existingFm.Created
	}

	// Update modified timestamp
	fm.UpdateTimestamps(false)

	// Format the complete document
	finalContent, err := FormatDocument(fm, body)
	if err != nil {
		return fmt.Errorf("failed to format document: %w", err)
	}

	return os.WriteFile(path, []byte(finalContent), 0644)
}

func (s *Store) Delete(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !strings.HasSuffix(name, ".md") {
		name = name + ".md"
	}

	path := filepath.Join(s.basePath, name)
	
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("memory %s not found", name)
	}

	return os.Remove(path)
}

func (s *Store) List() ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entries, err := os.ReadDir(s.basePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	var memories []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
			name := strings.TrimSuffix(entry.Name(), ".md")
			memories = append(memories, name)
		}
	}

	return memories, nil
}

func (s *Store) Search(query string) (map[string][]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entries, err := os.ReadDir(s.basePath)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string][]string{}, nil
		}
		return nil, err
	}

	results := make(map[string][]string)
	queryLower := strings.ToLower(query)

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
			path := filepath.Join(s.basePath, entry.Name())
			content, err := os.ReadFile(path)
			if err != nil {
				continue
			}

			// Parse document to search in body only (not frontmatter)
			_, body, _ := ParseDocument(string(content))
			
			lines := strings.Split(body, "\n")
			var matches []string
			
			for i, line := range lines {
				if strings.Contains(strings.ToLower(line), queryLower) {
					context := fmt.Sprintf("Line %d: %s", i+1, line)
					matches = append(matches, context)
				}
			}

			if len(matches) > 0 {
				name := strings.TrimSuffix(entry.Name(), ".md")
				results[name] = matches
			}
		}
	}

	return results, nil
}

// MemoryInfo contains both content and metadata
type MemoryInfo struct {
	Name        string
	Content     string
	Body        string
	Frontmatter *Frontmatter
}

// ReadWithMetadata reads a memory and returns both content and parsed metadata
func (s *Store) ReadWithMetadata(name string) (*MemoryInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if !strings.HasSuffix(name, ".md") {
		name = name + ".md"
	}

	path := filepath.Join(s.basePath, name)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("memory %s not found", name)
		}
		return nil, err
	}

	content := string(data)
	fm, body, err := ParseDocument(content)
	if err != nil {
		// If parsing fails, treat entire content as body
		fm = &Frontmatter{}
		body = content
	}

	return &MemoryInfo{
		Name:        strings.TrimSuffix(name, ".md"),
		Content:     content,
		Body:        body,
		Frontmatter: fm,
	}, nil
}

// SearchByTag finds all memories with a specific tag
func (s *Store) SearchByTag(tag string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entries, err := os.ReadDir(s.basePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	var results []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
			path := filepath.Join(s.basePath, entry.Name())
			content, err := os.ReadFile(path)
			if err != nil {
				continue
			}

			fm, _, err := ParseDocument(string(content))
			if err != nil {
				continue
			}

			if fm.HasTag(tag) {
				name := strings.TrimSuffix(entry.Name(), ".md")
				results = append(results, name)
			}
		}
	}

	return results, nil
}

// GetAllTags returns all unique tags across all memories
func (s *Store) GetAllTags() (map[string][]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entries, err := os.ReadDir(s.basePath)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string][]string{}, nil
		}
		return nil, err
	}

	tagMap := make(map[string][]string)
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
			path := filepath.Join(s.basePath, entry.Name())
			content, err := os.ReadFile(path)
			if err != nil {
				continue
			}

			fm, _, err := ParseDocument(string(content))
			if err != nil {
				continue
			}

			name := strings.TrimSuffix(entry.Name(), ".md")
			for tag := range fm.Tags {
				tagMap[tag] = append(tagMap[tag], name)
			}
		}
	}

	return tagMap, nil
}