package ratelimit

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	ctx = context.Background()
	key = "key"
)

func TestMemAllow(t *testing.T) {

	assert.True(t, memeRL().allow(key, 1))
	assert.True(t, memeRL().allow(key, 69999))
	assert.False(t, memeRL().allow(key, 1))
}

//func TestRedisAllow(t *testing.T) {
//	allowed, _, err := redisRL().AllowN(ctx, key, 1000)
//	assert.NoError(t, err)
//	assert.True(t, allowed)
//
//	allowed, _, err = redisRL().AllowN(ctx, key, 70000)
//	assert.NoError(t, err)
//	assert.False(t, allowed)
//}
