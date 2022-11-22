package cmd

import (
	c "easycoding/internal/config"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	defaultMigrationDir = "./migrations"
)

var (
	dryrun       bool
	migrationDir string
	migrateAll   bool
	configPath   string

	config *c.Config
)

var rootCmd = &cobra.Command{
	Use:   "migrate",
	Short: "migrate database",
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			fmt.Println(err)
		}
		os.Exit(1)
	},
}

type migrateLogger struct {
	logger *logrus.Logger
}

func newMigrateLogger() migrate.Logger {
	return &migrateLogger{
		logger: logrus.New(),
	}
}

var _ migrate.Logger = (*migrateLogger)(nil)

func (logger *migrateLogger) Printf(format string, v ...interface{}) {
	logger.logger.Infof(format, v...)
}

func (logger *migrateLogger) Verbose() bool {
	return true
}

func InitCmd() error {
	rootCmd.PersistentFlags().BoolVarP(&dryrun, "dry-run", "n", false, "--dry-run")
	rootCmd.PersistentFlags().BoolVar(&migrateAll, "all", false, "--all")
	rootCmd.PersistentFlags().StringVarP(
		&migrationDir, "migrate-path", "p", defaultMigrationDir, "--migrate-path")
	rootCmd.PersistentFlags().StringVar(
		&configPath, "config", "./config.yaml", "--config")
	config = c.LoadConfig(configPath)
	initStep()
	initDiff()
	initTo()
	return initVersion()
}

func Execute() error {
	return rootCmd.Execute()
}
