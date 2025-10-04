package sever

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)



func RecoverInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {

	fmt.Printf("RecoverInterceptor called for method: %s\n", info.FullMethod)
	return handler(ctx, req)
}