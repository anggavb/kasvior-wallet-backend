package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kasvior-wallet-backend/internal/controller"
	"github.com/kasvior-wallet-backend/internal/middleware"
	"github.com/kasvior-wallet-backend/internal/repository"
	"github.com/kasvior-wallet-backend/internal/service"
)

func TransactionRouter(router *gin.Engine, db *pgxpool.Pool) {
	authRepo := repository.NewAuthRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	transactionService := service.NewTransactionService(transactionRepo)
	transactionController := controller.NewTransactionController(transactionService)

	transactionRouter := router.Group("/transaction", middleware.VerifyToken(authRepo))

	transactionRouter.GET("/transfer/receivers", transactionController.FindReceivers)
}
