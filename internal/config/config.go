package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	DBHOST     string
	DBPORT     string
	DBUSER     string
	DBPASSWORD string
	DBNAME     string
}

func InitConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("err loading: %v", err)
	}

	return &Config{
		DBHOST:     os.Getenv("DB_HOST"),
		DBPORT:     os.Getenv("DB_PORT"),
		DBUSER:     os.Getenv("DB_USER"),
		DBPASSWORD: os.Getenv("DB_PASSWORD"),
		DBNAME:     os.Getenv("DB_NAME"),
	}
}
