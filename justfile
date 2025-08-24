# SimpleMem MCP Server Justfile

# Default recipe - show available commands
default:
    @just --list

# Run the MCP server in stdio mode
run:
    go run cmd/simplemem/main.go

# Build the binary
build:
    go build -o simplemem cmd/simplemem/main.go

# Run tests
test:
    go test ./...

# Run tests with coverage
test-coverage:
    go test -cover ./...

# Clean build artifacts and memories directory
clean:
    rm -f simplemem
    rm -rf .memories

# Install dependencies
deps:
    go mod download
    go mod tidy

# Format code
fmt:
    go fmt ./...

# Run linter (requires golangci-lint)
lint:
    golangci-lint run

# Run the server with logging
run-verbose:
    go run cmd/simplemem/main.go 2>&1 | tee simplemem.log

# Check for dependency updates
check-deps:
    go list -u -m all

# Test with custom JSON-RPC call and optional custom DB
# Usage: just test-json '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"list_memories","arguments":{}},"id":1}' [/path/to/test.db]
test-json JSON DB="":
    #!/usr/bin/env bash
    if [ -n "{{DB}}" ]; then
        echo '{{JSON}}' | go run cmd/simplemem/main.go --db {{DB}}
    else
        echo '{{JSON}}' | go run cmd/simplemem/main.go
    fi

# Test semantic backlinks functionality (creates multiple memories to trigger semantic linking)
test-backlinks DB="/tmp/backlinks-test.db":
    @echo "Testing semantic backlinks with database: {{DB}}"
    @echo "Removing existing test database..."
    rm -f {{DB}}
    @echo "Creating first memory..."
    just test-json '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"create_memory","arguments":{"content":"---\nname: test-ai-memory\ntitle: AI Memory Test\n---\n\n# AI and Machine Learning\n\nThis memory is about artificial intelligence and machine learning concepts.\n"}},"id":1}' {{DB}}
    @echo "\nCreating second memory..."
    just test-json '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"create_memory","arguments":{"content":"---\nname: test-ml-memory\ntitle: ML Memory Test\n---\n\n# Machine Learning Algorithms\n\nThis memory covers various machine learning algorithms and their applications.\n"}},"id":2}' {{DB}}
    @echo "\nListing memories to verify creation..."
    just test-json '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"list_memories","arguments":{}},"id":3}' {{DB}}

# Quick test with clean database
test-clean DB="/tmp/test-clean.db":
    rm -f {{DB}}
    just test-json '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"list_memories","arguments":{}},"id":1}' {{DB}}

# Initialize memories directory
init:
    mkdir -p .memories
    @echo "Initialized .memories directory"