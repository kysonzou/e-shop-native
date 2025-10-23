package intercepter

import (
	"context"

	"github.com/kyson/e-shop-native/pkg/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const TraceIDKey = "X-Trace-ID"

func TraceServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
    // 步骤 1: 从传入的 context 中提取 gRPC metadata
    // gRPC 框架已经帮我们把请求头解析好，并放在了 context 中。
    // 我们用 metadata.FromIncomingContext 来“解包”它。
    md, ok := metadata.FromIncomingContext(ctx)
    if !ok {
        // 如果没有 metadata，创建一个空的，以防后续代码 panic
        md = metadata.New(nil)
    }

    // 步骤 2: 从 metadata 中读取 trace_id
    // metadata 本质上是一个 map[string][]string
    var traceID string
    if vals := md.Get(TraceIDKey); len(vals) > 0 {
        // 如果上游服务传来了 trace_id，我们沿用它
        traceID = vals[0]
    } else {
        // 如果没有，说明这是调用链的起点，我们生成一个新的
        traceID = trace.NewTraceID()

		// 将traceId写到metadata中
        // 注意：Incoming metadata 是不可变的，所以我们需要创建一个它的副本
        mdCopy := md.Copy()
        mdCopy.Set(TraceIDKey, traceID)
        
        // 创建一个新的 context，它携带了这个更新后的 metadata
        // 这样，后续的中间件或 handler 调用 metadata.FromIncomingContext 就能拿到 trace-id 了
        ctx = metadata.NewIncomingContext(ctx, mdCopy)
    }

    // 步骤 3: 将 trace_id 注入到 Go 的 context 中
    // 这里就是 context.WithValue 的用武之地。我们创建了一个新的 context，
    // 它携带了 trace_id。
    newCtx := trace.ToContext(ctx, traceID)

    // 步骤 4: 调用下一个 handler，并将这个“增强后”的 newCtx 传递下去
    // 从这里开始，这个请求在我们的服务内部的所有函数调用（service, biz, data），
    // 只要它们接收 ctx 参数，就都能通过 trace.FromContext(ctx) 获取到 trace_id。
    return handler(newCtx, req)
}

func TraceClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
    // 步骤 1: 从当前 Go 的 context 中提取 trace_id
    // 这个 trace_id 是在入站拦截器中被放进去的。
    traceID, ok := trace.FromContext(ctx)
    if !ok {
        // 如果 context 中没有，可以生成一个新的，或者留空
        traceID = trace.NewTraceID()
    }

    // 步骤 2: 将 trace_id 添加到 gRPC 的出站 metadata 中
    // 我们需要修改即将被发送出去的请求的 context，把 metadata “附加”上去。
    // metadata.AppendToOutgoingContext 会创建一个新的 context，其中包含了要发送的头部信息。
    newCtx := metadata.AppendToOutgoingContext(ctx, TraceIDKey, traceID)

    // 步骤 3: 调用真正的 RPC invoker，并将这个“增强后”的 newCtx 传递下去
    // gRPC 客户端在发送请求时，会自动从 newCtx 中提取出站 metadata，
    // 并将其转换为 HTTP/2 的请求头，发送给下游服务。
    return invoker(newCtx, method, req, reply, cc, opts...)
}