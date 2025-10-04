package sever

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

var (
	HttpRequestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_total",
			Help: "Total number of HTTP requests",	
			// ConstLabels: prometheus.Labels{ // 通过它可以携带各种常量信息
			// 	"service": "user-srv",
			// },
		},
		[]string{"method", "path", "code"}, // 动态变量，需要在具体实现的时候给它赋值
	)	

	HttpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latencies in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

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
	prometheus.MustRegister(HttpRequestTotal)
	prometheus.MustRegister(HttpRequestDuration)
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
	fmt.Printf("MetricsInterceptor called for method: %s\n", info.FullMethod)
	return resp, err
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// 这是一个装饰器，自带的http.ResponseWriter它只有写的属性，但是我们需要读取它的状态码
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		duration := time.Since(start)	
		fmt.Printf("MetricsMiddleware called for method: %s\n", r.URL.Path)
		HttpRequestTotal.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(ww.Status())).Inc()
		HttpRequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration.Seconds())
	})
}
