package intercepter

import (
	"context"

	code "github.com/kyson/e-shop-native/pkg/code"

	"google.golang.org/grpc"
)

func ErrorInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	resp, err = handler(ctx, req)
	// 这里必须把所有的错误全部转换为grpc的标准错误
	if err != nil {
		return resp, code.FromError(err).GrpcError()
	}
	return resp, err
}
