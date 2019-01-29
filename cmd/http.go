package cmd

import (
	"github.com/g0lang/proxy/src/api"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(httpCmd)
}

var httpCmd = &cobra.Command{
	Use:   "run",
	Short: "run http server",
	Long:  `run server `,
	Run: func(cmd *cobra.Command, args []string) {
		api.Run()
	},
}
