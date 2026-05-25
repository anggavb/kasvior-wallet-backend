package jwttoken

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kasvior-wallet-backend/internal/response"
	"github.com/kasvior-wallet-backend/pkg"
)

func GetClaims(ctx *gin.Context) (pkg.Claims, bool) {
	claimsValue, ok := ctx.Get("claims")
	if !ok {
		log.Println("Error: Claims not found in context")
		response.JSONUnauthorized(ctx, "Unauthorized, please login!")
		return pkg.Claims{}, false
	}

	claims, ok := claimsValue.(pkg.Claims)
	if !ok {
		log.Println("Error: Invalid claims type")
		response.JSONUnauthorized(ctx, "Unauthorized, please login!")
		return pkg.Claims{}, false
	}

	return claims, true
}

func CheckAuthToken(ctx *gin.Context) (string, bool) {
	token, ok := ctx.Get("token")
	if !ok {
		log.Println("Error: token not found in context")
		response.JSONUnauthorized(ctx, "Unauthorized, please login!")
		return "", false
	}

	tokenString, ok := token.(string)
	if !ok {
		log.Println("Error: token type assertion failed")
		response.JSONUnauthorized(ctx, "Unauthorized, please login!")
		return "", false
	}
	return tokenString, true
}

func CheckExpiredToken(ctx *gin.Context) (*jwt.NumericDate, error) {
	claims, ok := GetClaims(ctx)
	if !ok {
		return nil, jwt.ErrTokenInvalidClaims
	}

	expiresAt, err := claims.GetExpirationTime()
	if err != nil || expiresAt == nil {
		if err != nil {
			log.Println("Error: ", err.Error())
		}
		log.Println("Error: expiresAt is nil")
		response.JSONUnauthorized(ctx, "Unauthorized, please login!")
		return nil, jwt.ErrTokenInvalidClaims
	}

	return expiresAt, nil
}
