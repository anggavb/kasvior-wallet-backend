package middleware

import (
	"errors"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kasvior-wallet-backend/internal/helper"
	"github.com/kasvior-wallet-backend/pkg"
)

func VerifyToken(ctx *gin.Context) {
	bearerToken := ctx.GetHeader("Authorization")
	if bearerToken == "" {
		helper.JSONAbortUnauthorized(ctx, "Unauthorized, please login!")
		return
	}

	splittedBearer := strings.Split(bearerToken, " ")
	if len(splittedBearer) != 2 {
		helper.JSONAbortUnauthorized(ctx, "Unauthorized, please login!")
		return
	}

	token := splittedBearer[1]

	var claims pkg.Claims
	if err := claims.VerifyJWT(token); err != nil {
		log.Println("Error: ", err.Error())
		if errors.Is(err, jwt.ErrTokenInvalidIssuer) || errors.Is(err, jwt.ErrTokenExpired) {
			helper.JSONAbortUnauthorized(ctx, err.Error())
			return
		}

		helper.JSONAbortUnauthorized(ctx, "Unauthorized, please login!")
		return
	}

	ctx.Set("claims", claims)
	ctx.Next()
}
