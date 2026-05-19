package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kasvior-wallet-backend/internal/dto"
	"github.com/kasvior-wallet-backend/internal/middleware"
)

func InitRouter(router *gin.Engine, db *pgxpool.Pool) {
	// middleware global
	router.Use(middleware.CORSMiddleware)

	AuthRouter(router, db)
	UserRouter(router, db)
	TransactionRouter(router, db)

	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, dto.Response{
			Message: "Invalid Route",
			Success: false,
			Error:   "Not Found",
		})
	})
}
