package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "nga_grep",
		Short: "nga爬虫&展示",
		Long:  "nga爬虫&展示",
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(runHttpServerCmd)
	rootCmd.AddCommand(syncCmd)
	rootCmd.AddCommand(migrateCmd)
}
