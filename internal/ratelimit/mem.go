package ratelimit

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type rlEntry struct {
	rl       *rate.Limiter
	lastSeen time.Time
}

func newRLStore() *rlStore {
	rl := &rlStore{
		m: make(map[string]*rlEntry),
	}
	go func() {
		defer func() {
			v := recover()
			slog.Info(fmt.Sprintf("panic for rl cleanup. detail:%v", v))
		}()
		rl.cleanup()
	}()
	return rl
}

type rlStore struct {
	m  map[string]*rlEntry
	mu sync.Mutex
}

func (r *rlStore) allow(key string, cost int) bool {
	e := r.get(key)

	return e.rl.AllowN(time.Now(), cost)
}

func (r *rlStore) get(key string) *rlEntry {
	r.mu.Lock()
	defer r.mu.Unlock()

	e, ok := r.m[key]
	if ok {
		e.lastSeen = time.Now()
		return e
	}
	e = &rlEntry{
		rl:       newLimiter(),
		lastSeen: time.Now(),
	}
	r.m[key] = e
	return e
}

func (r *rlStore) cleanup() {
	for {
		time.Sleep(time.Minute)
		now := time.Now()

		r.mu.Lock()
		total := 0
		cnt := 0
		for k, v := range r.m {
			if now.Sub(v.lastSeen) > 48*time.Hour {
				delete(r.m, k)
				cnt++
			}
			total++

			// 查 20 个 如果超时的大于 33% 则继续，否则退出
			if total < 20 {
				continue
			}
			if cnt*3 >= total {
				total = 0
				cnt = 0
			} else {
				break
			}
		}
		r.mu.Unlock()
	}
}
