package controller

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kasvior-wallet-backend/internal/helper"
	"github.com/kasvior-wallet-backend/internal/service"
	"github.com/kasvior-wallet-backend/pkg"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (uc *UserController) GetProfile(ctx *gin.Context) {
	claimsValue, ok := ctx.Get("claims")
	if !ok {
		helper.JSONUnauthorized(ctx, "Unauthorized, please login!")
		return
	}

	claims, ok := claimsValue.(pkg.Claims)
	if !ok {
		helper.JSONUnauthorized(ctx, "Unauthorized, please login!")
		return
	}

	res, err := uc.userService.GetProfile(ctx.Request.Context(), claims.UserId)
	if err != nil {
		log.Println("Error: ", err.Error())
		helper.JSONInternalServerError(ctx)
		return
	}

	helper.JSONSuccess(ctx, res, "Get Profile Successfully")
}

func (uc *UserController) GetDashboardInformation(ctx *gin.Context) {
	claimsValue, ok := ctx.Get("claims")
	if !ok {
		helper.JSONUnauthorized(ctx, "Unauthorized, please login!")
		return
	}

	claims, ok := claimsValue.(pkg.Claims)
	if !ok {
		helper.JSONUnauthorized(ctx, "Unauthorized, please login!")
		return
	}

	res, err := uc.userService.GetDashboardInformation(ctx.Request.Context(), claims.UserId)
	if err != nil {
		log.Println("Error: ", err.Error())
		helper.JSONInternalServerError(ctx)
		return
	}

	helper.JSONSuccess(ctx, res, "Get Dashboard Information Successfully")
}
