package loghelper

import (
	"context"

	"github.com/kyson/e-shop-native/pkg/trace"
	"go.uber.org/zap"
)

func FromContext(ctx context.Context, logger *zap.Logger) *zap.Logger{
	if traceID, ok := trace.FromContext(ctx); ok{
		logger = logger.With(zap.String("trace_id", traceID))
	}
	return logger
}

