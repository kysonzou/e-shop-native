//go:build tools
// +build tools

package tools

import (
	// 接口生成工具，测试时使用
	_ "github.com/golang/mock/mockgen"

	// 依赖注入工具
	_ "github.com/google/wire/cmd/wire"

	// grpc-gateway 和 protobuf 相关工具
	// _ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
	// _ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	// _ "google.golang.org/protobuf/cmd/protoc-gen-go"

	// 代码质量检查工具
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"

	// 代码规范
	_ "golang.org/x/tools/cmd/goimports"

	// buf 代码生成工具
	_ "github.com/bufbuild/buf/cmd/buf"
)
