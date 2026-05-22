package config

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB() (*pgxpool.Pool, error) {
	return pgxpool.New(context.Background(), os.Getenv("DB_URL"))
}
