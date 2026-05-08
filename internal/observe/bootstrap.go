package observe

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
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

func OTelAccessLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		InitLoggerOnce()

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
