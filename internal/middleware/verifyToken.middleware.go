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

func VerifyToken(authCache *repository.AuthCacheRepository) gin.HandlerFunc {
	return func(ctx *gin.Context) { // closure function
		token, ok := jwttoken.VerifyClientToken(ctx)
		if !ok {
			response.JSONAbortUnauthorized(ctx)
			return
		}

		var claims pkg.Claims
		if err := claims.VerifyJWT(token); err != nil {
			log.Println("Error: ", err.Error())
			if errors.Is(err, jwt.ErrTokenInvalidIssuer) || errors.Is(err, jwt.ErrTokenExpired) {
				response.JSONAbortUnauthorized(ctx)
				return
			}

			response.JSONAbortUnauthorized(ctx)
			return
		}

		isActive, err := authCache.IsTokenActive(ctx.Request.Context(), hashToken(token), claims.UserId)
		if !jwttoken.HandleTokenIsActive(ctx, isActive, err) {
			response.JSONAbortUnauthorized(ctx)
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
