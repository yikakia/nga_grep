package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yikakia/nga_grep/handler"
)

func init() {
	runHttpServerCmd.PersistentFlags().StringP("port", "p", ":11648", "Port to listen on")
	runHttpServerCmd.PersistentFlags().StringSlice("cors", []string{"localhost"}, "CORS Allow Origin 的域名, 逗号分割")
	runHttpServerCmd.PersistentFlags().String("db", "./nga.db", "db 路径")
}

var runHttpServerCmd = &cobra.Command{
	Use:   "api-server",
	Short: "启动http服务, 监听指定端口",
	Run: func(cmd *cobra.Command, args []string) {
		port, err := cmd.PersistentFlags().GetString("port")
		if err != nil {
			panic(err)
		}
		corsAllowOrigin, err := cmd.PersistentFlags().GetStringSlice("cors")
		if err != nil {
			panic(err)
		}
		db, err := cmd.PersistentFlags().GetString("db")
		if err != nil {
			panic(err)
		}
		handler.RunHttpServer(handler.RunHttpServerConfig{
			Port:            port,
			CorsAllowOrigin: corsAllowOrigin,
			DB:              db,
		})
	},
}
