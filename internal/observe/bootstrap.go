package observe

import (
	"log/slog"
	"time"

	"github.com/bytedance/gg/gcond"
	"github.com/gin-gonic/gin"
	"github.com/yikakia/nga_grep/internal/buildinfo"
	"github.com/yikakia/nga_grep/internal/env"
	"go.opentelemetry.io/otel/attribute"
)

func InitAll() error {
	if _, err := InitTracer(); err != nil {
		return err
	}
	if _, err := InitMeter(); err != nil {
		return err
	}
	if _, err := InitLoggerOnce(); err != nil {
		return err
	}
	return nil
}

func defaultAttributes() []attribute.KeyValue {
	kvs := []attribute.KeyValue{buildinfo.VCSAttribute()}
	kvs = append(kvs, attribute.String("deployment.environment", gcond.If(env.IsProduction(), "production", env.DEPLOYMENT.Get())))
	return kvs
}

func OTelAccessLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		InitLoggerOnce()

		start := time.Now()

		c.Next()

		slog.InfoContext(
			c.Request.Context(),
			"access log",
			"method", c.Request.Method,
			"host", c.Request.Host,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency_ms", time.Since(start).Milliseconds(),
			"client_ip", c.ClientIP(),
		)
	}
}
