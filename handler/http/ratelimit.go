package http

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/extra/redisotel-native/v9"
	"github.com/redis/go-redis/v9"
	"github.com/sourcegraph/conc/panics"
	"github.com/yikakia/nga_grep/internal/env"
	"github.com/yikakia/nga_grep/internal/observe"
	"github.com/yikakia/nga_grep/internal/ratelimit"
	"go.opentelemetry.io/otel/attribute"
	"golang.org/x/time/rate"
)

func isAllow(c *gin.Context, start, end time.Time, duration time.Duration) (isAllow bool) {
	ctx := c.Request.Context()
	defer func() { recordRateLimit(ctx, isAllow) }()
	ctx, sp := observe.Start(ctx, "ratelimit")
	defer sp.End()

	if duration <= 0 {
		slog.WarnContext(ctx, fmt.Sprintf("duration should > 0 but got %v", duration))
		return false
	}

	key := strings.Join([]string{"rl", c.ClientIP()}, ":")

	cost := int(end.Sub(start)/duration + 1)
	if cost <= 0 {
		slog.WarnContext(ctx, "end should > start.", slog.Time("start", start), slog.Time("end", end))
		return false
	}
	sp.SetAttributes(attribute.Int("cost", cost))

	slog.DebugContext(ctx, "call http allow.", "key", key, "cost", cost)
	// 如果查询量小，不是明细，则直接放行，但是计入额度
	if cost < 1000 {
		nctx := context.WithoutCancel(ctx)
		go panics.Try(func() {
			doAllow(nctx, key, cost)
		})
		return true
	}

	return doAllow(ctx, key, cost)
}

func initRateLimiter() {
	if env.REDIS_URL.Get() != "" {
		redisRL()
		return
	}
	memRL()
}

// 平均每三天70000个点，最大每三天 70000 个点
func doAllow(ctx context.Context, key string, cost int) bool {
	if env.REDIS_URL.Get() != "" {
		allow, _, err := redisRL().AllowN(ctx, key, cost)
		if err != nil {
			slog.WarnContext(ctx, "call redis failed", "err", err.Error())
			// 退化到本地限流
			return memRL().AllowN(key, cost)
		}
		return allow
	}

	return memRL().AllowN(key, cost)
}

const _3DaySeconds = 3 * 24 * 60 * 60
const burst = 70000

var memRL = sync.OnceValue(func() *ratelimit.RlStore {
	// 平均每三天70000个点，最大每三天 70000 个点
	return ratelimit.NewRLStore(rate.Every((3*24*time.Hour)/burst), burst)
})

var redisRL = sync.OnceValue(func() *ratelimit.TokenBucket {
	opt, _ := redis.ParseURL(env.REDIS_URL.Get())
	client := redis.NewClient(opt)

	observe.InitAll()
	regRedisOTEL()

	client.Ping(context.Background())

	slog.Info("redis client init")
	return ratelimit.NewTokenBucket(ratelimit.TokenBucketConfig{
		Capacity:   burst,
		RefillRate: float64(burst) / float64(_3DaySeconds) / 2,
		Client:     client,
	})
})

func regRedisOTEL() {
	i := redisotel.GetObservabilityInstance()
	config := redisotel.NewConfig().WithEnabled(true)

	err := i.Init(config)
	if err != nil {
		slog.Error("redis init with otel failed", "err", err.Error())
		return
	}
}
