package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var dryrun bool
var rootCmd = &cobra.Command{
	Use:   "lint",
	Short: "Lint source code",
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			fmt.Println(err)
		}
		os.Exit(1)
	},
}

func InitCmd() error {
	rootCmd.PersistentFlags().BoolVarP(&dryrun, "dry-run", "n", false, "--dry-run")
	initGo()
	initProto()
	return nil
}

func Execute() error {
	return rootCmd.Execute()
}
