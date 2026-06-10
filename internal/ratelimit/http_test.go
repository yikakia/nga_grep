package ratelimit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPAllow(t *testing.T) {
	key := "key"
	assert.True(t, HTTPAllow(key, 1))
	assert.True(t, HTTPAllow(key, 9999))
	assert.False(t, HTTPAllow(key, 1))
}
