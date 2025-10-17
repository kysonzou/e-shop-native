package sever_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/kyson/e-shop-native/internal/user-srv/auth"
	"github.com/kyson/e-shop-native/internal/user-srv/conf"
	"github.com/kyson/e-shop-native/internal/user-srv/sever"
	//apperrors "github.com/kyson/e-shop-native/internal/user-srv/errors"
)

// TestAuthInterceptor_Unit a pure unit test for the interceptor logic.
func TestAuthInterceptor_Unit(t *testing.T) {
	// 1. Arrange (准备)
	mockConfig := &conf.Auth{
		Algorithm:      "HS256",
		JwtKey:         "ahkPzSJ6auFD2WZHt5NFfixFSI3JmXm4isbTs8y29Zs=",
		ExpireDuration: 3600,
		Whitelist:      []string{"/test.Service/PublicMethod"},
	}
	authInstance, err := auth.NewAuth(mockConfig)
	require.NoError(t, err)

	// 获取拦截器函数
	interceptor := sever.AuthInterceptor(authInstance)

	// --- 定义我们的测试用例 ---

	tests := []struct {
		name                  string
		fullMethod            string          // 模拟 info.FullMethod
		ctx                   context.Context // 模拟传入的 context
		handlerShouldBeCalled bool            // 预期 handler 是否会被调用
		expectedErrCode       codes.Code      // 预期返回的 gRPC 错误码
		checkClaimsInCtx      bool            // 是否需要在 handler 中检查 claims
	}{
		{
			name:                  "Whitelisted method should pass through",
			fullMethod:            "/test.Service/PublicMethod",
			ctx:                   context.Background(),
			handlerShouldBeCalled: true,
			expectedErrCode:       codes.OK,
		},
		{
			name:                  "Protected method without metadata should fail",
			fullMethod:            "/test.Service/ProtectedMethod",
			ctx:                   context.Background(),
			handlerShouldBeCalled: false,
			expectedErrCode:       codes.Unauthenticated,
		},
		{
			name:                  "Protected method without auth header should fail",
			fullMethod:            "/test.Service/ProtectedMethod",
			ctx:                   metadata.NewIncomingContext(context.Background(), metadata.MD{}),
			handlerShouldBeCalled: false,
			expectedErrCode:       codes.Unauthenticated,
		},
		{
			name:                  "Protected method with invalid auth header format should fail",
			fullMethod:            "/test.Service/ProtectedMethod",
			ctx:                   metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "invalid format")),
			handlerShouldBeCalled: false,
			expectedErrCode:       codes.Unauthenticated,
		},
		{
			name:                  "Protected method with invalid token should fail",
			fullMethod:            "/test.Service/ProtectedMethod",
			ctx:                   metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "Bearer invalid-token")),
			handlerShouldBeCalled: false,
			expectedErrCode:       codes.Unauthenticated,
		},
		{
			name:       "Protected method with valid token should pass and inject claims",
			fullMethod: "/test.Service/ProtectedMethod",
			ctx: func() context.Context {
				token, _ := authInstance.GenerateToken(context.Background(), 42, "testuser")
				return metadata.NewIncomingContext(context.Background(), metadata.Pairs("authorization", "Bearer "+token))
			}(),
			handlerShouldBeCalled: true,
			expectedErrCode:       codes.OK,
			checkClaimsInCtx:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 2. Act (执行)

			handlerCalled := false // 用于跟踪 handler 是否被调用

			// 创建一个假的 handler
			mockHandler := func(ctx context.Context, req any) (any, error) {
				handlerCalled = true

				// 如果需要，检查 context 中是否存在 claims
				if tt.checkClaimsInCtx {
					claims, ok := auth.FromContext(ctx)
					assert.True(t, ok, "Handler context should contain claims")
					assert.NotNil(t, claims)
					assert.Equal(t, uint(42), claims.Id, "Claims should have correct user ID")
				}

				return "handler response", nil
			}

			// 模拟 grpc.UnaryServerInfo
			info := &grpc.UnaryServerInfo{
				FullMethod: tt.fullMethod,
			}

			// 直接调用拦截器函数
			_, err := interceptor(tt.ctx, "fake request", info, mockHandler)

			// 3. Assert (断言)

			// 断言 handler 的调用情况
			assert.Equal(t, tt.handlerShouldBeCalled, handlerCalled, "Handler call status mismatch")

			// 断言返回的错误
			if tt.expectedErrCode == codes.OK {
				assert.NoError(t, err)
			} else {
				require.Error(t, err)
				st, _ := status.FromError(err)
				assert.Equal(t, tt.expectedErrCode, st.Code())
			}
		})
	}
}
