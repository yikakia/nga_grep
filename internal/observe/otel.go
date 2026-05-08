package observe

import (
	"context"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/sdk/metric"
)

func InitAll() error {
	if _, err := initTracer(); err != nil {
		return err
	}
	if _, err := initMetric(); err != nil {
		return err
	}
	if _, err := initLogger(); err != nil {
		return err
	}
	return nil
}

var _initTracer = sync.OnceValues(func() (*sdktrace.TracerProvider, error) {
	exporter, err := otlptrace.New(context.Background(), otlptracehttp.NewClient())
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{}, propagation.Baggage{}))
	return tp, err
})

func initTracer() (*sdktrace.TracerProvider, error) {
	return _initTracer()
}

var _initMetric = sync.OnceValues(func() (*metric.MeterProvider, error) {
	exporter, err := otlpmetrichttp.New(context.Background())

	if err != nil {
		return nil, err
	}
	mp := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(exporter)),
	)
	otel.SetMeterProvider(mp) // 必须注册全局 MeterProvider
	return mp, nil
})

func initMetric() (*metric.MeterProvider, error) {
	return _initMetric()
}

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

func initLogger() (*log.LoggerProvider, error) {
	return _initLogger()
}

func GinMiddleware(r *gin.Engine) error {
	_, err := initTracer()
	if err != nil {
		return err
	}
	_, err = initMetric()
	if err != nil {
		return err
	}

	_, err = _initLogger()
	if err != nil {
		return err
	}

	r.Use(otelgin.Middleware("nga"))
	return nil
}

func OTelAccessLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		_initLogger()

		start := time.Now()

		c.Next()

		slog.InfoContext(
			c.Request.Context(),
			"access log",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency_ms", time.Since(start).Milliseconds(),
			"client_ip", c.ClientIP(),
		)
	}
}
