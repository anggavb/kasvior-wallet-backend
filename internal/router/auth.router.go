package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kasvior-wallet-backend/internal/controller"
	"github.com/kasvior-wallet-backend/internal/middleware"
	"github.com/kasvior-wallet-backend/internal/repository"
	"github.com/kasvior-wallet-backend/internal/service"
	"github.com/kasvior-wallet-backend/pkg"
	"github.com/redis/go-redis/v9"
)

func AuthRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	authRouter := router.Group("/auth")

	authRepo := repository.NewAuthRepository(db)
	authCache := repository.NewAuthCacheRepository(rdb)
	smtpMailer := pkg.NewSMTPMailerFromEnv()
	authService := service.NewAuthService(authRepo, authCache, smtpMailer)
	authController := controller.NewAuthController(authService)

	authRouter.POST("", authController.Login)
	authRouter.POST("/register", authController.Register)
	authRouter.POST("/forgot-password", authController.ForgotPassword)
	authRouter.POST("/reset-password", authController.ResetPassword)
	authRouter.DELETE("/logout", middleware.VerifyToken(authCache), authController.Logout)
}
