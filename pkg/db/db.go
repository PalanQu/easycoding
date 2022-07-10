package db

import (
	"database/sql"
	"easycoding/internal/config"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func CreateGdb(config *config.Config, logger *logrus.Logger) (*gorm.DB, error) {
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
