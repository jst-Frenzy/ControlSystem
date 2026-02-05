package dataBase

import (
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logrus.WithError(err).Info("Migration failed")
	}
}
