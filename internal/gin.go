package internal

import (
	"time"

	"github.com/gin-gonic/gin"
)

type afterMiddlewareWriter struct {
	gin.ResponseWriter
	c     *gin.Context
	start time.Time
	f     func(c *gin.Context)
}

func (w *afterMiddlewareWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.Header().Set("X-Response-Time", time.Since(w.start).String())
	w.ResponseWriter.WriteHeader(statusCode)
}

func ResponseTimeHeader() func(*gin.Context) {
	return func(c *gin.Context) {
		c.Writer = &afterMiddlewareWriter{ResponseWriter: c.Writer, c: c, start: time.Now()}
		c.Next()
	}
}
