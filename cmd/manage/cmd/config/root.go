package config

import (
	"os"

	"github.com/spf13/cobra"
)

const (
	dryrunKeyName = "dry-run"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage the config.yaml files",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
	},
}

func Init(rootCmd *cobra.Command) {
	rootCmd.AddCommand(configCmd)
	initGenerate()
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
