package ratelimit

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yikakia/nga_grep/internal/env"
	"golang.org/x/time/rate"
)

func Init() {
	if env.REDIS_URL.Get() != "" {
		redisRL()
		return
	}
	memeRL()
}

// 平均每三天70000个点，最大每三天 70000 个点
func HTTPAllow(ctx context.Context, key string, cost int) bool {
	if env.REDIS_URL.Get() != "" {
		allow, _, err := redisRL().AllowN(ctx, key, cost)
		if err != nil {
			slog.WarnContext(ctx, "call redis failed", "err", err.Error())
			return false
		}
		return allow
	}

	return memeRL().allow(key, cost)
}

const _3DaySeconds = 3 * 24 * 60 * 60
const burst = 70000

var memeRL = sync.OnceValue(newRLStore)

// 平均每三天70000个点，最大每三天 70000 个点
func newLimiter() *rate.Limiter {
	return rate.NewLimiter(rate.Every((3*24*time.Hour)/burst), burst)
}

var redisRL = sync.OnceValue(func() *TokenBucket {
	opt, _ := redis.ParseURL(env.REDIS_URL.Get())
	client := redis.NewClient(opt)

	return NewTokenBucket(TokenBucketConfig{
		Capacity:   burst,
		RefillRate: float64(burst) / float64(_3DaySeconds) / 2,
		Client:     client,
	})
})
