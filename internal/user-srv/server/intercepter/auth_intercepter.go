package intercepter

import (
	"context"
	"slices"
	"strings"

	"github.com/kyson/e-shop-native/internal/user-srv/auth"
	apperrors "github.com/kyson/e-shop-native/internal/user-srv/errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
    // AuthorizationHeader is the header name for authorization.
    AuthorizationHeader = "authorization"
    // BearerScheme is the prefix for bearer tokens.
    BearerScheme = "Bearer"
)

func AuthInterceptor(a auth.Auth) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		// 白名单
		whiteList := a.GetWhiteList()

		if slices.Contains(whiteList, info.FullMethod) {
			return handler(ctx, req)
		}

		// 解析token
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, apperrors.ErrTokenInvalid.WithMessage("authentication required").GrpcError()
		}

		authHeaders := md.Get(AuthorizationHeader)
		if len(authHeaders) == 0 {
			return nil, apperrors.ErrTokenInvalid.WithMessage("authentication required").GrpcError()
		}

		parts := strings.Split(authHeaders[0], " ")
		if len(parts) != 2 || parts[0] != BearerScheme {
			return nil, apperrors.ErrTokenInvalid.WithMessage("invalid token format").GrpcError()
		}

		tokenString := parts[1]
		if tokenString == "" {
			return nil, apperrors.ErrTokenInvalid.WithMessage("token is empty").GrpcError()
		}

		// 判断token的
		ctx, err = a.ParseAndSaveToken(ctx, tokenString)
		if err != nil {
			return nil, apperrors.ErrTokenInvalid.WithMessage("invalid or expired token")
		}

		return handler(ctx, req)
	}
}
