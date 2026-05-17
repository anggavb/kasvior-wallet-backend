package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kasvior-wallet-backend/internal/controller"
	"github.com/kasvior-wallet-backend/internal/middleware"
	"github.com/kasvior-wallet-backend/internal/repository"
	"github.com/kasvior-wallet-backend/internal/service"
)

func UserRouter(router *gin.Engine, db *pgxpool.Pool) {
	authRepo := repository.NewAuthRepository(db)
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	userRouter := router.Group("/users", middleware.VerifyToken(authRepo))

	userRouter.GET("/me", userController.GetProfile)
	userRouter.PATCH("/me", userController.UpdateProfile)
	userRouter.POST("/me/pin/check", userController.CheckPin)
	userRouter.GET("/me/wallet", userController.GetDashboardInformation)
	userRouter.GET("/me/transaction-report", userController.GetTransactionReport)
}
