package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/samber/lo"
	"golang.org/x/sync/singleflight"
)

type Cache[T any] struct {
	*singleflight.Group
	getter Getter[T]
	name   KeyName
	*cache.Cache
}

type Getter[T any] func(args map[string]any) (T, error)
type KeyName func(args map[string]any) (string, error)

func NewCache[T any](getter Getter[T], name KeyName) *Cache[T] {
	return &Cache[T]{
		Group:  &singleflight.Group{},
		getter: getter,
		name:   name,
		Cache:  cache.New(time.Minute, time.Minute*5),
	}
}

func (c *Cache[T]) Get(args map[string]any) (T, error) {
	key, err := c.name(args)
	if err != nil {
		return lo.Empty[T](), err
	}

	v, exit := c.Cache.Get(key)
	if exit {
		return v.(T), nil
	}
	v, err, _ = c.Group.Do(key, func() (any, error) {
		v, err := c.getter(args)
		if err != nil {
			return nil, err
		}
		c.Cache.Set(key, v, cache.DefaultExpiration)
		return v, nil
	})
	if err != nil {
		return lo.Empty[T](), err
	}

	return v.(T), nil
}
