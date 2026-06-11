package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kasvior-wallet-backend/internal/controller"
	"github.com/kasvior-wallet-backend/internal/middleware"
	"github.com/kasvior-wallet-backend/internal/repository"
	"github.com/kasvior-wallet-backend/internal/service"
	"github.com/redis/go-redis/v9"
)

func UserRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	authCache := repository.NewAuthCacheRepository(rdb)
	userRepo := repository.NewUserRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	userService := service.NewUserService(userRepo, transactionRepo, authCache, db)
	userController := controller.NewUserController(userService)

	userRouter := router.Group("/users", middleware.VerifyToken(authCache))

	{
		meRouter := userRouter.Group("/me")

		meRouter.GET("", userController.GetProfile)
		meRouter.PATCH("", userController.UpdateProfile)
		meRouter.PATCH("/password", userController.UpdatePassword)
		meRouter.PATCH("/pin", userController.UpdatePin)
		meRouter.POST("/pin/check", userController.CheckPin)
		meRouter.GET("/wallet", userController.GetDashboardInformation)
		meRouter.GET("/transaction-report", userController.GetTransactionReport)
	}
}
