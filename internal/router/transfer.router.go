package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kasvior-wallet-backend/internal/controller"
	"github.com/kasvior-wallet-backend/internal/middleware"
	"github.com/kasvior-wallet-backend/internal/repository"
	"github.com/kasvior-wallet-backend/internal/service"
)

func TransferRouter(router *gin.Engine, db *pgxpool.Pool) {
	authRepo := repository.NewAuthRepository(db)
	transferRepo := repository.NewTransferRepository(db)
	transferService := service.NewTransferService(transferRepo)
	transferController := controller.NewTransferController(transferService)

	transferRouter := router.Group("/transfers", middleware.VerifyToken(authRepo))

	transferRouter.GET("/receivers", transferController.FindReceivers)
}
