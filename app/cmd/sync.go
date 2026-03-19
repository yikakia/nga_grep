package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	// sync 命令作为兼容入口：复用 api-server 的 sync 配置与实现。
	syncCmd.Flags().StringVar(&dbPath, "db", "./nga.db", "db 路径")
	addSyncFlags(syncCmd)
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "爬取帖子",
	RunE: func(cmd *cobra.Command, args []string) error {
		runApiServerWithModes(apiServerModes{Sync: true})
		return nil
	},
}
