package server

import (
	"net/http"

	"google.golang.org/grpc"
)

type BusinessGRPCServer struct {
	*grpc.Server
}
type BusinessHTTPServer struct {
	*http.Server
}
type AdminHTTPServer struct {
	*http.Server
}
