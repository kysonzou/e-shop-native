package middleware

import (
	"context"
	"io"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	v1 "github.com/kyson/e-shop-native/api/protobuf/user/v1"
	"go.uber.org/zap"
	//"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CustomErrorHandle 是一个 grpc-gateway 的错误处理器，
// 它将 gRPC 错误转换为自定义的 JSON 错误响应。

func CustomErrorHandle (logger *zap.Logger) func(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error){

	return func (ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
		// 1. 将传入的 error 转换为 gRPC 的 status 对象
		s := status.Convert(err)

		// 2. 从 status 中提取 gRPC code 和 message
		code := s.Code().String()
		msg := s.Message()

		// 尝试从 status 中提取 details。
		// 在我们的 ErrorInterceptor 中，我们将 ecode.Detail() 作为 details 附加了。
		if len(s.Details()) > 0 {
			details := s.Details()[0]
			biz_err, ok := details.(*v1.UserErr)
			if ok {
				code = biz_err.Code
				msg = biz_err.Message
			}
		}

		// 3. 将 gRPC code 映射为 HTTP status code
		httpStatus := runtime.HTTPStatusFromCode(s.Code())

		// 4. 组装自定义的 JSON 错误响应体
		// 这个结构可以根据您的前端需求进行调整
		httpErr := &v1.UserErr{
			Code:    code,
			Message: msg,
		}

		// 5. 使用 marshaler 将自定义错误结构序列化为 JSON
		// 设置响应体数据类型JOSN
		w.Header().Set("Content-Type", marshaler.ContentType(httpErr)) 

		// 设置HTTP status code
		w.WriteHeader(httpStatus)
		
		// 使用 marshaler 将我们自定义的 httpErr 结构体序列化成 JSON 字符串，并写入到 HTTP 响应体中。
		// 这里还考虑如果JOSN序列化失败的回退
		if err := marshaler.NewEncoder(w).Encode(httpErr); err != nil { 
			logger.Error("Failed to marshal and write error response", zap.Error(err))
			// 如果序列化失败，提供一个最终的回退
			w.WriteHeader(http.StatusInternalServerError)

			const fallbackErrorResponse = `{"code":-1, "message":"An internal error occurred while processing the error response"}`
			if _, wErr := io.WriteString(w, fallbackErrorResponse); wErr != nil {
				logger.Error("Failed to write fallback error response", zap.Error(err))
			}
		}
	}
}


