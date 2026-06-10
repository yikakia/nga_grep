package handler

import (
	"fmt"
	"sync"

	"github.com/yikakia/nga_grep/client"
	"github.com/yikakia/nga_grep/model/gen"
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
		gen.SetDefault(client.NewDB(dbPath))
	})

	if defaultDBPath != dbPath {
		panic(fmt.Errorf("default db path mismatch: first=%q now=%q", defaultDBPath, dbPath))
	}
}
