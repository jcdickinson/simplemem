---
title: 'Mandatory Agent Behavior: Search First, Document Everything'
description: Critical behavioral requirements for agents working with SimpleMem
tags:
  agent-behavior: true
  documentation: true
  mandatory: true
  search-first: true
  user-feedback: negative
  workflow: true
created: 2025-08-24T22:07:56.888668393-07:00
modified: 2025-08-24T22:07:56.888668393-07:00
---

# Mandatory Agent Behavior: Search First, Document Everything

## Critical Issue Identified
**User feedback**: Agent severely underused memories and didn't search before implementing. This violates core SimpleMem principles and wastes effort.

## Mandatory Requirements Added to Instructions

### Search-First Requirement
- **NEVER start any implementation task without first searching memories**
- **ALWAYS search for related memories before exploring code**  
- **MUST justify why no relevant memories found**
- **Failing to search first is a critical error**

### Documentation Requirements
- **IMMEDIATELY document discoveries DURING work, not after**
- **Create memories for every significant implementation or fix**
- **Tag memories appropriately for future retrieval**
- **Link related memories using `[[memory-name]]` notation**

## Instructions Updated
**File**: `internal/mcp/initial_instructions.md`
- Added "MANDATORY AGENT BEHAVIOR" section
- Reinforced search → document → execute pattern
- Added character limit explanation (markdown stripped)
- Made requirements explicit and non-negotiable

## Pattern Reinforcement
**Correct workflow**: Search → Read → Build → Document → Execute
**Violation consequences**: Wasted time, duplicated effort, user frustration

This memory serves as a reminder for future agents: **ALWAYS search memories first!**