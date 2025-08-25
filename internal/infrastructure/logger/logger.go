package logger

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func InitLogger() error {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{
		"logs/service.log",
		"stderr",
	}
	config.ErrorOutputPaths = []string{
		"logs/error.log",
		"stderr",
	}
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	var err error
	Log, err = config.Build()
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	return nil
}

func SyncLogger() {
	_ = Log.Sync()
}

func InitTestLogger() error {
	config := zap.NewDevelopmentConfig()
	config.OutputPaths = []string{"stdout"}
	config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)

	logger, err := config.Build()
	if err != nil {
		return fmt.Errorf("failed to initialize test logger: %w", err)
	}

	Log = logger

	return nil
}

type ctxLoggerKey struct{}

func WithLogger(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxLoggerKey{}, l)
}

func FromContext(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return Log
	}
	if l, ok := ctx.Value(ctxLoggerKey{}).(*zap.Logger); ok && l != nil {
		return l
	}

	return Log
}

func WithTraceID(ctx context.Context, l *zap.Logger) *zap.Logger {
	if l == nil {
		return l
	}

	sc := trace.SpanContextFromContext(ctx)
	if sc.TraceID().IsValid() {
		return l.With(zap.String("traceID", sc.TraceID().String()))
	}

	return l
}
