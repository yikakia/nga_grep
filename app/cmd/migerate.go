package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yikakia/nga_grep/client"
	"github.com/yikakia/nga_grep/model"
)

var migrateDB string

func init() {
	migrateCmd.Flags().StringVar(&migrateDB, "db", "", "db 路径")
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "初始化数据库",
	PreRun: func(cmd *cobra.Command, args []string) {
		err := cmd.MarkFlagRequired("db")
		if err != nil {
			panic(err)
		}
		return
	},
	Run: func(cmd *cobra.Command, args []string) {
		db := client.NewDB(migrateDB)
		err := db.AutoMigrate(model.ThreadLatestData{}, model.ThreadCount{})
		if err != nil {
			panic(err)
		}
	},
}
