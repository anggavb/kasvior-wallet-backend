package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kasvior-wallet-backend/internal/dto"
)

func JSONAbortUnauthorized(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
		Message: "Unauthorized",
		Success: false,
	})
}

func JSONAbortInternalServerError(ctx *gin.Context) {
	ctx.AbortWithStatusJSON(http.StatusInternalServerError, dto.Response{
		Message: "Internal Server Error",
		Success: false,
	})
}
