package intercepter

import (
	"context"
	//"fmt"
	//"net/http"
	//"strconv"
	"time"

	//"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

var (
	GRPCRequestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_request_total",
			Help: "Total number of gRPC requests",
		},
		[]string{"method", "code"},
	)
	GRPCRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_request_duration_seconds",
			Help:    "gRPC request latencies in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)
)

func init() {
	// 目的是让 prometheus.Handler 知道有它们的存在，未来通过metric可以获取到对应的数据
	prometheus.MustRegister(GRPCRequestTotal)
	prometheus.MustRegister(GRPCRequestDuration)
}

func MetricsInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	start := time.Now()

	resp, err = handler(ctx, req)

	duration := time.Since(start)
	statusCode := status.Code(err).String()

	// 记录指标
	GRPCRequestTotal.WithLabelValues(info.FullMethod, statusCode).Inc()
	GRPCRequestDuration.WithLabelValues(info.FullMethod).Observe(duration.Seconds())
	return resp, err
}
