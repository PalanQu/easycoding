package cmd

import (
	"database/sql"
	"easycoding/common/workspace"
	c "easycoding/internal/config"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
)

const (
	maxStep = 1000
)

var stepLatest bool
var stepReverse bool

var stepCmd = &cobra.Command{
	Use:   "step [step_number]",
	Short: "Migrate step",
	RunE: func(cmd *cobra.Command, args []string) error {
		step := 0
		if len(args) != 1 && !stepLatest {
			return errors.New("input the step number or use --latest")
		}
		if !stepLatest {
			var err error
			step, err = strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("step invalid %s, %s", args[0], err.Error())
			}
			if step > maxStep {
				return fmt.Errorf("max step %v", maxStep)
			}
		}
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
		for _, model := range models {
			dir := filepath.Join(workspace.GetWorkspace(), migrationDir, model.TableName())
			fileUri := fmt.Sprintf("file://%s", dir)
			m, err := migrate.NewWithDatabaseInstance(fileUri, "mysql", driver)
			m.Log = newMigrateLogger()
			if err != nil {
				return err
			}
			if stepLatest {
				if !stepReverse {
					if err := m.Up(); err != nil {
						return err
					}
				} else {
					if err := m.Down(); err != nil {
						return err
					}
				}
			} else {
				if !stepReverse {
					if err := m.Steps(step); err != nil {
						return err
					}
				} else {
					if err := m.Steps(step * -1); err != nil {
						return err
					}
				}
			}
		}
		return nil
	},
}

func initUp() {
	rootCmd.AddCommand(stepCmd)
	stepCmd.Flags().BoolVar(&stepLatest, "latest", false, "--latest")
	stepCmd.Flags().BoolVar(&stepReverse, "reverse", false, "--reverse")
	config = c.LoadConfig(configPath)
}
