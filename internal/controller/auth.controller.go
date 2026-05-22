package controller

import (
	"errors"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/kasvior-wallet-backend/internal/binder"
	"github.com/kasvior-wallet-backend/internal/dto"
	"github.com/kasvior-wallet-backend/internal/jwttoken"
	"github.com/kasvior-wallet-backend/internal/response"
	"github.com/kasvior-wallet-backend/internal/service"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (ac *AuthController) Register(ctx *gin.Context) {
	var body dto.AuthRequest
	if !binder.BindFormat(ctx, &body, binding.JSON) {
		return
	}

	res, err := ac.authService.RegisterUser(ctx.Request.Context(), body)
	if err != nil {
		log.Println("Error: ", err.Error())
		if strings.Contains(err.Error(), "users_email_key") {
			response.JSONDuplicate(ctx, "Email Already Used")
			return
		}
		response.JSONInternalServerError(ctx)
		return
	}

	response.JSONCreated(ctx, res, "Register Successfully")
}

func (ac *AuthController) Login(ctx *gin.Context) {
	var body dto.AuthRequest
	if !binder.BindFormat(ctx, &body, binding.JSON) {
		return
	}

	res, err := ac.authService.LoginUser(ctx.Request.Context(), body)
	if err != nil {
		log.Println("Error: ", err.Error())
		if strings.Contains(err.Error(), "wrong password") || strings.Contains(err.Error(), "no rows") {
			response.JSONUnauthorized(ctx, "Invalid email or password")
			return
		}
		response.JSONInternalServerError(ctx)
		return
	}

	response.JSONSuccess(ctx, res, "Login Successfully")
}

func (ac *AuthController) Logout(ctx *gin.Context) {
	tokenString, ok := jwttoken.CheckAuthToken(ctx)
	if !ok {
		return
	}

	expiresAt, err := jwttoken.CheckExpiredToken(ctx)
	if err != nil {
		return
	}

	if err := ac.authService.LogoutUser(ctx.Request.Context(), tokenString, &expiresAt.Time); err != nil {
		log.Println("Error: ", err.Error())
		if errors.Is(err, service.ErrTokenAlreadyExpired) {
			response.JSONUnauthorized(ctx, "Token already expired")
			return
		}
		response.JSONInternalServerError(ctx)
		return
	}

	response.JSONSuccess(ctx, nil, "Logout Successfully")
}
