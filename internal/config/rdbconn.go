package config

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("RDB_ADDR"),
		Username: os.Getenv("RDB_USER"),
		Password: os.Getenv("RDB_PASS"),
	})

	cmdErr := rdb.Ping(context.Background())

	return rdb, cmdErr.Err()
}
