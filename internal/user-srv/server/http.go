package server

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	v1 "github.com/kyson/e-shop-native/api/protobuf/user/v1"
	"github.com/kyson/e-shop-native/internal/user-srv/conf"
	intercepter "github.com/kyson/e-shop-native/internal/user-srv/server/intercepter"
	middleware "github.com/kyson/e-shop-native/internal/user-srv/server/middleware"
)

func NewHTTPServer(c *conf.Server, logger *zap.Logger) (*BusinessHTTPServer, error) {
	// 初始化gateway
	mux := runtime.NewServeMux(runtime.WithErrorHandler(middleware.CustomErrorHandle(logger)))

	// GRPC客户端
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(intercepter.TraceClientInterceptor),
	}

	err := v1.RegisterUserServiceHandlerFromEndpoint(context.Background(), mux, c.GRPC.Addr, opts)
	if err != nil {
		return nil, err
	}

	// 这里可以直接把mux挂载到http.Server上，但是这样的话就不能实现中间件了，所以引入一个轻量http库
	chi := chi.NewRouter()
	//chi.Use() //可以挂载各种中间件
	chi.Use(middleware.TraceMiddleware)
	chi.Use(chiMiddleware.Logger)         // 记录请求的完整生命周期
	chi.Use(chiMiddleware.Recoverer)      // 终极保护，必须在最外层之一，捕获一切panic
	chi.Use(middleware.MetricsMiddleware) // 指标

	chi.Mount("/", mux) //把gateway挂载到chi上，也就是请求先到chi，然后chi再根据这里的挂载规则转发到gateway

	http_server := &http.Server{
		Addr:    c.HTTP.Addr,
		Handler: chi,
	}
	return &BusinessHTTPServer{Server: http_server}, nil
}
