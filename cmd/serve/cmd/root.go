package cmd

import (
	"easycoding/internal/app"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	configPath string
)

var rootCmd = &cobra.Command{
	Use:   "go-template-server",
	Short: "This is the main command",
	Run: func(cmd *cobra.Command, args []string) {
		runServer(cmd, args)
	},
}

func InitCmd() error {
	rootCmd.PersistentFlags().StringVar(
		&configPath, "config", "./config.yaml", "--config")
	if err := initServe(); err != nil {
		return err
	}
	if err := initVersion(); err != nil {
		return err
	}
	return nil
}

func Execute() error {
	return rootCmd.Execute()
}

func boot() *app.Kernel {
	kernel, err := app.New(configPath)
	if err != nil {
		logger := log.New(os.Stderr, "debug", 1)
		logger.Fatalf("failed to boot kernel: %s", err)
		os.Exit(2)
	}
	return kernel
}
