package logger

import (
	"context"

	"go.uber.org/zap"
)

type Interface interface {
	Info(ctx context.Context, message string, args ...zap.Field)
	Warn(ctx context.Context, message string, args ...zap.Field)
	Error(ctx context.Context, message string, args ...zap.Field)
	Fatal(message string, args ...zap.Field)
	Debug(ctx context.Context, message string, args ...zap.Field)
}
