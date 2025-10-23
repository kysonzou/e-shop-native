package middleware

import (
	"net/http"

	"github.com/kyson/e-shop-native/pkg/trace"
)

const TraceIDHeader = "X-Trace-ID"

// TraceMiddleware 是一个 HTTP 中间件，用于处理 trace_id
func TraceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 尝试从请求头获取 trace_id
		traceID := r.Header.Get(TraceIDHeader)
		if traceID == "" {
			// 如果没有，就生成一个新的
			traceID = trace.NewTraceID()
		}

		// 将 trace_id 存入 context
		ctx := trace.ToContext(r.Context(), traceID)
		
		// 将带有 trace_id 的新 context 传递下去
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}