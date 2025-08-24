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

# Create a test memory for testing
test-create-memory:
    echo '{"jsonrpc":"2.0","method":"tools/call","params":{"name":"create_memory","arguments":{"name":"test","content":"# Test Memory\n\nThis is a test memory."}},"id":1}' | go run cmd/simplemem/main.go

# Initialize memories directory
init:
    mkdir -p .memories
    @echo "Initialized .memories directory"