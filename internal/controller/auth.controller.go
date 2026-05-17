package controller

import (
	"errors"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/kasvior-wallet-backend/internal/dto"
	"github.com/kasvior-wallet-backend/internal/helper"
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
	if !helper.BindFormat(ctx, &body, binding.JSON) {
		return
	}

	res, err := ac.authService.RegisterUser(ctx.Request.Context(), body)
	if err != nil {
		log.Println("Error: ", err.Error())
		if strings.Contains(err.Error(), "users_email_key") {
			helper.JSONDuplicate(ctx, "Email Already Used")
			return
		}
		helper.JSONInternalServerError(ctx)
		return
	}

	helper.JSONCreated(ctx, res, "Register Successfully")
}

func (ac *AuthController) Login(ctx *gin.Context) {
	var body dto.AuthRequest
	if !helper.BindFormat(ctx, &body, binding.JSON) {
		return
	}

	res, err := ac.authService.LoginUser(ctx.Request.Context(), body)
	if err != nil {
		log.Println("Error: ", err.Error())
		if strings.Contains(err.Error(), "wrong password") || strings.Contains(err.Error(), "no rows") {
			helper.JSONUnauthorized(ctx, "Invalid email or password")
			return
		}
		helper.JSONInternalServerError(ctx)
		return
	}

	helper.JSONSuccess(ctx, res, "Login Successfully")
}

func (ac *AuthController) Logout(ctx *gin.Context) {
	tokenString, ok := helper.CheckAuthToken(ctx)
	if !ok {
		return
	}

	expiresAt, err := helper.CheckExpiredToken(ctx)
	if err != nil {
		return
	}

	if err := ac.authService.LogoutUser(ctx.Request.Context(), tokenString, &expiresAt.Time); err != nil {
		log.Println("Error: ", err.Error())
		if errors.Is(err, service.ErrTokenAlreadyExpired) {
			helper.JSONUnauthorized(ctx, "Token already expired")
			return
		}
		helper.JSONInternalServerError(ctx)
		return
	}

	helper.JSONSuccess(ctx, nil, "Logout Successfully")
}
