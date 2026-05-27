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

	{
		meRouter := userRouter.Group("/me")

		meRouter.GET("/", userController.GetProfile)
		meRouter.PATCH("/", userController.UpdateProfile)
		meRouter.PATCH("/password", userController.UpdatePassword)
		meRouter.PATCH("/pin", userController.UpdatePin)
		meRouter.POST("/pin/check", userController.CheckPin)
		meRouter.GET("/wallet", userController.GetDashboardInformation)
		meRouter.GET("/transaction-report", userController.GetTransactionReport)
	}
}
