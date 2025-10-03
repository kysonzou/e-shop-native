package sever

import (
	v1 "github.com/kyson/e-shop-native/api/protobuf/user/v1" // Update to the correct import path for your generated gRPC code
	"github.com/kyson/e-shop-native/internal/user-srv/conf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func NewGRPCServer(c *conf.Server, src v1.UserServiceServer) *grpc.Server {

	// Create the gRPC server
	server := grpc.NewServer()
	// Register your gRPC services here
	v1.RegisterUserServiceServer(server, src)

	// Enable reflection for debugging (optional)
	reflection.Register(server)

	// Return the gRPC server instance
	return server
}