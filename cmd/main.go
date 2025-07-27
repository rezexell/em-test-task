package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/rezexell/em-test-task/internal/config"
	"log"
	"os"
)

func main() {
	cfg := config.InitConfig()

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUSER, cfg.DBPASSWORD, cfg.DBHOST, cfg.DBPORT, cfg.DBNAME)

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(context.Background())

	var greeting string
	err = conn.QueryRow(context.Background(), "select 'Hello, world!'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	log.Println(greeting)
}
