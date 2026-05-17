package helper

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kasvior-wallet-backend/internal/dto"
)

// Status 401 - Abort when Unauthorized
func JSONAbortUnauthorized(ctx *gin.Context, message string) {
	ctx.AbortWithStatusJSON(http.StatusUnauthorized, dto.Response{
		Message: message,
		Error:   "Unauthorized",
	})
}
