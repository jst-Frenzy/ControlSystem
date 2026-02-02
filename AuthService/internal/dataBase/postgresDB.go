package dataBase

import (
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var PostgresDB *gorm.DB

func InitPostgres() {
	dsn := os.Getenv("AUTH_DB_DSN")
	var err error
	if PostgresDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{}); err != nil {
		logrus.WithError(err).Fatal("Failed to connect to database")
	}
}
