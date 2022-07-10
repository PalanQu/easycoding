package cmd

import (
	"easycoding/internal/app"
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the current version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(app.Version())
	},
}

func initVersion() error {
	rootCmd.AddCommand(versionCmd)
	return nil
}
