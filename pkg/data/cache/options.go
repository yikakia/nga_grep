package cache

import (
	"time"
)

// WithExpiration 用于设置缓存的过期时间
// 0 表示使用默认的过期时间
func WithExpiration[T any](fn func(args map[string]any) time.Duration) func(*Cache[T]) error {
	return func(c *Cache[T]) error {
		c.expirationFn = fn
		return nil
	}
}
