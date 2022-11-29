package db

import (
	"context"
	"database/sql"
	"easycoding/internal/config"
	"easycoding/pkg/ent"
	"easycoding/pkg/errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	entsql "entgo.io/ent/dialect/sql"
)

func CreateDBClient(config *config.Config, tracer trace.Tracer) (*ent.Client, error) {
	if config.Database.CreateDatabase {
		if err := CreateDatabase(config); err != nil {
			return nil, err
		}
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8",
		config.Database.User,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.DBName,
	)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	drv := entsql.OpenDB("mysql", db)

	client := ent.NewClient(ent.Driver(drv))
	client.Use(func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			opType := m.Op().String()
			tablesName := m.Type()
			spanName := fmt.Sprintf("%s-%s", tablesName, opType)
			newCtx, span := tracer.Start(ctx, spanName)
			span.SetAttributes(
				attribute.String("table", tablesName),
				attribute.String("operation", opType),
			)
			defer span.End()
			return next.Mutate(newCtx, m)
		})
	})
	return client, nil
}

func CreateDatabase(config *config.Config) error {
	db, err := sql.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/",
		config.Database.User,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
	))
	if err != nil {
		return errors.ErrInvalidf(err, "fail to open connection with mysql")
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf(
		"CREATE DATABASE IF NOT EXISTS %s ", config.Database.DBName))
	if err != nil {
		return errors.ErrInvalid(err)
	}
	return nil
}
