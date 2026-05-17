package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kasvior-wallet-backend/internal/helper"
	"github.com/kasvior-wallet-backend/internal/repository"
	"github.com/kasvior-wallet-backend/pkg"
)

func VerifyToken(authRepository *repository.AuthRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
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

		isActive, err := authRepository.IsTokenActive(ctx.Request.Context(), hashToken(token))
		if err != nil {
			log.Println("Error: ", err.Error())
			helper.JSONAbortUnauthorized(ctx, "Unauthorized, please login!")
			return
		}

		if !isActive {
			helper.JSONAbortUnauthorized(ctx, "Unauthorized, please login!")
			return
		}

		ctx.Set("claims", claims)
		ctx.Set("token", token)
		ctx.Next()
	}
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
