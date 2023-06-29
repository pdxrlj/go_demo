package cmd

import (
	"github.com/spf13/cobra"

	"remote_share_windows/config"
	"remote_share_windows/server"
)

var rootCmd = &cobra.Command{
	Use:   "remote_share_windows",
	Short: "remote_share_windows",
	Long:  "remote_share_windows",
	Run: func(cmd *cobra.Command, args []string) {
		server.NewServer(config.C).Run()
	},
}

func Run() error {
	return rootCmd.Execute()
}

func init() {
	c := config.C
	rootCmd.PersistentFlags().StringVar(&c.Addr, "addr", ":8090", "服务器监听地址")
}
