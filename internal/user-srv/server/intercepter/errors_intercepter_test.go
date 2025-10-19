package intercepter_test

import (
	"context"
	"errors"
	"testing"

	v1 "github.com/kyson/e-shop-native/api/protobuf/user/v1"
	apperrors "github.com/kyson/e-shop-native/internal/user-srv/errors"
	intercepter "github.com/kyson/e-shop-native/internal/user-srv/server/intercepter"
	"github.com/kyson/e-shop-native/pkg/ecode"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestErrorInterceptor(t *testing.T) {
	// 1. Arrange (准备)
	interceptor := intercepter.ErrorInterceptor
	info := &grpc.UnaryServerInfo{FullMethod: "/test.Service/TestMethod"}
	defaultInternalError := ecode.New(500, "internal server error", codes.Internal)

	tests := []struct {
		name            string
		mockHandler     grpc.UnaryHandler
		expectedResp    any
		expectedErrCode codes.Code
		expectedErrMsg  string
		checkDetails    bool
		expectedDetails *v1.UserErr
	}{
		{
			name: "Handler returns no error",
			mockHandler: func(ctx context.Context, req any) (any, error) {
				return "success", nil
			},
			expectedResp:    "success",
			expectedErrCode: codes.OK,
		},
		{
			name: "Handler returns a standard Go error",
			mockHandler: func(ctx context.Context, req any) (any, error) {
				return nil, errors.New("a generic error")
			},
			expectedResp:    nil,
			expectedErrCode: defaultInternalError.GrpcCode(),
			expectedErrMsg:  defaultInternalError.Message(),
			checkDetails:    true,
			expectedDetails: defaultInternalError.Detail(),
		},
		{
			name: "Handler returns a custom ecode error",
			mockHandler: func(ctx context.Context, req any) (any, error) {
				return nil, apperrors.ErrUserNotFound
			},
			expectedResp:    nil,
			expectedErrCode: apperrors.ErrUserNotFound.GrpcCode(),
			expectedErrMsg:  apperrors.ErrUserNotFound.Message(),
			checkDetails:    true,
			expectedDetails: apperrors.ErrUserNotFound.Detail(),
		},
		{
			name: "Handler returns a grpc status error",
			mockHandler: func(ctx context.Context, req any) (any, error) {
				return nil, status.Error(codes.PermissionDenied, "permission denied")
			},
			expectedResp:    nil,
			expectedErrCode: defaultInternalError.GrpcCode(), // ecode.FromError 将其视为普通 error
			expectedErrMsg:  defaultInternalError.Message(),
			checkDetails:    true,
			expectedDetails: defaultInternalError.Detail(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 2. Act (执行)
			resp, err := interceptor(context.Background(), "fake request", info, tt.mockHandler)

			// 3. Assert (断言)
			assert.Equal(t, tt.expectedResp, resp, "Response should match expected")

			if tt.expectedErrCode == codes.OK {
				assert.NoError(t, err, "Expected no error")
			} else {
				require.Error(t, err, "Expected an error")

				st, ok := status.FromError(err)
				require.True(t, ok, "Error should be a gRPC status error")

				assert.Equal(t, tt.expectedErrCode, st.Code(), "gRPC status code should match")
				assert.Equal(t, tt.expectedErrMsg, st.Message(), "gRPC status message should match")

				if tt.checkDetails {
					details := st.Details()
					require.NotEmpty(t, details, "Expected error details")

					detail, ok := details[0].(*v1.UserErr)
					require.True(t, ok, "Detail should be of type *v1.UserErr")

					assert.Equal(t, tt.expectedDetails.Code, detail.Code, "Detail code should match")
					assert.Equal(t, tt.expectedDetails.Message, detail.Message, "Detail message should match")
				}
			}
		})
	}
}
