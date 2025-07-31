package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

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

func InitDB(cfg *config.Config, logger *slog.Logger) *gorm.DB {
	dsn := getConnString(cfg)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error("Unable to create GORM connection", slog.Any("err", err.Error()))
		os.Exit(1)
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Error("Failed to get underlying DB:", slog.Any("err", err.Error()))
		os.Exit(1)
	}

	if err := sqlDB.PingContext(context.Background()); err != nil {
		logger.Error("Unable to ping database:", slog.Any("err", err.Error()))
		os.Exit(1)
	}

	logger.Info("Database connection successfully established")
	return db
}

func ApplyMigrations(cfg *config.Config, logger *slog.Logger) {

	m, err := migrate.New(
		"file:///app/migrations",
		getConnString(cfg))
	if err != nil {
		logger.Error("Failed to initialize migrate:", slog.Any("err", err.Error()))
		os.Exit(1)
	}
	defer func(m *migrate.Migrate) {
		_, _ = m.Close()
	}(m)

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.Error("Failed to apply migrations:", slog.Any("err", err.Error()))
		os.Exit(1)
	}

	version, dirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		logger.Warn("Failed to get migration version", slog.Any("err", err.Error()))
	} else {
		logger.Info("Database migrations applied", slog.Int("version", int(version)), slog.Bool("dirty", dirty))
	}
}
