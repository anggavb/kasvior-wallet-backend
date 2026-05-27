package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/kasvior-wallet-backend/internal/binder"
	"github.com/kasvior-wallet-backend/internal/config"
	"github.com/kasvior-wallet-backend/internal/router"
)

// @title						Backend Kasvior Wallet API
// @version						1.0
// @description					API documentation for Kasvior Wallet backend application

// @license.name				MIT

// @host						localhost:8080
// @BasePath					/

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description					Bearer token used for authorization. Example: Bearer <token>
func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading env. \ncause: %s", err.Error())
	}

	app := gin.Default()

	// PostgreSQL Connect
	conn, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("DB connection error. \ncause: %s", err.Error())
	}
	defer conn.Close()
	log.Println("DB Connected")

	// Redis Connect
	rdb, err := config.ConnectRedis()
	if err != nil {
		log.Fatalf("Redis connection error. \ncause: %s", err.Error())
	}
	defer rdb.Close()
	log.Println("Redis Connected")

	// install router
	router.InitRouter(app, conn)

	addr := fmt.Sprintf("%s:%s", os.Getenv("APP_HOST"), os.Getenv("APP_PORT"))
	app.Run(addr)
}
