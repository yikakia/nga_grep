package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/yikakia/nga_grep/handler"
)

var (
	// http
	port            string
	corsAllowOrigin []string

	// common
	dbPath string

	// api-server mode
	mode string

	// sync
	cid                 string
	uid                 string
	url                 string
	thresholdLow        int
	thresholdHigh       int
	thresholdLowFactor  float64
	thresholdHighFactor float64
	loopMin             time.Duration
	loopMax             time.Duration
)

func init() {
	// http
	runHttpServerCmd.Flags().StringVarP(&port, "port", "p", ":11648", "Port to listen on")
	runHttpServerCmd.Flags().StringSliceVar(&corsAllowOrigin, "cors", []string{"localhost"}, "CORS Allow Origin 的域名, 逗号分割")

	// common
	runHttpServerCmd.Flags().StringVar(&dbPath, "db", "", "db 路径")
	runHttpServerCmd.Flags().StringVar(&mode, "mode", "http", "启动模式，逗号分隔：http,sync")

	// sync flags (used when --mode contains sync)
	addSyncFlags(runHttpServerCmd)

	if err := runHttpServerCmd.MarkFlagRequired("db"); err != nil {
		panic(err)
	}
}

type apiServerModes struct {
	HTTP bool
	Sync bool
}

func parseApiServerModes(s string) (apiServerModes, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return apiServerModes{}, fmt.Errorf("--mode 不能为空，可选：http,sync（逗号分隔）")
	}

	var m apiServerModes
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		switch part {
		case "http":
			m.HTTP = true
		case "sync":
			m.Sync = true
		default:
			return apiServerModes{}, fmt.Errorf("非法 --mode=%q，可选：http,sync（逗号分隔）", part)
		}
	}

	if !m.HTTP && !m.Sync {
		return apiServerModes{}, fmt.Errorf("--mode 必须至少包含一个：http 或 sync")
	}
	return m, nil
}

func runApiServerWithModes(m apiServerModes) {
	switch {
	case m.Sync && m.HTTP:
		go handler.SyncServer(buildSyncConfig())
		handler.RunHttpServer(buildHttpConfig())
	case m.HTTP:
		handler.RunHttpServer(buildHttpConfig())
	case m.Sync:
		handler.SyncServer(buildSyncConfig())
	}
}

func buildHttpConfig() handler.RunHttpServerConfig {
	return handler.RunHttpServerConfig{
		Port:            port,
		CorsAllowOrigin: corsAllowOrigin,
		DB:              dbPath,
	}
}

func buildSyncConfig() handler.SyncServerConfig {
	return handler.SyncServerConfig{
		Cid:                 cid,
		Uid:                 uid,
		Url:                 url,
		DB:                  dbPath,
		ThresholdLow:        thresholdLow,
		ThresholdHigh:       thresholdHigh,
		ThresholdLowFactor:  thresholdLowFactor,
		ThresholdHighFactor: thresholdHighFactor,
		LoopMin:             loopMin,
		LoopMax:             loopMax,
	}
}

func addSyncFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&cid, "cid", "cid", "cookie 中的 ngaPassportCid")
	cmd.Flags().StringVar(&uid, "uid", "uid", "cookie 中的 ngaPassportUid")
	cmd.Flags().StringVar(&url, "url", "https://bbs.nga.cn", "nga域名")
	cmd.Flags().IntVar(&thresholdLow, "th-l", 10, "下限阈值 低于这个值，下次调度时间会修改")
	cmd.Flags().IntVar(&thresholdHigh, "th-h", 20, "上限阈值 高于这个值，下次调度时间会修改")
	cmd.Flags().Float64Var(&thresholdLowFactor, "th-lf", 2, "低于下限阈值后，调度时间的倍数")
	cmd.Flags().Float64Var(&thresholdHighFactor, "th-hf", 0.5, "高于上限阈值后，调度时间的倍数")
	cmd.Flags().DurationVar(&loopMin, "loop-min", time.Second*30, "最小调度时间")
	cmd.Flags().DurationVar(&loopMax, "loop-max", time.Minute*8, "最大调度时间")
}

var runHttpServerCmd = &cobra.Command{
	Use:   "api-server",
	Short: "启动服务（HTTP API / 同步爬取）",
	RunE: func(cmd *cobra.Command, args []string) error {
		modes, err := parseApiServerModes(mode)
		if err != nil {
			return err
		}
		runApiServerWithModes(modes)
		return nil
	},
}
