package http

import (
	"log/slog"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/yikakia/nga_grep/internal/observe"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func newGinEngine(cfg RunHttpServerConfig) (*gin.Engine, error) {
	r := gin.New()

	config := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"}, // 允许的方法
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},      // 允许的头
		ExposeHeaders:    []string{"Content-Length"},                                // 暴露的头
		AllowCredentials: true,
		AllowOriginWithContextFunc: func(c *gin.Context, origin string) bool {
			for _, s := range cfg.CorsAllowOrigin {
				if strings.Contains(origin, s) {
					return true
				}
			}
			slog.WarnContext(c.Request.Context(), "hit cors for origin:"+origin)
			return false
		},
		MaxAge: 12 * time.Hour, // 预检请求的缓存时间
	}

	middlewares := []gin.HandlerFunc{cors.New(config)}

	err := observe.InitAll()
	if err != nil {
		panic(err)
	}

	middlewares = append(middlewares, otelgin.Middleware("nga-api"))

	middlewares = append(middlewares, gin.Logger(), observe.OTelAccessLogMiddleware(), gin.Recovery())

	middlewares = append(middlewares, func(c *gin.Context) {
		ctx := c.Request.Context()
		// 打印所有的 header
		attrs := make([]slog.Attr, 0, len(c.Request.Header))
		for k, v := range c.Request.Header {
			attrs = append(attrs, slog.String(k, strings.Join(v, ",")))
		}

		slog.DebugContext(ctx, "request headers", slog.GroupAttrs("headers", attrs...))
		c.Next()

	})

	r.Use(middlewares...)

	return r, nil
}
