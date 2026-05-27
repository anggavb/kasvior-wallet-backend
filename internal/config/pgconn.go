package config

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB() (*pgxpool.Pool, error) {
	pg, _ := pgxpool.New(context.Background(), os.Getenv("DB_URL"))
	return pg, pg.Ping(context.Background())
}
