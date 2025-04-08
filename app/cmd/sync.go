package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yikakia/nga_grep/handler"
)

var (
	cid                 string
	uid                 string
	url                 string
	db                  string
	thresholdLow        int
	thresholdHigh       int
	thresholdLowFactor  float64
	thresholdHighFactor float64
)

func init() {
	syncCmd.Flags().StringVar(&cid, "cid", "cid", "cookie 中的 ngaPassportCid")
	syncCmd.Flags().StringVar(&uid, "uid", "uid", "cookie 中的 ngaPassportUid")
	syncCmd.Flags().StringVar(&url, "url", "https://bbs.nga.cn", "nga域名")
	syncCmd.Flags().StringVar(&db, "db", "./nga.db", "db 路径")
	syncCmd.Flags().IntVar(&thresholdLow, "th-l", 10, "下限阈值 低于这个值，下次调度时间会修改")
	syncCmd.Flags().IntVar(&thresholdHigh, "th-h", 20, "上限阈值 高于这个值，下次调度时间会修改")
	syncCmd.Flags().Float64Var(&thresholdLowFactor, "th-lf", 2, "低于下限阈值后，调度时间的倍数")
	syncCmd.Flags().Float64Var(&thresholdHighFactor, "th-hf", 0.5, "高于上限阈值后，调度时间的倍数")
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "爬取帖子",
	Run: func(cmd *cobra.Command, args []string) {
		handler.SyncServer(handler.SyncServerConfig{
			Cid:                 cid,
			Uid:                 uid,
			Url:                 url,
			DB:                  db,
			ThresholdLow:        thresholdLow,
			ThresholdHigh:       thresholdHigh,
			ThresholdLowFactor:  thresholdLowFactor,
			ThresholdHighFactor: thresholdHighFactor,
		})
	},
}
