---
title: '[TODO] Fix Semantic Backlinks SQL Constraint Issue'
tags:
  area: database
  component: semantic-backlinks
  priority: medium
  resolution: index-removed
  status: completed
  todo: true
created: 2025-08-24T01:11:00-07:00
modified: 2025-08-24T01:46:46.350205891-07:00
---

# ✅ RESOLVED: Semantic Backlinks SQL Constraint Issue

## Problem
When creating semantic backlinks, the system encountered this SQL error:
```
failed to upsert semantic backlink: Binder Error: Can not assign to column 'similarity_score' because it has a UNIQUE/PRIMARY KEY constraint or is referenced by an INDEX
```

## Root Cause ✅ IDENTIFIED
The issue was caused by the `similarity_score` column having an INDEX in the database schema:

```sql
CREATE INDEX IF NOT EXISTS idx_semantic_backlinks_score ON semantic_backlinks (similarity_score)
```

In DuckDB, when a column has an INDEX, it creates a constraint that prevents updates in `ON CONFLICT` clauses, treating the indexed column as if it has a constraint.

## Solution ✅ IMPLEMENTED
**Removed the unnecessary index** on the `similarity_score` column from the database schema in `internal/db/duckdb.go`:

```diff
- `CREATE INDEX IF NOT EXISTS idx_semantic_backlinks_score ON semantic_backlinks (similarity_score)`,
+ // Note: Removed idx_semantic_backlinks_score index because it prevents ON CONFLICT updates in DuckDB
```

## Verification ✅ TESTED
Using the improved justfile testing approach:

```bash
just test-backlinks
```

**Result**: Semantic backlinks are now creating successfully with no errors:
- ✅ `Created semantic backlink between 2 and 1 (similarity: 0.675)`
- ✅ `Created semantic backlink between 3 and 2 (similarity: 0.687)`
- ✅ `Created semantic backlink between 3 and 1 (similarity: 0.631)`
- ✅ Multiple successful semantic backlink creations throughout processing

## Technical Details

### Original UpsertSemanticBacklink Function
The function using `INSERT ... ON CONFLICT` works perfectly now:

```sql
INSERT INTO semantic_backlinks (id, memory_a_id, memory_b_id, similarity_score)
VALUES (nextval('seq_backlink_id'), ?, ?, ?)
ON CONFLICT (memory_a_id, memory_b_id) DO UPDATE SET
    similarity_score = EXCLUDED.similarity_score
```

### Performance Impact
- Removing the `similarity_score` index has minimal performance impact
- The primary access patterns use `memory_a_id` and `memory_b_id` indexes (retained)
- Semantic similarity queries are rare compared to relationship lookups

## Impact
- **High**: Core semantic search functionality now works perfectly
- Semantic backlinks provide enhanced memory relationships
- System processes all memories without SQL constraint errors
- RAG functionality operates at full capacity

## Files Modified
- `internal/db/duckdb.go` - Removed problematic similarity_score index
- `justfile` - Enhanced with comprehensive testing approach

## Testing Tools Created
- `just test-backlinks` - Comprehensive semantic backlinks testing
- `just test-json` - Universal JSON-RPC testing command
- Enhanced debug memory with justfile patterns

---
📝 **Title:** [TODO] Fix Semantic Backlinks SQL Constraint Issue
🏷️ **Tags:** todo, area: database, component: semantic-backlinks, priority: medium, status: completed, resolution: index-removed
📅 **Created:** 2025-08-24 01:11:00
🔄 **Modified:** 2025-08-24 01:12:51