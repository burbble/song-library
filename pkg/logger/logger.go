package logger

import (
	"context"
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	RequestIdTrackerKey = "request_id"
)

type Logger struct {
	logger *zap.Logger
}

func New(level string) *Logger {
	var l zapcore.Level

	switch strings.ToLower(level) {
	case "error":
		l = zap.ErrorLevel
	case "warn":
		l = zap.WarnLevel
	case "info":
		l = zap.InfoLevel
	case "debug":
		l = zap.DebugLevel
	default:
		l = zap.InfoLevel
	}

	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(l),
		Development: false,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "msg",
			LevelKey:       "level",
			TimeKey:        "ts",
			CallerKey:      "caller",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	logger, _ := cfg.Build()
	defer func() {
		_ = logger.Sync()
	}()

	return &Logger{
		logger: logger,
	}
}

func (l *Logger) Debug(ctx context.Context, message string, fields ...zap.Field) {
	fields = l.appendRequestId(ctx, fields...)
	l.log("debug", message, fields...)
}

func (l *Logger) Info(ctx context.Context, message string, fields ...zap.Field) {
	fields = l.appendRequestId(ctx, fields...)
	l.log("info", message, fields...)
}

func (l *Logger) Warn(ctx context.Context, message string, fields ...zap.Field) {
	fields = l.appendRequestId(ctx, fields...)
	l.log("warn", message, fields...)
}

func (l *Logger) Error(ctx context.Context, message string, fields ...zap.Field) {
	fields = l.appendRequestId(ctx, fields...)
	l.log("error", message, fields...)
}

func (l *Logger) Fatal(message string, fields ...zap.Field) {
	l.log("fatal", message, fields...)
	os.Exit(1)
}

func (l *Logger) appendRequestId(ctx context.Context, fields ...zap.Field) []zap.Field {
	requestIdTracker := ctx.Value(RequestIdTrackerKey)
	if requestIdTracker != nil {
		fields = append(fields, zap.String(RequestIdTrackerKey, requestIdTracker.(string)))
	}

	return fields
}

func (l *Logger) log(level string, message string, fields ...zap.Field) {
	switch strings.ToLower(level) {
	case "debug":
		l.logger.Debug(message, fields...)
	case "info":
		l.logger.Info(message, fields...)
	case "warn":
		l.logger.WithOptions(zap.AddCallerSkip(2)).Warn(message, fields...)
	case "error":
		l.logger.WithOptions(zap.AddCallerSkip(2)).Error(message, fields...)
	case "fatal":
		l.logger.WithOptions(zap.AddCallerSkip(2)).Fatal(message, fields...)
	default:
		l.logger.WithOptions(zap.AddCallerSkip(2)).Info(fmt.Sprintf("%s message: %v has unknown level %v", level, message, level))
	}
}
