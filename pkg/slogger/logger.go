package slogger

import (
	"github.com/rezexell/em-test-task/internal/config"
	"log/slog"
	"os"
)

func InitLogger(cfg *config.Config) *slog.Logger {
	var logger *slog.Logger
	switch cfg.LOGLEVEL {
	case "dev":
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	case "prod":
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelWarn,
		}))
	}
	return logger
}
