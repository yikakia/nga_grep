package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yikakia/nga_grep/handler"
)

var (
	port            string
	corsAllowOrigin []string
	rDB             string
)

func init() {
	runHttpServerCmd.PersistentFlags().StringVarP(&port, "port", "p", ":11648", "Port to listen on")
	runHttpServerCmd.PersistentFlags().StringSliceVar(&corsAllowOrigin, "cors", []string{"localhost"}, "CORS Allow Origin 的域名, 逗号分割")
	runHttpServerCmd.PersistentFlags().StringVar(&rDB, "db", "", "db 路径")
}

var runHttpServerCmd = &cobra.Command{
	Use:   "api-server",
	Short: "启动http服务, 监听指定端口",
	PreRun: func(cmd *cobra.Command, args []string) {
		err := cmd.MarkFlagRequired("db")
		if err != nil {
			panic(err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		handler.RunHttpServer(handler.RunHttpServerConfig{
			Port:            port,
			CorsAllowOrigin: corsAllowOrigin,
			DB:              rDB,
		})
	},
}
