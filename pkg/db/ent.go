package db

import (
	"database/sql"
	"easycoding/internal/config"
	"easycoding/pkg/ent"
	"easycoding/pkg/errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	entsql "entgo.io/ent/dialect/sql"
)

func CreateDBClient(config *config.Config) (*ent.Client, error) {
	if config.Database.CreateDatabase {
		if err := createDatabase(config); err != nil {
			return nil, err
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
		return nil, err
	}
	drv := entsql.OpenDB("mysql", db)
	return ent.NewClient(ent.Driver(drv)), nil
}

func createDatabase(config *config.Config) error {
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
