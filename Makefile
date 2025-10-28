# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=e-shop-server
API_PROTO_FILES=$(shell find api/protobuf -name *.proto)

# 工具二进制文件目录
TOOLS_BIN_DIR := $(CURDIR)/bin

# 声明工具变量
PROTOC_GEN_GO_PATH      = "$(TOOLS_BIN_DIR)/protoc-gen-go"
PROTOC_GEN_GO_GRPC_PATH = "$(TOOLS_BIN_DIR)/protoc-gen-go-grpc"
PROTOC_GEN_GATEWAY_PATH = "$(TOOLS_BIN_DIR)/protoc-gen-grpc-gateway"
MOCKGEN_PATH            = "$(TOOLS_BIN_DIR)/mockgen"
WIRE_PATH               = "$(TOOLS_BIN_DIR)/wire"
GOLANGCI_LINT_PATH      = "$(TOOLS_BIN_DIR)/golangci-lint"
GOIMPORTS_PATH          = "$(TOOLS_BIN_DIR)/goimports"


# 将工具目录添加到 PATH
export PATH := $(TOOLS_BIN_DIR):$(PATH)

# Default target
all: help

# Run the server
run: build
	@echo ">> running server..."
	./$(BINARY_NAME) 
	@echo "<< server stopped."

# Build the binary
build: api wire
	@echo ">> building binary..."
	$(GOBUILD) -o $(BINARY_NAME) ./cmd/user-srv 
	@echo "<< binary built."

# Test all packages
test:
	@echo ">> running tests..."
	$(GOTEST) -v ./...
	@echo "<< tests completed."

# 从 tools.go 安装工具（版本由 go.mod 管理）
.PHONY: tools
tools:
	@echo ">> installing tools from tools.go"
	@mkdir -p $(TOOLS_BIN_DIR)
	@echo "Installing protoc-gen-go..."
	@GOBIN=$(TOOLS_BIN_DIR) $(GOCMD) install google.golang.org/protobuf/cmd/protoc-gen-go
	@echo "Installing protoc-gen-go-grpc..."
	@GOBIN=$(TOOLS_BIN_DIR) $(GOCMD) install google.golang.org/grpc/cmd/protoc-gen-go-grpc
	@echo "Installing protoc-gen-grpc-gateway..."
	@GOBIN=$(TOOLS_BIN_DIR) $(GOCMD) install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
	@echo "Installing mockgen..."
	@GOBIN=$(TOOLS_BIN_DIR) $(GOCMD) install github.com/golang/mock/mockgen
	@echo "Installing wire..."
	@GOBIN=$(TOOLS_BIN_DIR) $(GOCMD) install github.com/google/wire/cmd/wire
	@echo "Installing golangci-lint..."
	@GOBIN=$(TOOLS_BIN_DIR) $(GOCMD) install github.com/golangci/golangci-lint/cmd/golangci-lint
	@echo "Installing goimports..."
	@GOBIN=$(TOOLS_BIN_DIR) $(GOCMD) install golang.org/x/tools/cmd/goimports@latest
	@echo "<< tools installed"

# Generate code from proto files
.PHONY: api
api: tools
	@echo ">> generating api code from proto"
	protoc -I. -I./third_party \
	       --go_out=. --go_opt=paths=source_relative \
	       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	       --grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative \
	       $(API_PROTO_FILES)
	@echo "<< api code generated."

# Generate dependency injection code with Wire
.PHONY: wire
wire: tools
	@echo ">> generating wire code"
	@$(WIRE_PATH) ./cmd/user-srv 
	@echo "<< wire code generated"

# Generate mocks
.PHONY: mock
mockgen: tools
	@echo ">> generating mocks"
	@$(MOCKGEN_PATH) -source=./internal/user-srv/biz/user.go -destination=./internal/user-srv/biz/mock/mocker_user.go -package=mock
	@echo ">> mocks generated"

# ----------------------
# 格式化代码
# ----------------------
.PHONY: format
format: tools
	@echo ">> formatting code with gofmt and goimports..."
# 格式化 go 文件并简化表达式，
# -s (simplify)，自动简化 Go 代码中可以优化的表达式
# -w (write)，直接把格式化后的代码写回原文件
	@find . -name "*.go" -not -path "./bin/*" -not -path "./third_party/*" | xargs gofmt -s -w
# 自动修复 import 排序和未使用的 import
	@find . -name "*.go" -not -path "./bin/*" -not -path "./third_party/*" | xargs $(GOIMPORTS_PATH) -w
	@echo "<< code formatted"

# Run code linter
.PHONY: lint
lint: tools 
	@echo ">> running code linter..."
	@$(GOLANGCI_LINT_PATH) run ./...
	@echo "<< linting completed."

# Run linter with auto-fix
.PHONY: lint-fix
lint-fix: tools format
	@echo ">> running linter with auto-fix..."
	@$(GOLANGCI_LINT_PATH) run --fix ./...
	@echo "<< auto-fix completed."

# Check tool versions
.PHONY: tools-version
tools-version:
	@echo ">> Tool versions:"
	@echo -n "protoc-gen-go: "
	@$(PROTOC_GEN_GO_PATH) --version 2>/dev/null || echo "not installed"
	@echo -n "wire: "
	@$(WIRE_PATH) version 2>/dev/null || echo "not installed"
	@echo -n "mockgen: "
	@$(MOCKGEN_PATH) --version 2>/dev/null || echo "not installed"
	@echo -n "golangci-lint: "
	@$(GOLANGCI_LINT_PATH) version 2>/dev/null || echo "not installed"

# Clean build artifacts
.PHONY: clean
clean:
	@echo ">> cleaning up..."
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -rf $(TOOLS_BIN_DIR)
	@echo "<< cleaned up."

# Help message
.PHONY: help
help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Primary Targets:"
	@echo "  run           Build and run the server"
	@echo "  build         Build the binary without running"
	@echo "  test          Run all tests"
	@echo ""
	@echo "Code Generation Targets:"
	@echo "  tools         Install development tools locally (versions from go.mod)"
	@echo "  api           Generate Go code from .proto files"
	@echo "  wire          Generate Go code from wire.go files"
	@echo "  mock          Generate mock code for interfaces"
	@echo ""
	@echo "Code Quality:"
	@echo "  lint          Run code linter to check for issues"
	@echo "  lint-fix      Run linter and automatically fix issues"
	@echo ""
	@echo "Utilities:"
	@echo "  tools-version Show installed tool versions"
	@echo "  clean         Clean up build artifacts and generated code"
	@echo ""
	@echo "Example workflow:"
	@echo "  make tools    # First time setup"
	@echo "  make lint     # Check code quality"
	@echo "  make test     # Run tests"
	@echo "  make run      # Build and run"

.PHONY: all run build test