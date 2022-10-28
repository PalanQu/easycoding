package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"easycoding/common/workspace"
	"easycoding/pkg/ent"

	"entgo.io/ent/dialect/sql/schema"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/spf13/cobra"
)

var applyDiffToDB bool

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Show diff between schema and db",
	RunE: func(cmd *cobra.Command, args []string) error {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?multiStatements=true",
			config.Database.User,
			config.Database.Password,
			config.Database.Host,
			config.Database.Port,
			config.Database.DBName,
		)
		client, err := ent.Open("mysql", dsn)
		if err != nil {
			return err
		}
		defer client.Close()
		migrateOptions := []schema.MigrateOption{
			schema.WithDropColumn(true),
			schema.WithDropIndex(true),
		}
		if err := printDiff(client, migrateOptions); err != nil {
			log.Fatalf("failed printing schema changes: %v", err)
			return err
		}
		if applyDiffToDB {
			if err := applyDiff(client, migrateOptions); err != nil {
				log.Fatalf("failed apply diff: %v", err)
				return err
			}
			if err := dirtyToFalse(dsn); err != nil {
				log.Fatalf("failed change dirty to false: %v", err)
				return err
			}
			log.Println("apply success")
		}
		return nil
	},
}

func printDiff(client *ent.Client, migrateOptions []schema.MigrateOption) error {
	if err := client.Schema.WriteTo(
		context.Background(),
		os.Stdout,
		migrateOptions...,
	); err != nil {
		return err
	}
	return nil
}

func applyDiff(client *ent.Client, migrateOptions []schema.MigrateOption) error {
	if err := client.Schema.Create(
		context.Background(),
		migrateOptions...,
	); err != nil {
		return err
	}
	return nil
}

func dirtyToFalse(dsn string) error {
	dir := filepath.Join(workspace.GetWorkspace(), migrationDir)
	fileUri := fmt.Sprintf("file://%s", dir)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(fileUri, "mysql", driver)
	m.Log = newMigrateLogger()
	if err != nil {
		return err
	}
	v, _, _ := m.Version()
	if err := m.Force(int(v)); err != nil {
		return err
	}
	return nil
}

func initDiff() {
	rootCmd.AddCommand(diffCmd)
	diffCmd.Flags().BoolVar(&applyDiffToDB, "apply", false, "--apply")
}
