package env

import (
	"strings"
)

// dev production
const DEPLOYMENT ENV = "DEPLOYMENT"

// 没有设置则默认为生产
func IsProduction() bool {
	return DEPLOYMENT.Get() == "" || strings.ToLower(DEPLOYMENT.Get()) == "production"
}
