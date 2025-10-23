package intercepter

import (
	"context"
	//"fmt"
	"runtime/debug"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RecoverInterceptor(log *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		defer func() {
			if r := recover(); r != nil {
				// 捕获到panic
				log.Error("panic recovered",
					zap.Any("panic", r),
					zap.String("stacktrace", string(debug.Stack())),
				)
				// 返回一个grpc标准错误（这里使用的是命名返回值的形式）
				// 等同于 return nil, err
				err = status.Errorf(codes.Internal, "内部服务错误")
			}
		}()
		// 正常调用下一个handler
		return handler(ctx, req)
	}
}
