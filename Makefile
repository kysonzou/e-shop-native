# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=e-shop-server
API_PROTO_FILES=$(shell find api/protobuf -name *.proto)

# ====================================================================================
# 工具链管理 (Toolchain Management)
# ====================================================================================

# 工具二进制文件目录
TOOLS_BIN_DIR := $(CURDIR)/bin
# 将工具目录添加到 PATH 的最前端，确保优先使用本地安装的工具
export PATH := $(TOOLS_BIN_DIR):$(PATH)

# 声明每个工具的最终二进制文件路径
PROTOC_GEN_GO_PATH      := $(TOOLS_BIN_DIR)/protoc-gen-go
PROTOC_GEN_GO_GRPC_PATH := $(TOOLS_BIN_DIR)/protoc-gen-go-grpc
PROTOC_GEN_GATEWAY_PATH := $(TOOLS_BIN_DIR)/protoc-gen-grpc-gateway
MOCKGEN_PATH            := $(TOOLS_BIN_DIR)/mockgen
WIRE_PATH               := $(TOOLS_BIN_DIR)/wire
GOLANGCI_LINT_PATH      := $(TOOLS_BIN_DIR)/golangci-lint
GOIMPORTS_PATH          := $(TOOLS_BIN_DIR)/goimports
BUF_PATH                := $(TOOLS_BIN_DIR)/buf



# tools 目标是一个伪目标，它的前提条件是所有工具的二进制文件
# 当运行 `make tools` 或任何依赖它的目标时，make 会检查每个工具文件是否存在
# 如果某个文件不存在，make 就会查找并执行对应的安装规则
.PHONY: tools
tools: $(PROTOC_GEN_GO_PATH) $(PROTOC_GEN_GO_GRPC_PATH) $(PROTOC_GEN_GATEWAY_PATH) $(MOCKGEN_PATH) $(WIRE_PATH) $(GOLANGCI_LINT_PATH) $(GOIMPORTS_PATH)\
	   $(BUF_PATH)
	@echo ">> All tools are installed and up to date."

# 每个工具的安装规则
# 目标是二进制文件路径，前提条件是 go.mod 和 tools.go
# 这意味着如果 go.mod/tools.go 更新了，或者二进制文件不存在，此规则就会被触发
# GOBIN=$(TOOLS_BIN_DIR) 指定了 `go install` 的安装目录
$(PROTOC_GEN_GO_PATH): go.mod tools.go
	@echo ">> Installing protoc-gen-go (version from go.mod)..."
	@GOBIN=$(TOOLS_BIN_DIR) $(GOCMD) install google.golang.org/protobuf/cmd/protoc-gen-go

$(PROTOC_GEN_GO_GRPC_PATH): go.mod tools.go
	@echo ">> Installing protoc-gen-go-grpc..."
	@GOBIN=$(TOOLS_BIN_DIR) $(GOCMD) install google.golang.org/grpc/cmd/protoc-gen-go-grpc

$(PROTOC_GEN_GATEWAY_PATH): go.mod tools.go
	@echo ">> Installing protoc-gen-grpc-gateway..."
	@GOBIN=$(TOOLS_BIN_DIR) $(GOCMD) install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway

$(MOCKGEN_PATH): go.mod tools.go
	@echo ">> Installing mockgen..."
	@GOBIN=$(TOOLS_BIN_DIR) $(GOCMD) install github.com/golang/mock/mockgen

$(WIRE_PATH): go.mod tools.go
	@echo ">> Installing wire..."
	@GOBIN=$(TOOLS_BIN_DIR) $(GOCMD) install github.com/google/wire/cmd/wire

$(GOLANGCI_LINT_PATH): go.mod tools.go
	@echo ">> Installing golangci-lint..."
	@GOBIN=$(TOOLS_BIN_DIR) $(GOCMD) install github.com/golangci/golangci-lint/cmd/golangci-lint

$(GOIMPORTS_PATH): go.mod tools.go
	@echo ">> Installing goimports..."
	@GOBIN=$(TOOLS_BIN_DIR) $(GOCMD) install golang.org/x/tools/cmd/goimports

$(BUF_PATH): go.mod tools.go
	@echo ">> Installing buf..."
	@GOBIN=$(TOOLS_BIN_DIR) $(GOCMD) install github.com/bufbuild/buf/cmd/buf

# ====================================================================================
# 主要开发流程 (Main Development Workflow)
# ====================================================================================

# Default target
all: help

# Run the server
run: build
	@echo ">> Running server..."
	./$(BINARY_NAME)
	@echo "<< Server stopped."

# Build the binary
build: api wire
	@echo ">> Building binary..."
	$(GOBUILD) -o $(BINARY_NAME) ./cmd/user-srv
	@echo "<< Binary built."

# Test all packages
test:
	@echo ">> Running tests..."
	$(GOTEST) -v -cover ./...
	@echo "<< Tests completed."

# ====================================================================================
# 代码生成 (Code Generation)
# ====================================================================================

