package sever

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	v1 "github.com/kyson/e-shop-native/api/protobuf/user/v1"
	"github.com/kyson/e-shop-native/internal/user-srv/conf"
)

func NewHTTPServer(c *conf.Server) *http.Server {
	// 初始化gateway
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	v1.RegisterUserServiceHandlerFromEndpoint(context.Background(), mux, c.GRPC.Addr, opts)

	// 这里可以直接把mux挂载到http.Server上，但是这样的话就不能实现中间件了，所以引入一个轻量http库
	chi := chi.NewRouter()
	//chi.Use() //可以挂载各种中间件

	chi.Mount("/", mux) //把gateway挂载到chi上

	return &http.Server{
		Addr:    c.HTTP.Addr,
		Handler: chi,
	}

}
