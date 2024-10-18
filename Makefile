# Makefile for Godis

# Variables
BINARY_NAME=build/godis-server
CLI_NAME=build/godis-cli
GO_CMD=go
SRC_DIR=./cmd/server
CLI_DIR=./cmd/client

# Default target: build both the server and CLI
all: build-dir build-server build-cli

# Create the build directory if it doesn't exist
build-dir:
	@mkdir -p build

# Build the server and place the binary in the build folder
build-server: build-dir
	@echo "Building Godis Server..."
	$(GO_CMD) build -o $(BINARY_NAME) $(SRC_DIR)

# Build the CLI and place the binary in the build folder
build-cli: build-dir
	@echo "Building Godis CLI..."
	$(GO_CMD) build -o $(CLI_NAME) $(CLI_DIR)

# Run the server
run: build-server
	@echo "Running Godis Server..."
	./$(BINARY_NAME)

# Run the CLI
run-cli: build-cli
	@echo "Running Godis CLI..."
	./$(CLI_NAME)

# Clean up binaries
clean:
	@echo "Cleaning up..."
	rm -rf build

# Test the application
test:
	@echo "Running tests..."
	$(GO_CMD) test ./...

# AOF persistence setup (example task)
persistence:
	@echo "Setting up AOF Persistence..."
	# Add commands for setting up AOF persistence if needed

# Help message
help:
	@echo "Usage:"
	@echo "  make           - Build the project (both server and CLI)"
	@echo "  make run-server - Build and run the server"
	@echo "  make run-cli   - Build and run the CLI"
	@echo "  make clean     - Remove binaries"
	@echo "  make test      - Run all tests"
	@echo "  make persistence - Setup AOF persistence"