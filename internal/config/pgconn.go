package config

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB() (*pgxpool.Pool, error) {
	log.Println(os.Getenv("DB_URL"))
	return pgxpool.New(context.Background(), os.Getenv("DB_URL"))
}
