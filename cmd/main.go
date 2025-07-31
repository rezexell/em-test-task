package main

import (
	"github.com/rezexell/em-test-task/internal/config"
	"github.com/rezexell/em-test-task/internal/handler"
	"github.com/rezexell/em-test-task/internal/repository"
	"github.com/rezexell/em-test-task/internal/service"
	"github.com/rezexell/em-test-task/pkg/postgres"
	"github.com/rezexell/em-test-task/pkg/slogger"
	"os"
)

func main() {
	cfg := config.InitConfig()
	logger := slogger.InitLogger(cfg)
	logger.Info("Logger initialized")

	postgres.ApplyMigrations(cfg, logger)

	db := postgres.InitDB(cfg, logger)

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	h := handler.NewHandler(services, logger)

	server := h.InitRouter()
	if err := server.Run(":3000"); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
