package data

import (
	"context"
	"sync"
	"time"

	"github.com/dgraph-io/ristretto/v2"
	"github.com/yikakia/cachalot"
	"github.com/yikakia/cachalot/core/cache"
	store_ristretto "github.com/yikakia/cachalot/stores/ristretto"
)

var risStore = sync.OnceValue(func() cache.Store {
	client, err := ristretto.NewCache(&ristretto.Config[string, any]{
		NumCounters: 1 << 10,
		MaxCost:     1 << 20,
		BufferItems: 64,
	})
	if err != nil {
		panic(err)
	}
	return store_ristretto.New(client, store_ristretto.WithStoreName("basic-ristretto"))
})

var missLoader = func(ctx context.Context, key string, opts ...cache.CallOption) (int, error) {
	cfg := cache.ApplyOptions(opts...)
	t := cfg.CustomField["t"].(time.Time)
	duration := cfg.CustomField["duration"].(time.Duration)

	data, err := getTimePointData(t, duration)
	if err != nil {
		return 0, err
	}
	return data, nil
}

var cachalotCache = sync.OnceValue(func() cache.Cache[int] {
	b, err := cachalot.NewBuilder[int]("cache", risStore())
	if err != nil {
		panic(err)
	}
	b.WithCacheMissLoader(missLoader).
		WithLogicExpireLoader(missLoader)

	build, err := b.Build()
	if err != nil {
		panic(err)
	}

	return build
})
