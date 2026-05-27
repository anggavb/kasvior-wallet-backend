package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/kasvior-wallet-backend/docs"
	"github.com/kasvior-wallet-backend/internal/dto"
	"github.com/kasvior-wallet-backend/internal/middleware"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	router.Use(middleware.CORSMiddleware)

	router.Static("/image", "public/img")
	router.Static("/logo", "public/payment")

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	AuthRouter(router, db, rdb)
	UserRouter(router, db, rdb)
	TransactionRouter(router, db, rdb)

	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, dto.Response{
			Message: "Invalid Route",
			Success: false,
			Error:   "Not Found",
		})
	})
}
