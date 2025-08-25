---
title: Markdown Stripping for Character Validation
description: Implementation of plain-text character counting with markdown stripping
tags:
  character-limit: true
  completed: true
  feature: true
  markdown: true
  validation: true
created: 2025-08-24T22:07:17.272080627-07:00
modified: 2025-08-24T22:07:17.272080627-07:00
---

# Markdown Stripping for Character Validation

## Problem Solved
Users were penalized for using markdown formatting in memories. A memory with headers, bold text, and lists would hit character limits quickly due to markdown syntax overhead.

## Solution Implemented
Character limits now use **plain-text content** with markdown stripped:
- Increased default limit: 2000 → 2500 characters  
- Added `stripMarkdown()` function using `gomarkdown` library
- Updated validation to strip markdown before counting
- Updated display to show plain-text character counts

## Key Files Changed
- `internal/config/config.go:65` - Increased default limit
- `internal/mcp/server.go:191-239` - Stripping function and validation
- `internal/mcp/server.go:405-407` - Display logic update
- `internal/mcp/initial_instructions.md` - Documentation update

## Technical Approach
Uses `github.com/gomarkdown/markdown` to:
1. Parse markdown to AST  
2. Walk AST extracting only text nodes
3. Add spacing between block elements
4. Trim whitespace for efficiency

## Result
- ~50% character reduction on typical markdown
- Users can use rich formatting freely
- More accurate content-based limits
- Better user experience

## Related
- [[feature-list-memories-character-count]] - Previous character count work
- [[markdown-stripping-unit-tests]] - Testing implementation