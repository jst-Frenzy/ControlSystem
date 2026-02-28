package dataBase

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var PostgresDB *gorm.DB

func InitPostgres() {
	dsn := os.Getenv("ORDER_DB_DSN")
	var err error
	if PostgresDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{}); err != nil {
		logrus.WithError(err).Fatal("Failed to connect to database")
	}
	runMigrations(os.Getenv("MIGRATE_DB_DSN"))
}

func runMigrations(dsn string) {
	m, errCreate := migrate.New("file:///app/migrations", dsn)
	if errCreate != nil {
		logrus.WithError(errCreate).Info("cant initialize migrations")
		return
	}
	defer func() {
		if _, err := m.Close(); err != nil {
			logrus.WithError(errCreate).Info("cant close migrations")
		}
	}()

	err := m.Up()
	if err != nil {
		logrus.WithError(err).Info("Migration failed")
	}
}
