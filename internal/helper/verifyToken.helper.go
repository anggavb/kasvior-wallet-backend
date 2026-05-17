package helper

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

func VerifyClientToken(ctx *gin.Context) (string, bool) {
	bearerToken := ctx.GetHeader("Authorization")
	if bearerToken == "" {
		log.Println("Error: Authorization header is missing")
		JSONAbortUnauthorized(ctx, "Unauthorized, please login!")
		return "", false
	}

	splittedBearer := strings.Split(bearerToken, " ")
	if len(splittedBearer) != 2 {
		log.Println("Error: Invalid Authorization header format")
		JSONAbortUnauthorized(ctx, "Unauthorized, please login!")
		return "", false
	}

	return splittedBearer[1], true
}

func HandleTokenIsActive(ctx *gin.Context, isActive bool, err error) bool {
	if err != nil {
		log.Println("Error: ", err.Error())
		JSONAbortUnauthorized(ctx, "Unauthorized, please login!")
		return false
	}

	if !isActive {
		log.Println("Error: Token is not active")
		JSONAbortUnauthorized(ctx, "Unauthorized, please login!")
		return false
	}

	return true
}
