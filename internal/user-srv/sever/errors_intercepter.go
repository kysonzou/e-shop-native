package sever

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

func ErrorInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {

	fmt.Printf("ErrorInterceptor called for method: %s\n", info.FullMethod)
	return handler(ctx, req)
}
