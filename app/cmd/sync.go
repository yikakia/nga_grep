package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yikakia/nga_grep/handler"
)

func init() {
	syncCmd.PersistentFlags().String("cid", "cid", "cookie 中的 ngaPassportCid")
	syncCmd.PersistentFlags().String("uid", "uid", "cookie 中的 ngaPassportUid")
	syncCmd.PersistentFlags().String("url", "https://bbs.nga.cn", "nga域名")
	syncCmd.PersistentFlags().String("db", "./nga.db", "db 路径")
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "爬取帖子",
	Run: func(cmd *cobra.Command, args []string) {

		cid, err := cmd.PersistentFlags().GetString("cid")
		if err != nil {
			panic(err)
		}
		uid, err := cmd.PersistentFlags().GetString("uid")
		if err != nil {
			panic(err)
		}
		url, err := cmd.PersistentFlags().GetString("url")
		if err != nil {
			panic(err)
		}
		db, err := cmd.PersistentFlags().GetString("db")
		if err != nil {
			panic(err)
		}
		handler.SyncServer(handler.SyncServerConfig{
			Cid: cid,
			Uid: uid,
			Url: url,
			DB:  db,
		})
	},
}
