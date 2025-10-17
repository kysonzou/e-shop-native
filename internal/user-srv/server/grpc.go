package server

import (
	v1 "github.com/kyson/e-shop-native/api/protobuf/user/v1" // Update to the correct import path for your generated gRPC code
	"github.com/kyson/e-shop-native/internal/user-srv/auth"
	"github.com/kyson/e-shop-native/internal/user-srv/conf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewGRPCServer(c *conf.Server, src v1.UserServiceServer, auth auth.Auth) *BusinessGRPCServer {
	// options
	opts := grpc.ChainUnaryInterceptor(
		RecoverInterceptor,
		LoggingInterceptor,
		AuthInterceptor(auth),
		ErrorInterceptor,
		MetricsInterceptor,
	)

	// Create the gRPC server
	server := grpc.NewServer(opts)

	// Register your gRPC services here
	v1.RegisterUserServiceServer(server, src)

	// Enable reflection for debugging (optional)
	reflection.Register(server)

	// Return the gRPC server instance
	return &BusinessGRPCServer{Server: server}
}
