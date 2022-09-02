package db

import (
	"os"

	c "easycoding/internal/config"

	"github.com/spf13/cobra"
)

const (
	defaultMigrationDir = "./migrations"
	defaultConfigPath   = "./cmd/manage/cmd/db_migrate/config.yaml"
)

var (
	dryrun       bool
	migrationDir string
	configPath   string

	config *c.Config
)

var dbMigrateCmd = &cobra.Command{
	Use:   "db-migrate",
	Short: "Manage database migration files",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
	},
}

func Init(rootCmd *cobra.Command) {
	rootCmd.AddCommand(dbMigrateCmd)
	rootCmd.PersistentFlags().StringVarP(
		&migrationDir, "migrate-path", "p", defaultMigrationDir, "--migrate-path")
	rootCmd.PersistentFlags().StringVar(
		&configPath, "config", defaultConfigPath, "--config")
	config = c.LoadConfig(configPath)
	initGenerate()
}
