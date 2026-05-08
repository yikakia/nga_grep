package observe

import (
	"context"
	"log/slog"
	"os"
	"sync"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
)

var _initLogger = sync.OnceValues(func() (*log.LoggerProvider, error) {
	ctx := context.Background()
	exp, err := otlploghttp.New(ctx)
	if err != nil {
		return nil, err
	}
	processor := log.NewBatchProcessor(exp)
	provider := log.NewLoggerProvider(log.WithProcessor(processor))
	global.SetLoggerProvider(provider)

	otelSlogHandler := otelslog.NewHandler("nga", otelslog.WithLoggerProvider(provider), otelslog.WithSource(true))

	final := slog.NewMultiHandler(otelSlogHandler, slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	}))

	slog.SetDefault(slog.New(final))

	return provider, nil
})

func InitLoggerOnce() (*log.LoggerProvider, error) {
	return _initLogger()
}
