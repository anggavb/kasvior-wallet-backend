package jwttoken

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kasvior-wallet-backend/internal/response"
)

func VerifyClientToken(ctx *gin.Context) (string, bool) {
	bearerToken := ctx.GetHeader("Authorization")
	if bearerToken == "" {
		log.Println("Error: Authorization header is missing")
		response.JSONUnauthorized(ctx, "Unauthorized, please login!")
		return "", false
	}

	splittedBearer := strings.Fields(bearerToken)
	if len(splittedBearer) == 1 {
		log.Println("Error: Invalid Authorization header format")
		response.JSONUnauthorized(ctx, "Unauthorized, please login!")
		return "", false
	}

	if len(splittedBearer) != 2 || !strings.EqualFold(splittedBearer[0], "Bearer") {
		log.Println("Error: Invalid Authorization header format")
		response.JSONUnauthorized(ctx, "Unauthorized, please login!")
		return "", false
	}

	return splittedBearer[1], true
}

func HandleTokenIsActive(ctx *gin.Context, isActive bool, err error) bool {
	if err != nil {
		log.Println("Error: ", err.Error())
		response.JSONUnauthorized(ctx, "Unauthorized, please login!")
		return false
	}

	if !isActive {
		log.Println("Error: Token is not active")
		response.JSONUnauthorized(ctx, "Unauthorized, please login!")
		return false
	}

	return true
}
