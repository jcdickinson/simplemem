---
title: MCP Frontmatter API Improvements
description: Two key improvements needed for the SimpleMem MCP interface - automatic timestamp handling and name field support
tags:
  api: improvements
  priority: high
  todo: true
created: 2025-08-24T01:16:20.58842268-07:00
modified: 2025-08-24T01:16:20.58842268-07:00
---

# MCP Frontmatter API Improvements

## Current Issues

### 1. Redundant Timestamp Fields in Request
Currently, the MCP create_memory tool accepts `created` and `modified` timestamps in the request frontmatter. This doesn't make sense since:
- These should be automatically set server-side using `time.Now()`
- The client shouldn't control when something was "created" or "modified" 
- The timestamps should be added to the response for reference

### 2. JSON Envelope Requirement
Currently, memories must be submitted through the MCP with a separate `name` parameter in the JSON request. A cleaner approach would be:
- Allow the `name` field to be specified in the frontmatter itself
- Remove the requirement for the separate JSON parameter
- Parse the name from frontmatter when present

## Implementation Plan

### Phase 1: Automatic Timestamp Handling
- ✅ Current code already handles this in `store.go:49` via `fm.UpdateTimestamps(true)`
- ✅ The `UpdateTimestamps` method automatically sets timestamps using `time.Now()`
- ✅ Response already includes the timestamps in read operations

**Status: Already implemented correctly**

### Phase 2: Name Field Support
- Add logic to extract `name` from frontmatter if present
- Fall back to JSON parameter if frontmatter name is not provided
- Update MCP tool description to document this capability

## Current Implementation Status

The timestamp handling is already working correctly. The main issue is that users can pass timestamps in requests, but they get overwritten anyway by the `UpdateTimestamps` call. This is actually good behavior - the user input is ignored and server-generated timestamps are used.

For the name field support, we need to modify the MCP handler to check frontmatter for a name field first.