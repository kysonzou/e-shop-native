package intercepter

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

func LoggingInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {

	fmt.Printf("LoggingInterceptor called for method: %s\n", info.FullMethod)
	return handler(ctx, req)
}
