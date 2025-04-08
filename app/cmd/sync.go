package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yikakia/nga_grep/handler"
)

var (
	cid string
	uid string
	url string
	db  string
)

func init() {
	syncCmd.Flags().StringVar(&cid, "cid", "cid", "cookie 中的 ngaPassportCid")
	syncCmd.Flags().StringVar(&uid, "uid", "uid", "cookie 中的 ngaPassportUid")
	syncCmd.Flags().StringVar(&url, "url", "https://bbs.nga.cn", "nga域名")
	syncCmd.Flags().StringVar(&db, "db", "./nga.db", "db 路径")
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "爬取帖子",
	Run: func(cmd *cobra.Command, args []string) {
		handler.SyncServer(handler.SyncServerConfig{
			Cid: cid,
			Uid: uid,
			Url: url,
			DB:  db,
		})
	},
}
