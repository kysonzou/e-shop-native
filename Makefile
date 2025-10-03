# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=e-shop-server
API_PROTO_FILES=$(shell find api/protobuf -name *.proto)

# Default target
all: help

# Run the server
run_user_srv:
	@echo "Building and running the server..."
	$(GOBUILD) -o $(BINARY_NAME) ./cmd/user-srv
	./$(BINARY_NAME) 
	@echo "Server stopped."

# Build the binary
build:
	@echo "Building the binary..."
	$(GOBUILD) -o $(BINARY_NAME) ./cmd/server
	@echo "Done."

# Generate code from proto files
api:
	@echo "Generating Go code from .proto files..."
	@protoc -I. -I./third_party \
	       --go_out=. --go_opt=paths=source_relative \
	       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
		   --grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative \
	       $(API_PROTO_FILES)
	@echo "Done."

# Generate dependency injection code with Wire
wire:
	@echo "--- Generating wire code ---"
# go generate 会自动在 cmd/server 目录下找到 //go:generate wire 并执行
#	$(GOCMD) generate ./..
# 或者可以用wire实现
	wire ./...
	@echo "Done."

# Clean build artifacts
clean:
	@echo "Cleaning up..."
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	@echo "Done."

# Test all packages
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...
	@echo "Tests completed."

# Help message
help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Targets:"
	@echo "  run        Build and run the server"
	@echo "  build      Build the binary"
	@echo "  api        Generate Go code from .proto files"
	@echo "  wire       Generate Go code from wire.go files"
	@echo "  clean      Clean up build artifacts"
	@echo "  test       Run all tests"

.PHONY: all run build api clean test help