package http

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yikakia/nga_grep/internal/ratelimit"
)

func isAllow(c *gin.Context, start, end time.Time, duration time.Duration) (isAllow bool) {
	ctx := c.Request.Context()
	defer func() { recordRateLimit(ctx, isAllow) }()

	if duration <= 0 {
		slog.WarnContext(ctx, fmt.Sprintf("duration should > 0 but got %v", duration))
		return false
	}

	key := strings.Join([]string{"rl", c.ClientIP()}, ":")

	cost := end.Sub(start)/duration + 1
	if cost <= 0 {
		slog.WarnContext(ctx, "end should > start.", slog.Time("start", start), slog.Time("end", end))
		return false
	}

	slog.DebugContext(ctx, "call http allow.", "key", key, "cost", cost)

	return ratelimit.HTTPAllow(ctx, key, int(cost))
}
