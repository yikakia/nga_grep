package observe

import (
	"context"
	"log/slog"
	"os"
	"sync"

	"github.com/yikakia/nga_grep/internal/buildinfo"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
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

	stdoutExporter, err := stdoutlog.New(stdoutlog.WithWriter(os.Stdout))
	if err != nil {
		return nil, err
	}
	stdoutProcessor := log.NewBatchProcessor(stdoutExporter)

	provider := log.NewLoggerProvider(log.WithProcessor(stdoutProcessor), log.WithProcessor(processor))
	global.SetLoggerProvider(provider)

	// 这里的name是 scope name
	otelSlogHandler := otelslog.NewHandler("",
		otelslog.WithLoggerProvider(provider), otelslog.WithSource(true),
		otelslog.WithAttributes(buildinfo.VCSAttribute()),
	)

	slog.SetDefault(slog.New(otelSlogHandler))

	return provider, nil
})

func InitLoggerOnce() (*log.LoggerProvider, error) {
	return _initLogger()
}
