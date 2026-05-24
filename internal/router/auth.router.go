package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kasvior-wallet-backend/internal/controller"
	"github.com/kasvior-wallet-backend/internal/middleware"
	"github.com/kasvior-wallet-backend/internal/repository"
	"github.com/kasvior-wallet-backend/internal/service"
	"github.com/kasvior-wallet-backend/pkg"
)

func AuthRouter(router *gin.Engine, db *pgxpool.Pool) {
	authRouter := router.Group("/auth")

	authRepo := repository.NewAuthRepository(db)
	smtpMailer := pkg.NewSMTPMailerFromEnv()
	authService := service.NewAuthService(authRepo, smtpMailer)
	authController := controller.NewAuthController(authService)

	authRouter.POST("", authController.Login)
	authRouter.POST("/register", authController.Register)
	authRouter.POST("/forgot-password", authController.ForgotPassword)
	authRouter.POST("/reset-password", authController.ResetPassword)
	authRouter.DELETE("/logout", middleware.VerifyToken(authRepo), authController.Logout)
}
