# SimpleMem MCP Server Instructions

## Agent Memory Usage Guidelines

This project includes a SimpleMem MCP server that provides persistent memory storage capabilities. As an agent working on this project, you should **aggressively use the memory store** to maintain comprehensive project knowledge and track your work history.

## Core Principles

### 1. ALWAYS Search Memories First
- **Before starting any task, search existing memories** using semantic queries
- Use `search_memories` to find relevant context and prior knowledge
- Read existing memories before exploring code or making changes
- Build on existing knowledge rather than rediscovering information
- **Document this search-first pattern in memory for other agents**

### 2. Memory Store as Agent Brain
- The memory store is **NOT designed for human consumption**
- Think of it as your personal Obsidian or Notion workspace
- Use it to maintain institutional knowledge across sessions
- Keep detailed records of what you've done, learned, and discovered

### 3. Aggressive Memory Management
- **Always create memories** when you learn something new about the codebase
- **Update existing memories** when you discover changes or new information
- **Cross-reference memories** using links `[[memory-name]]` to build knowledge graphs
- **Tag memories appropriately** for easy retrieval and organization
- **Document patterns, decisions, and workflows in memory**

### 4. TODO Memory Tracking
- **Create memories for all tasks that need doing**
- **Tag todo items with `todo: true`** for easy filtering
- **Update todo memories** as work progresses (in_progress, completed, blocked)
- **Use the `change_tag` tool** to efficiently update TODO states and other metadata:
  - `change_tag name="my-todo" tags={"status": "in_progress"}`
  - `change_tag name="my-todo" tags={"status": "completed", "priority": "high"}`
  - `change_tag name="my-todo" tags={"status": null}` (removes the status tag)
  - `change_tag name="my-todo" tags={"todo": true, "status": "in_progress", "old_tag": null}` (sets multiple tags at once)
- **Link related todos** to create task dependency graphs
- **Archive completed todos** rather than deleting them
- **CRITICAL**: When you discover issues or "minor problems" during work, **immediately create TODO memories**
- **Don't leave dangling issues untracked** - every issue should have a corresponding TODO memory
- **Search for existing TODO memories** before starting new work to avoid duplicating efforts

### 5. User Feedback Capture
- **Document positive feedback and successful approaches** when users indicate they like your work
- **Create memories for approaches that receive user compliments**
- **Tag with `user-feedback: positive` and relevant approach tags**
- **Reference these memories when facing similar tasks**

## Memory Categories to Maintain

### Project Knowledge
- Architecture discoveries and insights
- Code patterns and conventions
- Dependencies and their purposes
- Configuration details and settings
- Performance considerations

### Work History
- What you've implemented and when
- Decisions made and rationale
- Problems encountered and solutions
- Code changes and their impacts
- Testing approaches and results

### Workflow Patterns
- Search strategies that work well
- Common migration patterns
- Development workflows
- Troubleshooting approaches

### Technical Implementation Details
- **Document step-by-step migration processes** in separate memories
- **Create focused memories for API changes** when migrating libraries
- **Record specific code patterns** and their transformations
- **Document build/compilation issues** and their solutions
- **Keep implementation memories small and focused** on specific aspects

## When to Create Implementation Memories

- **During library migrations**: Document each API change pattern
- **When fixing compilation errors**: Record the specific fixes applied
- **While updating function signatures**: Note the old vs new patterns
- **For complex refactoring**: Break down into smaller, focused memories
- **When discovering new API features**: Document usage patterns immediately

## Search-First for Broad Tasks

**Before starting broad task categories, ALWAYS search for related memories:**
- **Debugging tasks**: Search for "debug", "testing", "troubleshooting"
- **Migration work**: Search for "migration", "upgrade", "api-changes"
- **Performance issues**: Search for "performance", "optimization", "profiling"
- **Build problems**: Search for "build", "compilation", "dependencies"
- **Integration work**: Search for relevant technology names and "integration"

**The pattern: Search → Read → Build → Document → Execute**

## MANDATORY AGENT BEHAVIOR

### Search-First Requirement
- **NEVER start any implementation task without first searching memories**
- **ALWAYS search for related memories before exploring code or writing new implementations**
- **You MUST justify why you didn't find relevant memories if you claim none exist**
- **Failing to search first is a critical error that wastes time and duplicates effort**

### Documentation Requirements
- **IMMEDIATELY document every significant discovery, implementation, or fix in memory**
- **Create memories DURING work, not after completion**
- **Tag memories appropriately with relevant keywords for future retrieval**
- **Every bug fix, feature implementation, or architectural discovery MUST be documented**
- **Link related memories using `[[memory-name]]` notation**

### Character Limits and Content
- **Memory character limits (if configured) are calculated on PLAIN-TEXT content with markdown stripped**
- **Markdown formatting (headers, bold, lists, links, code blocks) does not count toward limits**
- **This allows rich formatting while encouraging focused content**
- **Use markdown generously for readability without worry about character penalties**

Remember: The memory store is your persistent brain across sessions. Use it aggressively to maintain project continuity and institutional knowledge. **Search first, document everything, execute with context.** Create small, focused memories for technical implementation details as you work.
