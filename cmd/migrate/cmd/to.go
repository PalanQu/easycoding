package cmd

import (
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"

	"easycoding/common/workspace"
	pkg_db "easycoding/pkg/db"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/spf13/cobra"
)

var toCmd = &cobra.Command{
	Use:   "to [version]",
	Short: "to specific version",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("please input the target version")
		}
		toVersion, err := strconv.Atoi(args[0])
		if err != nil {
			return errors.New("the version must be a number")
		}
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?multiStatements=true",
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
		currentVersion, dirty, err := m.Version()
		if err != nil {
			return err
		}
		if dirty {
			return errors.New("database is dirty, patch the diff first")
		}
		direction, step, err := pkg_db.MigrationGenerate(dir, int(currentVersion), toVersion)
		if err != nil {
			return err
		}
		if direction == pkg_db.MigrationDirectionUP {
			if err := m.Steps(step); err != nil {
				return err
			}
		} else {
			if err := m.Steps(step * -1); err != nil {
				return err
			}
		}
		return nil
	},
}

func initTo() {
	rootCmd.AddCommand(toCmd)
}
