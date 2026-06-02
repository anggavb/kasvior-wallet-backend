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

func TransactionRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	authCache := repository.NewAuthCacheRepository(rdb)
	transactionRepo := repository.NewTransactionRepository(db)
	transactionService := service.NewTransactionService(transactionRepo, db)
	transactionController := controller.NewTransactionController(transactionService)

	transactionRouter := router.Group("/transaction", middleware.VerifyToken(authCache))

	{ // use for scoping route
		transactionRouter.GET("/history", transactionController.FindHistory)
		transactionRouter.GET("/payment-methods", transactionController.FindPaymentMethods)
		transactionRouter.POST("/topup", transactionController.CreateTopup)
	}

	{ // use for scoping route
		transferRouter := transactionRouter.Group("/transfer")
		transferRouter.POST("", transactionController.CreateTransfer)
		transferRouter.GET("/receivers", transactionController.FindReceivers)
	}
}
