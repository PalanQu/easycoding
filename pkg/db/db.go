package db

import (
	"database/sql"
	"easycoding/internal/config"
	"easycoding/pkg/errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func CreateGdb(config *config.Config, logger *logrus.Logger) (*gorm.DB, error) {
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
	dial := mysql.Open(dsn)
	return gorm.Open(dial, &gorm.Config{
		Logger: newGormLogger(logger),
	})
}

func CreateTestingGdb(db *sql.DB) (*gorm.DB, error) {
	return gorm.Open(mysql.New(mysql.Config{
		Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
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
