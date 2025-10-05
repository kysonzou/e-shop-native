package sever

import (
	"context"
	"slices"
	"strings"

	apperrors "github.com/kyson/e-shop-native/internal/user-srv/errors"	
	"github.com/kyson/e-shop-native/internal/user-srv/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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
			return nil, apperrors.ErrTokenInvalid
		}

		authHeaders := md.Get("authorization")
		if len(authHeaders) == 0 {
			return nil, apperrors.ErrTokenInvalid
		}

		parts := strings.Split(authHeaders[0], " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return nil, apperrors.ErrTokenInvalid
		}

		tokenString := parts[1]
		// 判断token的
		ctx, err = a.ParseAndSaveToken(ctx, tokenString)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}
