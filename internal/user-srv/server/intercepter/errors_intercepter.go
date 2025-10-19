package intercepter

import (
	"context"

	ecode "github.com/kyson/e-shop-native/pkg/ecode"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func ErrorInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	resp, err = handler(ctx, req)
	// 这里必须把所有的错误全部转换为grpc的标准错误
	if err != nil {
		ecodeErr, _ := ecode.FromError(err)
		st := status.New(ecodeErr.GrpcCode(), ecodeErr.Message()) // 创建基础 status

		// WithDetails 会返回一个新的 status 对象，必须接收它
		stWithDetails, detailErr := st.WithDetails(ecodeErr.Detail())
		if detailErr != nil {
			// 如果附加 detail 失败（虽然很少见），则返回不带 detail 的原始 status 错误
			return resp, st.Err()
		}
		return resp, stWithDetails.Err()
	}
	return resp, err
}
