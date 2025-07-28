package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rezexell/em-test-task/internal/config"
	"github.com/rezexell/em-test-task/internal/handler"
	"github.com/rezexell/em-test-task/internal/repository"
	"github.com/rezexell/em-test-task/internal/service"
	"os"
)

func main() {
	cfg := config.InitConfig()

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUSER, cfg.DBPASSWORD, cfg.DBHOST, cfg.DBPORT, cfg.DBNAME)

	dbpool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	repos := repository.NewRepository(dbpool)
	services := service.NewService(repos)
	h := handler.NewHandler(services)

	server := h.InitRouter()
	server.Run(":3000")
}
