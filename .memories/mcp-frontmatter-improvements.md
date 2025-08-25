---
title: MCP Frontmatter API Improvements
description: Two key improvements needed for the SimpleMem MCP interface
tags:
  api: improvements
  priority: high
  todo: true
created: 2025-08-24T01:16:20.58842268-07:00
modified: 2025-08-24T18:48:17.48628234-07:00
---

# MCP Frontmatter API Improvements

## Current Issues

### 1. Redundant Timestamp Fields
Currently accepts `created` and `modified` timestamps in request frontmatter. These should be automatically set server-side using `time.Now()`.

### 2. JSON Envelope Requirement
Memories must be submitted with separate `name` parameter in JSON request. A cleaner approach would allow the `name` field in frontmatter itself.

## Implementation Plan

### Phase 1: Automatic Timestamp Handling
- ✅ Current code already handles this via `fm.UpdateTimestamps(true)`
- ✅ `UpdateTimestamps` method automatically sets timestamps using `time.Now()`
- ✅ Response includes timestamps in read operations

**Status: Already implemented correctly**

### Phase 2: Name Field Support
- Add logic to extract `name` from frontmatter if present
- Fall back to JSON parameter if frontmatter name not provided
- Update MCP tool description to document capability

## Current Status

Timestamp handling works correctly. Users can pass timestamps in requests, but they get overwritten by `UpdateTimestamps` call, which is the correct behavior.

For name field support, we need to modify the MCP handler to check frontmatter for name field first.