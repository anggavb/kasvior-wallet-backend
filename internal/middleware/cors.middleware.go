package middleware

import (
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware(ctx *gin.Context) {
	allowedOrigin := []string{
		"http://localhost:5173",
		"http://127.0.0.1:5173",
		"http://localhost:5174",
		"http://127.0.0.1:5174",
	}
	currentOrigin := ctx.GetHeader("Origin")
	if slices.Contains(allowedOrigin, currentOrigin) {
		ctx.Header("Access-Control-Allow-Origin", currentOrigin)
		ctx.Header("Access-Control-Allow-Credentials", "true")
	}

	allowedHeaders := []string{"Content-Type", "Authorization", "X-Reset-Token"}
	ctx.Header("Access-Control-Allow-Headers", strings.Join(allowedHeaders, ", "))

	allowedMethods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch, http.MethodOptions}
	ctx.Header("Access-Control-Allow-Methods", strings.Join(allowedMethods, ", "))

	if ctx.Request.Method == http.MethodOptions {
		ctx.AbortWithStatus(http.StatusNoContent)
		return
	}

	ctx.Next()
}
