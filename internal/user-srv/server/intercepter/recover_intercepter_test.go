package intercepter_test

import (
	"context"
	"errors"
	"testing"

	intercepter "github.com/kyson/e-shop-native/internal/user-srv/server/intercepter"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// 单次逻辑测试
func TestRecoverInterceptor_Unit(t *testing.T) {
	// 使用 observer 来验证日志
	core, logs := observer.New(zap.ErrorLevel)
	logger := zap.New(core)
	interceptor := intercepter.RecoverInterceptor(logger)

	tests := []struct {
		name               string
		mockHandle         func(ctx context.Context, req any) (resp any, err error)
		expectErr          bool
		expectInternalCode bool
		expectResp         any
		expectLogCount     int
	}{
		{
			name: "测试string类型panic",
			mockHandle: func(ctx context.Context, req any) (resp any, err error) {
				panic("test panic")
			},
			expectErr:          true,
			expectInternalCode: true,
			expectLogCount:     1,
		},
		{
			name: "测试error类型panic",
			mockHandle: func(ctx context.Context, req any) (resp any, err error) {
				panic(errors.New("test error panic"))
			},
			expectErr:          true,
			expectInternalCode: true,
			expectLogCount:     1,
		},
		{
			name: "测试Handle正常调用",
			mockHandle: func(ctx context.Context, req any) (resp any, err error) {
				return "called", nil
			},
			expectErr:      false,
			expectResp:     "called",
			expectLogCount: 0,
		},
		{
			name: "测试Handle返回业务错误",
			mockHandle: func(ctx context.Context, req any) (resp any, err error) {
				return nil, status.Error(codes.InvalidArgument, "invalid argument")
			},
			expectErr:      true,
			expectLogCount: 0,
		},
		{
			name: "测试Handle返回普通错误",
			mockHandle: func(ctx context.Context, req any) (resp any, err error) {
				return nil, errors.New("test error")
			},
			expectErr:      true,
			expectLogCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 重置日志计数
			logs.TakeAll()

			info := &grpc.UnaryServerInfo{
				FullMethod: "/test.Service/TestMethod",
			}

			resp, err := interceptor(context.Background(), nil, info, tt.mockHandle)

			// 验证错误状态
			if (err != nil) != tt.expectErr {
				t.Fatalf("expect error status mismatch. Got error: %v, Expected error: %v", err, tt.expectErr)
			}

			// 验证错误码
			if tt.expectInternalCode {
				s, ok := status.FromError(err)
				if !ok {
					t.Errorf("expected grpc status error, got: %v", err)
				} else if s.Code() != codes.Internal {
					t.Errorf("expected Internal error code, got: %v", s.Code())
				}
				// 验证错误消息
				if s.Message() != "internal server error" {
					t.Errorf("expected error message 'internal server error', got: %v", s.Message())
				}
			}

			// 验证正常返回值
			if resp != tt.expectResp {
				t.Errorf("expected response %v, got %v", tt.expectResp, resp)
			}

			// 验证日志记录
			logEntries := logs.TakeAll()
			if len(logEntries) != tt.expectLogCount {
				t.Errorf("expected %d log entries, got %d", tt.expectLogCount, len(logEntries))
			}

			// 如果期望有日志，验证日志内容
			if tt.expectLogCount > 0 && len(logEntries) > 0 {
				logEntry := logEntries[0]

				if logEntry.Message != "panic recovered" {
					t.Errorf("expected log message 'panic recovered', got: %v", logEntry.Message)
				}

				// 验证日志字段
				fields := logEntry.ContextMap()
				if _, ok := fields["panic"]; !ok {
					t.Error("expected 'panic' field in log")
				}
				if _, ok := fields["stacktrace"]; !ok {
					t.Error("expected 'stacktrace' field in log")
				}
			}
		})
	}
}

// 测试多次调用
func TestRecoverInterceptor_MultipleCalls(t *testing.T) {
	logger := zap.NewNop()
	interceptor := intercepter.RecoverInterceptor(logger)
	info := &grpc.UnaryServerInfo{
		FullMethod: "/test.Service/TestMethod",
	}

	// 第一次调用：正常
	resp, err := interceptor(context.Background(), nil, info, func(ctx context.Context, req any) (resp any, err error) {
		return "first", nil
	})
	if err != nil || resp != "first" {
		t.Errorf("first call failed: resp=%v, err=%v", resp, err)
	}

	// 第一次调用：panic
	_, err = interceptor(context.Background(), nil, info, func(ctx context.Context, req any) (resp any, err error) {
		panic("second panic")
	})
	if err == nil {
		t.Error("expected error from panic, got nil")
	}

	// 第三次调用：正常（验证拦截器是否可以正常工作）
	resp, err = interceptor(context.Background(), nil, info, func(ctx context.Context, req any) (resp any, err error) {
		return "third", nil
	})
	if err != nil || resp != "third" {
		t.Errorf("third call failed: resp=%v, err=%v", resp, err)
	}
}

// 性能测试
func BenchmarkRecoverInterceptor(b *testing.B) {
	logger := zap.NewExample()
	interceptor := intercepter.RecoverInterceptor(logger)
	info := &grpc.UnaryServerInfo{
		FullMethod: "/test.Service/TestMethod",
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return "response", nil
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := interceptor(context.Background(), nil, info, handler)
		if err != nil {
			b.Fatalf("interceptor failed: %v", err)
		}
	}
}