# Generate code from proto files
# .PHONY: api
# api: tools
# 	@echo ">> Generating api code from proto..."
# 	@protoc -I. -I./third_party \
# 	       --go_out=. --go_opt=paths=source_relative \
# 	       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
# 	       --grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative \
# 	       $(API_PROTO_FILES)
# 	@echo "<< API code generated."

# Generate dependency injection code with Wire
.PHONY: wire
wire: tools
	@echo ">> Generating wire code..."
	@$(WIRE_PATH) ./cmd/user-srv
	@echo "<< Wire code generated."

# Generate mocks
.PHONY: mock
mock: tools
	@echo ">> Generating mocks..."
	@$(MOCKGEN_PATH) -source=./internal/user-srv/biz/user.go -destination=./internal/user-srv/biz/mock/mocker_user.go -package=mock
	@echo "<< Mocks generated."


# ====================================================================================
# 代码质量 (Code Quality)
# ====================================================================================

# Format code
.PHONY: format
format: tools
	@echo ">> Formatting code..."
	@find . -name "*.go" -not -path "./bin/*" -not -path "./third_party/*" | xargs gofmt -s -w
	@find . -name "*.go" -not -path "./bin/*" -not -path "./third_party/*" | xargs $(GOIMPORTS_PATH) -w
	@echo "<< Code formatted."

# Run code linter
.PHONY: lint
lint: tools 
	@echo ">> Running linter..."
	@$(GOLANGCI_LINT_PATH) run ./...
	@echo "<< Linting completed."

# Run linter with auto-fix
.PHONY: lint-fix
lint-fix: tools format
	@echo ">> Running linter with auto-fix..."
	@$(GOLANGCI_LINT_PATH) run --fix ./...
	@echo "<< Auto-fix completed."

# ====================================================================================
# Protobuf
# ====================================================================================
# 生成API代码
.PHONY: api
api:
	@echo ">> Generating api code with Buf..."
	@$(BUF_PATH) generate
	@echo "<< API code generated."

# 运行proto文件的lint检查
.PHONY: proto-lint
proto-lint:
	@echo ">> Linting proto files..."
	@$(BUF_PATH) lint
	@echo "<< Proto files linted."
# 格式化proto文件
.PHONY: proto-format
proto-format:
	@echo ">> Formatting proto files..."
	@$(BUF_PATH) format -w
	@echo "<< Proto files formatted."
# 更新proto依赖
.PHONY: proto-dep-update
proto-dep-update:
	@echo ">> Updating proto dependencies..."
	@$(BUF_PATH) dep update
	@echo "<< Proto dependencies updated in buf.lock."
# 初始化buf配置
.PHONY: buf-config-init
buf-config-init:
	@echo ">> Initializing buf configuration..."
	@$(BUF_PATH) config init
	@echo "<< Buf configuration initialized."

# ====================================================================================
# npm
# ====================================================================================
.PHONY: npm-install
npm-install:
	@echo ">> Installing npm packages..."
	@npm install
	@echo "<< Npm packages installed."

# ====================================================================================
# 工具与清理 (Utilities & Cleanup)
# ====================================================================================

# Check tool versions
.PHONY: version
version: tools
	@echo ">> Checking tool versions:"
	@echo "-------------------------"
	@echo -n "protoc-gen-go:       "; $(PROTOC_GEN_GO_PATH) --version
	@echo -n "protoc-gen-go-grpc:  "; $(PROTOC_GEN_GO_GRPC_PATH) --version
	@echo -n "wire:                "; $(WIRE_PATH) --version
	@echo -n "mockgen:             "; $(MOCKGEN_PATH) --version
	@echo -n "golangci-lint:       "; $(GOLANGCI_LINT_PATH) --version
	@echo "-------------------------"


# Clean build artifacts
.PHONY: clean
clean:
	@echo ">> Cleaning up..."
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -rf $(TOOLS_BIN_DIR)
	@echo "<< Cleaned up."


# Help message
.PHONY: help
help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Main Targets:"
	@echo "  run           Build and run the server"
	@echo "  build         Build the binary without running"
	@echo "  test          Run all tests"
	@echo ""
	@echo "Code Generation:"
	@echo "  tools         Install/update development tools locally (from go.mod)"
	@echo "  api           Generate Go code from .proto files"
	@echo "  wire          Generate Go dependency injection code"
	@echo "  mock          Generate mock code for interfaces"
	@echo "  proto-lint    Lint .proto files using Buf"
	@echo "  proto-format  Format .proto files using Buf"
	@echo "  proto-dep-update  Update proto dependencies in buf.lock"
	@echo "  buf-config-init   Initialize Buf configuration file"
	@echo ""
	@echo "Code Quality:"
	@echo "  format        Format all Go code in the project"
	@echo "  lint          Run code linter to check for issues"
	@echo "  lint-fix      Run linter and automatically fix issues"
	@echo ""
	@echo "Utilities:"
	@echo "  version       Show installed local tool versions"
	@echo "  clean         Clean up build artifacts and installed tools"
	@echo "  help          Show this help message"