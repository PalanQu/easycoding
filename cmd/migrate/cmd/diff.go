package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"easycoding/pkg/ent"

	"entgo.io/ent/dialect/sql/schema"
	"github.com/spf13/cobra"
)

var applyDiff bool

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
		if err := client.Schema.WriteTo(
			context.Background(),
			os.Stdout,
			migrateOptions...,
		); err != nil {
			log.Fatalf("failed printing schema changes: %v", err)
		}
		if applyDiff {
			if err := client.Schema.Create(
				context.Background(),
				migrateOptions...,
			); err != nil {
				log.Fatalf("failed creating schema resources: %v", err)
				return err
			}
		}
		log.Println("apply success")
		return nil
	},
}

func initDiff() {
	rootCmd.AddCommand(diffCmd)
	diffCmd.Flags().BoolVar(&applyDiff, "apply", false, "--apply")
}
