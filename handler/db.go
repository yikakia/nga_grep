package handler

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/yikakia/nga_grep/client"
	"github.com/yikakia/nga_grep/internal/env"
	"github.com/yikakia/nga_grep/model/gen"
	"gorm.io/plugin/opentelemetry/tracing"
)

var (
	defaultDBOnce sync.Once
	defaultDBPath string
)

// InitDefaultDB 初始化 gorm/gen 的默认 DB。
//
// 由于 HTTP 与 sync 可能并发启动，为避免对 [`gen.SetDefault()`](model/gen/gen.go:24) 的并发写造成数据竞争，
// 这里使用 sync.Once 确保只初始化一次。
//
// 注意：若在同一进程中传入不同的 dbPath，将直接 panic。
func InitDefaultDB(dbPath string) {
	defaultDBOnce.Do(func() {
		defaultDBPath = dbPath
		db := client.NewDB(dbPath)

		var traceOpts []tracing.Option
		// 开发环境才显示查询参数
		if env.IsProduction() {
			slog.Info("disable query variables for non-dev env")
			traceOpts = append(traceOpts, tracing.WithoutQueryVariables())
		}
		err := db.Use(tracing.NewPlugin(traceOpts...))
		if err != nil {
			panic(err)
		}

		gen.SetDefault(db)
	})

	if defaultDBPath != dbPath {
		panic(fmt.Errorf("default db path mismatch: first=%q now=%q", defaultDBPath, dbPath))
	}
}
