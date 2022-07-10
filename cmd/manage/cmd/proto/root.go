package proto

import (
	"os"

	"github.com/spf13/cobra"
)

const (
	dryrunKeyName = "dry-run"
)

var protoCmd = &cobra.Command{
	Use:   "proto",
	Short: "Manage the api/*.proto files",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(1)
		}
	},
}

func getDryrun(cmd *cobra.Command) bool {
	dryrun := true
	var err error
	dryrun, err = cmd.Flags().GetBool("dry-run")
	if err != nil {
		return true
	}
	return dryrun
}

func Init(rootCmd *cobra.Command) {
	rootCmd.AddCommand(protoCmd)
	initGo()
	initSwagger()
}
