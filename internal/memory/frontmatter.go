package memory

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Frontmatter represents the YAML frontmatter of a memory document
type Frontmatter struct {
	Name        string                 `yaml:"name,omitempty"`
	Title       string                 `yaml:"title,omitempty"`
	Description string                 `yaml:"description,omitempty"`
	Tags        map[string]interface{} `yaml:"tags,omitempty"`
	Created     time.Time              `yaml:"created,omitempty"`
	Modified    time.Time              `yaml:"modified,omitempty"`
	Links       []string               `yaml:"links,omitempty"`
	Metadata    map[string]interface{} `yaml:",inline"`
}

// ParseDocument separates frontmatter from content, handling multiple consecutive frontmatter blocks
func ParseDocument(content string) (*Frontmatter, string, error) {
	// Check if content starts with frontmatter delimiter
	if !strings.HasPrefix(content, "---\n") && !strings.HasPrefix(content, "---\r\n") {
		// No frontmatter, return empty frontmatter and full content
		return &Frontmatter{}, content, nil
	}

	lines := strings.Split(content, "\n")
	var mergedFrontmatter Frontmatter
	currentLine := 0
	
	// Process multiple consecutive frontmatter blocks
	for currentLine < len(lines) {
		// Check if we're at a frontmatter start
		if strings.TrimSpace(lines[currentLine]) != "---" {
			break
		}
		
		// Find the end of this frontmatter block
		endIndex := -1
		for i := currentLine + 1; i < len(lines); i++ {
			if strings.TrimSpace(lines[i]) == "---" {
				endIndex = i
				break
			}
		}
		
		if endIndex == -1 {
			// No closing delimiter found, treat remaining content as body
			bodyLines := lines[currentLine:]
			body := strings.Join(bodyLines, "\n")
			return &mergedFrontmatter, body, nil
		}
		
		// Extract this frontmatter block
		frontmatterLines := lines[currentLine+1:endIndex]
		frontmatterStr := strings.Join(frontmatterLines, "\n")
		
		// Skip empty frontmatter blocks
		if strings.TrimSpace(frontmatterStr) == "" {
			currentLine = endIndex + 1
			continue
		}
		
		// Parse this frontmatter block
		var fm Frontmatter
		if err := yaml.Unmarshal([]byte(frontmatterStr), &fm); err != nil {
			// If parsing fails, log warning and skip this block
			currentLine = endIndex + 1
			continue
		}
		
		// Merge frontmatter (later blocks override earlier ones)
		if fm.Name != "" {
			mergedFrontmatter.Name = fm.Name
		}
		if fm.Title != "" {
			mergedFrontmatter.Title = fm.Title
		}
		if fm.Description != "" {
			mergedFrontmatter.Description = fm.Description
		}
		if !fm.Created.IsZero() {
			mergedFrontmatter.Created = fm.Created
		}
		if !fm.Modified.IsZero() {
			mergedFrontmatter.Modified = fm.Modified
		}
		if len(fm.Tags) > 0 {
			if mergedFrontmatter.Tags == nil {
				mergedFrontmatter.Tags = make(map[string]interface{})
			}
			for k, v := range fm.Tags {
				mergedFrontmatter.Tags[k] = v
			}
		}
		if len(fm.Links) > 0 {
			mergedFrontmatter.Links = append(mergedFrontmatter.Links, fm.Links...)
		}
		if len(fm.Metadata) > 0 {
			if mergedFrontmatter.Metadata == nil {
				mergedFrontmatter.Metadata = make(map[string]interface{})
			}
			for k, v := range fm.Metadata {
				mergedFrontmatter.Metadata[k] = v
			}
		}
		
		currentLine = endIndex + 1
	}
	
	// Extract body (everything after the last frontmatter block)
	bodyLines := lines[currentLine:]
	body := strings.Join(bodyLines, "\n")
	body = strings.TrimLeft(body, "\n\r")

	return &mergedFrontmatter, body, nil
}

// FormatDocument combines frontmatter and content into a complete document
func FormatDocument(fm *Frontmatter, content string) (string, error) {
	if fm == nil || (fm.Title == "" && fm.Description == "" && len(fm.Tags) == 0 && 
		fm.Created.IsZero() && fm.Modified.IsZero() && len(fm.Links) == 0 && len(fm.Metadata) == 0) {
		// No frontmatter to add
		return content, nil
	}

	// Create a copy for formatting (exclude name field from document output)
	formatFm := *fm
	formatFm.Name = "" // Don't include name in the document itself

	// Marshal frontmatter to YAML
	var buf bytes.Buffer
	encoder := yaml.NewEncoder(&buf)
	encoder.SetIndent(2)
	if err := encoder.Encode(&formatFm); err != nil {
		return "", fmt.Errorf("failed to encode frontmatter: %w", err)
	}

	// Combine with content
	result := "---\n" + buf.String() + "---\n\n" + content
	return result, nil
}

// ExtractTags returns all tags from the frontmatter
func (fm *Frontmatter) ExtractTags() map[string]interface{} {
	if fm.Tags == nil {
		return make(map[string]interface{})
	}
	return fm.Tags
}

// HasTag checks if a specific tag exists
func (fm *Frontmatter) HasTag(tag string) bool {
	if fm.Tags == nil {
		return false
	}
	_, exists := fm.Tags[tag]
	return exists
}

// GetTagValue returns the value associated with a tag
func (fm *Frontmatter) GetTagValue(tag string) (interface{}, bool) {
	if fm.Tags == nil {
		return nil, false
	}
	value, exists := fm.Tags[tag]
	return value, exists
}

// AddTag adds a tag with an optional value
func (fm *Frontmatter) AddTag(tag string, value interface{}) {
	if fm.Tags == nil {
		fm.Tags = make(map[string]interface{})
	}
	if value == nil {
		fm.Tags[tag] = true
	} else {
		fm.Tags[tag] = value
	}
}

// RemoveTag removes a tag
func (fm *Frontmatter) RemoveTag(tag string) {
	if fm.Tags != nil {
		delete(fm.Tags, tag)
	}
}

// UpdateTimestamps updates created and modified timestamps
func (fm *Frontmatter) UpdateTimestamps(isNew bool) {
	now := time.Now()
	if isNew && fm.Created.IsZero() {
		fm.Created = now
	}
	fm.Modified = now
}