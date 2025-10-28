package server

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/kyson/e-shop-native/internal/user-srv/conf"
)

func NewAdminServer(c *conf.Server) *AdminHTTPServer {
	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.Handler())

	http_server := &http.Server{
		Addr:    c.Admin.Addr,
		Handler: mux,
	}
	return &AdminHTTPServer{Server: http_server}
}
