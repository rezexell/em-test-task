package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rezexell/em-test-task/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getConnString(cfg *config.Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUSER, cfg.DBPASSWORD, cfg.DBHOST, cfg.DBPORT, cfg.DBNAME)
}

func InitDB(cfg *config.Config) *gorm.DB {
	dsn := getConnString(cfg)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Unable to create GORM connection: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get underlying DB: %v", err)
	}

	if err := sqlDB.PingContext(context.Background()); err != nil {
		log.Fatalf("Unable to ping database: %v", err)
	}

	log.Println("Database connection established")
	return db
}

func ApplyMigrations(cfg *config.Config) {

	m, err := migrate.New(
		"file:///app/migrations",
		getConnString(cfg))
	if err != nil {
		log.Fatalf("Failed to initialize migrate: %v", err)
	}
	defer func(m *migrate.Migrate) {
		_, _ = m.Close()
	}(m)

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	version, dirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		log.Printf("Failed to get migration version: %v", err)
	} else {
		log.Printf("Database migrations applied. Version: %d, Dirty: %v", version, dirty)
	}
}
