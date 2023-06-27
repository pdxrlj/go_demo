package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "remote_share_windows",
	Short: "remote_share_windows",
	Long:  "remote_share_windows",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Run() error {
	return rootCmd.Execute()
}
