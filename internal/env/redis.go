package env

import (
	"os"
)

type ENV string

func (e ENV) Get() string {
	return os.Getenv(string(e))
}

const (
	REDIS_URL ENV = "REDIS_URL"
)
