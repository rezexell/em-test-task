package main

import (
	"github.com/rezexell/em-test-task/internal/config"
	"github.com/rezexell/em-test-task/internal/handler"
	"github.com/rezexell/em-test-task/internal/repository"
	"github.com/rezexell/em-test-task/internal/service"
	"github.com/rezexell/em-test-task/pkg/postgres"
	"log"
)

func main() {
	cfg := config.InitConfig()

	postgres.ApplyMigrations(cfg)

	pool := postgres.InitDB(cfg)
	defer pool.Close()

	repos := repository.NewRepository(pool)
	services := service.NewService(repos)
	h := handler.NewHandler(services)

	server := h.InitRouter()
	if err := server.Run(":3000"); err != nil {
		log.Fatal(err)
	}
}
