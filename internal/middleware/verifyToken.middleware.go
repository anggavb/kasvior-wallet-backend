package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kasvior-wallet-backend/internal/jwttoken"
	"github.com/kasvior-wallet-backend/internal/repository"
	"github.com/kasvior-wallet-backend/internal/response"
	"github.com/kasvior-wallet-backend/pkg"
)

func VerifyToken(authRepository *repository.AuthRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) { // closure function
		token, ok := jwttoken.VerifyClientToken(ctx)
		if !ok {
			return
		}

		var claims pkg.Claims
		if err := claims.VerifyJWT(token); err != nil {
			log.Println("Error: ", err.Error())
			if errors.Is(err, jwt.ErrTokenInvalidIssuer) || errors.Is(err, jwt.ErrTokenExpired) {
				response.JSONUnauthorized(ctx, err.Error())
				return
			}

			response.JSONUnauthorized(ctx, "Unauthorized, please login!")
			return
		}

		isActive, err := authRepository.IsTokenActive(ctx.Request.Context(), hashToken(token))
		if !jwttoken.HandleTokenIsActive(ctx, isActive, err) {
			return
		}

		ctx.Set("claims", claims)
		ctx.Next()
	}
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
