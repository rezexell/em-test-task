package main

import (
	"github.com/rezexell/em-test-task/internal/config"
	"github.com/rezexell/em-test-task/internal/database"
	"github.com/rezexell/em-test-task/internal/handler"
	"github.com/rezexell/em-test-task/internal/repository"
	"github.com/rezexell/em-test-task/internal/service"
	"log"
)

func main() {
	cfg := config.InitConfig()

	database.ApplyMigrations(cfg)

	pool := database.InitDB(cfg)
	defer pool.Close()

	repos := repository.NewRepository(pool)
	services := service.NewService(repos)
	h := handler.NewHandler(services)

	server := h.InitRouter()
	if err := server.Run(":3000"); err != nil {
		log.Fatal(err)
	}
}
