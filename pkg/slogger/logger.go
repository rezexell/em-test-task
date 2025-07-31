package slogger

import (
	"log/slog"
	"os"
)

func InitLogger() *slog.Logger {
	var logger *slog.Logger
	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	return logger
}
