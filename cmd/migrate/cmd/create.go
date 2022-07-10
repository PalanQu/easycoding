package cmd

import (
	"easycoding/common/workspace"
	"easycoding/pkg/orm"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/sunary/sqlize"
	sqlize_utils "github.com/sunary/sqlize/utils"
)

var createCmd = &cobra.Command{
	Use:   "create [table]",
	Short: "Create migrate sql files",
	RunE: func(cmd *cobra.Command, args []string) error {
		if migrateAll {
			for _, model := range models {
				if err := migrateModel(model); err != nil {
					return err
				}
			}
			return nil
		}
		if len(args) == 0 {
			return errors.New("please input a table name or use --all to create all tables")
		}
		tableName := args[0]
		for _, model := range models {
			if model.TableName() == tableName {
				if err := migrateModel(model); err != nil {
					return err
				}
			}
		}
		return fmt.Errorf("invalid table name %s", tableName)
	},
}

func migrateModel(model orm.Model) error {
	folder := filepath.Join(workspace.GetWorkspace(), migrationDir, model.TableName())
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		os.Mkdir(folder, os.ModePerm)
	}
	oldMigration := sqlize.NewSqlize(sqlize.WithMigrationFolder(folder))
	if err := oldMigration.FromMigrationFolder(); err != nil {
		return err
	}
	newMigration := sqlize.NewSqlize(
		sqlize.WithSqlTag("gorm"), sqlize.WithMigrationFolder(folder))
	if err := newMigration.FromObjects(model); err != nil {
		return err
	}
	newMigration.Diff(*oldMigration)
	if dryrun {
		verboseMigrationFile(newMigration)
		return nil
	}
	fmt.Printf("new files %s\n", filepath.Join(folder, model.TableName()))
	if err := writeFiles(folder, model.TableName(), newMigration); err != nil {
		return err
	}
	return nil
}

func writeFiles(dir, name string, migration *sqlize.Sqlize) error {
	fileName := sqlize_utils.MigrationFileName(name)
	upSubfix := ".up.sql"
	upFileName := fileName + upSubfix

	if err := ioutil.WriteFile(
		filepath.Join(dir, upFileName), []byte(migration.StringUp()), 0644); err != nil {
		return err
	}

	downSubfix := ".down.sql"
	downFileName := fileName + downSubfix
	if err := ioutil.WriteFile(
		filepath.Join(dir, downFileName), []byte(migration.StringDown()), 0644); err != nil {
		return err
	}
	return nil
}

func verboseMigrationFile(migration *sqlize.Sqlize) {
	fmt.Println("==================")
	fmt.Println("migration up")
	fmt.Println(migration.StringUp())
	fmt.Println("==================")
	fmt.Println("migration down")
	fmt.Println(migration.StringDown())
	fmt.Println("==================")
}

func initCreate() {
	rootCmd.AddCommand(createCmd)
}
