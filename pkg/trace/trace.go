package trace

import (
	"context"

	"github.com/google/uuid"
)

func NewTraceID() string{
	return uuid.NewString()
}

type traceIdKey struct{}

func WithTraceID(ctx context.Context, traceId string) context.Context{
	return context.WithValue(ctx, traceIdKey{}, traceId)
}

func FromContext(ctx context.Context) (string, bool){
	id, ok := ctx.Value(traceIdKey{}).(string)
	return id, ok
}