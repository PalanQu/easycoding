package cmd

import (
	"database/sql"
	"easycoding/common/workspace"
	c "easycoding/internal/config"
	"fmt"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Migrate version",
	RunE: func(cmd *cobra.Command, args []string) error {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			config.Database.User,
			config.Database.Password,
			config.Database.Host,
			config.Database.Port,
			config.Database.DBName,
		)
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			return err
		}
		driver, err := mysql.WithInstance(db, &mysql.Config{})
		if err != nil {
			return err
		}
		dir := filepath.Join(workspace.GetWorkspace(), migrationDir)
		fileUri := fmt.Sprintf("file://%s", dir)
		m, err := migrate.NewWithDatabaseInstance(fileUri, "mysql", driver)
		if err != nil {
			return err
		}
		v, dirty, err := m.Version()
		if err != nil {
			return err
		}
		fmt.Printf("Version: %v, Dirty: %v\n", v, dirty)
		return nil
	},
}

func initVersion() {
	rootCmd.AddCommand(versionCmd)
	config = c.LoadConfig(configPath)
}
