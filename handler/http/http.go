package http

import (
	"log/slog"

	"github.com/yikakia/nga_grep/handler"
)

type RunHttpServerConfig struct {
	Port            string
	CorsAllowOrigin []string
	DB              string
}

func RunHttpServer(cfg RunHttpServerConfig) {
	handler.InitDefaultDB(cfg.DB)

	r, err := newGinEngine(cfg)
	if err != nil {
		panic(err)
	}

	// 监听 /my-path 路径
	r.GET("/api/timeseries", timeSeries)

	slog.Info("apiserver start...")
	err = r.Run(cfg.Port)
	if err != nil {
		panic(err)
	}
}
