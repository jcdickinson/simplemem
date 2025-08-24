---
title: SimpleMem Frontmatter System
description: Analysis of the YAML frontmatter parsing and metadata management system
tags:
  architecture: true
  golang: true
  metadata: true
  parsing: true
  yaml: true
created: 2025-08-24T00:21:19.107685432-07:00
modified: 2025-08-24T00:21:19.107685432-07:00
---

# SimpleMem Frontmatter System

The frontmatter system (`internal/memory/frontmatter.go`) provides sophisticated YAML metadata parsing for memory documents, supporting multiple frontmatter blocks and flexible metadata structures.

## Frontmatter Structure

### Core Frontmatter Type

```go
type Frontmatter struct {
    Title       string                 `yaml:"title,omitempty"`
    Description string                 `yaml:"description,omitempty"`
    Tags        map[string]interface{} `yaml:"tags,omitempty"`
    Created     time.Time              `yaml:"created,omitempty"`
    Modified    time.Time              `yaml:"modified,omitempty"`
    Links       []string               `yaml:"links,omitempty"`
    Metadata    map[string]interface{} `yaml:",inline"`
}
```

### Field Descriptions

#### Standard Fields
- **Title**: Human-readable memory title (optional)
- **Description**: Brief description of memory content (optional)
- **Tags**: Key-value metadata for categorization and filtering
- **Created**: Timestamp when memory was created
- **Modified**: Timestamp when memory was last updated
- **Links**: Explicit links to other memories

#### Flexible Metadata
- **Metadata**: Inline YAML fields for custom attributes
- **Interface{} Values**: Support for any YAML-compatible data types
- **Extensible**: New fields automatically captured in Metadata map

## Parsing Algorithm

### Multi-Block Support

The parser supports multiple consecutive frontmatter blocks in a single document:

```markdown
---
title: "My Memory"
tags:
  category: notes
---
---
description: "Updated description"
tags:
  priority: high
---

# Memory Content

Rest of the markdown content...
```

### Parsing Process

```go
func ParseDocument(content string) (*Frontmatter, string, error)
```

1. **Delimiter Detection**: Check for starting `---` delimiter
2. **Block Iteration**: Process each consecutive frontmatter block
3. **End Detection**: Find closing `---` delimiter for each block
4. **YAML Parsing**: Parse each block as valid YAML
5. **Merging Strategy**: Later blocks override earlier fields
6. **Body Extraction**: Return remaining content as document body

### Merging Strategy

#### Field Override Rules
- **Simple Fields**: Later values completely replace earlier ones
- **Tags Map**: Merge all tag key-value pairs (later values override)
- **Links Array**: Append all links together
- **Metadata Map**: Merge all inline fields (later values override)
- **Timestamps**: Later timestamps override earlier ones

#### Error Handling
- **Invalid YAML**: Skip malformed blocks, continue parsing
- **Missing Delimiters**: Treat as body content if no closing `---`
- **Empty Blocks**: Skip empty frontmatter blocks
- **Graceful Degradation**: Return partial results with warnings

## Integration Points

### Memory Store Integration

#### Document Creation
```go
// Automatic frontmatter generation on create
fm := &Frontmatter{
    Created:  time.Now(),
    Modified: time.Now(),
}
```

#### Document Updates
```go  
// Preserve creation time, update modification time
if existing.Created.IsZero() {
    fm.Created = time.Now()
} else {
    fm.Created = existing.Created
}
fm.Modified = time.Now()
```

### Database Synchronization

#### Metadata Extraction
- Parse frontmatter during database sync
- Store structured data in normalized tables
- Index tags and timestamps for efficient queries
- Handle type conversion for different YAML value types

#### Tag Processing
```go
// Convert interface{} values to appropriate database types
for key, value := range frontmatter.Tags {
    switch v := value.(type) {
    case bool:
        // Store boolean flags
    case string:
        // Store string values
    case int, int64, float64:
        // Store numeric values
    default:
        // Convert to JSON string for complex types
    }
}
```

## Advanced Features

### Document Serialization

#### Generate Frontmatter
```go
func (fm *Frontmatter) ToYAML() ([]byte, error)
```
- Convert frontmatter struct back to YAML
- Used for document updates and creation
- Maintains consistent formatting

#### Document Assembly
```go
func AssembleDocument(fm *Frontmatter, body string) string
```
- Combine frontmatter and body into complete document
- Add proper YAML delimiters
- Handle empty frontmatter gracefully

### Link Extraction

The system can extract and track document relationships:

#### Explicit Links
- **Links Field**: Direct references in frontmatter
- **Wiki-style Links**: `[[memory-name]]` patterns in content
- **Markdown Links**: Standard `[text](memory-name)` patterns

#### Link Processing
- Extract during document parsing
- Store in database for backlink analysis
- Update when documents change
- Support both absolute and relative references

## Validation and Error Handling

### YAML Validation
- **Strict Parsing**: Reject invalid YAML syntax
- **Type Safety**: Handle mixed-type tag values safely
- **Encoding Handling**: Support UTF-8 content correctly

### Content Validation
- **Frontmatter Position**: Must be at document start
- **Delimiter Format**: Require exact `---` format
- **Line Ending Handling**: Support both Unix and Windows line endings

### Error Recovery
- **Partial Success**: Return valid data even if some blocks fail
- **Informative Errors**: Provide context for parsing failures
- **Graceful Defaults**: Use empty frontmatter if parsing completely fails

The frontmatter system provides a robust foundation for memory metadata management while maintaining flexibility for diverse use cases and ensuring reliable parsing even with complex or partially invalid input.