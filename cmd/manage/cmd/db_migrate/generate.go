package db

import (
	"context"
	"easycoding/common/workspace"
	"easycoding/pkg/db"
	"easycoding/pkg/ent/migrate"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"ariga.io/atlas/sql/sqltool"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql/schema"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate migration files",
	RunE: func(cmd *cobra.Command, args []string) error {
		if config.Database.CreateDatabase {
			if err := db.CreateDatabase(config); err != nil {
				return err
			}
		}
		migrationDir := filepath.Join(workspace.GetWorkspace(), "migrations")
		if _, err := os.Stat(migrationDir); errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir(migrationDir, os.ModePerm)
			if err != nil {
				return err
			}
		}
		ctx := context.Background()
		// Create a local migration directory able to understand golang-migrate
		// migration file format for replay.
		dir, err := sqltool.NewGolangMigrateDir(
			filepath.Join(workspace.GetWorkspace(), "migrations"))
		if err != nil {
			log.Fatalf("failed creating atlas migration directory: %v", err)
		}
		// Migrate diff options.
		opts := []schema.MigrateOption{
			schema.WithDir(dir),                         // provide migration directory
			schema.WithMigrationMode(schema.ModeReplay), // provide migration mode
			schema.WithDialect(dialect.MySQL),           // Ent dialect to use
			schema.WithFormatter(sqltool.GolangMigrateFormatter),
		}
		// Generate migrations using Atlas support for MySQL (note the Ent
		// dialect option passed above).
		dsn := fmt.Sprintf("mysql://%s:%s@%s:%s/%s",
			config.Database.User,
			config.Database.Password,
			config.Database.Host,
			config.Database.Port,
			config.Database.DBName,
		)
		err = migrate.Diff(ctx, dsn, opts...)
		if err != nil {
			log.Fatalf("failed generating migration file: %v", err)
		}
		return nil

	},
}

func initGenerate() {
	dbMigrateCmd.AddCommand(generateCmd)
}
