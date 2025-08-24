---
title: '[TODO] Fix Semantic Backlinks SQL Constraint Issue'
tags:
  area: database
  component: semantic-backlinks
  priority: medium
  status: pending
  todo: true
created: 2025-08-24T01:11:00-07:00
modified: 2025-08-24T01:12:51.116592173-07:00
---

# Semantic Backlinks SQL Constraint Issue

## Problem
When creating semantic backlinks, the system encounters this SQL error:
```
failed to upsert semantic backlink: Binder Error: Can not assign to column 'similarity_score' because it has a UNIQUE/PRIMARY KEY constraint or is referenced by an INDEX
```

## Root Cause
The `UpsertSemanticBacklink` function in `internal/db/duckdb.go:471-476` attempts to update the `similarity_score` column, but there appears to be a constraint preventing this update.

## Current Behavior
- Semantic search works perfectly
- Memory processing completes successfully
- Only the backlink creation fails (non-critical)
- Error occurs during the ON CONFLICT UPDATE clause

## Investigation Needed
1. Check the `semantic_backlinks` table schema and constraints
2. Verify if `similarity_score` has an index that prevents updates
3. Determine if the UNIQUE constraint on `(memory_a_id, memory_b_id)` is interfering

## Impact
- **Low**: Core semantic search functionality works fine
- Semantic backlinks provide enhanced memory relationships but aren't critical
- System continues to operate normally despite this error

## Related Files
- `internal/db/duckdb.go` - UpsertSemanticBacklink function
- Database schema creation around line 115-127

## Next Steps
- [ ] Examine table schema and constraints
- [ ] Test different ON CONFLICT strategies
- [ ] Consider if similarity_score updates are actually needed
- [ ] Implement fix and test

---
📝 **Title:** [TODO] Fix Semantic Backlinks SQL Constraint Issue
🏷️ **Tags:** todo, database, semantic-backlinks, pending
📅 **Created:** 2025-08-24 01:11:00
🔄 **Modified:** 2025-08-24 01:11:00