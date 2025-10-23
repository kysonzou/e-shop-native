package server

import (
	v1 "github.com/kyson/e-shop-native/api/protobuf/user/v1" // Update to the correct import path for your generated gRPC code
	"github.com/kyson/e-shop-native/internal/user-srv/auth"
	"github.com/kyson/e-shop-native/internal/user-srv/conf"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	intercepter "github.com/kyson/e-shop-native/internal/user-srv/server/intercepter"

)

func NewGRPCServer(c *conf.Server, src v1.UserServiceServer, auth auth.Auth, log *zap.Logger) *BusinessGRPCServer {
	// options
	opts := grpc.ChainUnaryInterceptor(
		intercepter.TraceServerInterceptor,
		intercepter.LoggingInterceptor,
		intercepter.RecoverInterceptor(log),
		intercepter.MetricsInterceptor,
		intercepter.AuthInterceptor(auth),
		intercepter.ErrorInterceptor,
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
