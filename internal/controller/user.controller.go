package controller

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kasvior-wallet-backend/internal/helper"
	"github.com/kasvior-wallet-backend/internal/service"
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
	claims, ok := helper.CheckClaims(ctx)
	if !ok {
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
	claims, ok := helper.CheckClaims(ctx)
	if !ok {
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

func (uc *UserController) GetTransactionReport(ctx *gin.Context) {
	duration := ctx.DefaultQuery("duration", "7d")
	if duration != "7d" {
		helper.JSONBadRequest(ctx)
		return
	}

	reportType := ctx.DefaultQuery("type", "all")
	if reportType != "all" && reportType != "income" && reportType != "expense" {
		helper.JSONBadRequest(ctx)
		return
	}

	claims, ok := helper.CheckClaims(ctx)
	if !ok {
		return
	}

	res, err := uc.userService.GetTransactionReport(ctx.Request.Context(), claims.UserId, reportType)
	if err != nil {
		log.Println("Error: ", err.Error())
		helper.JSONInternalServerError(ctx)
		return
	}

	helper.JSONSuccess(ctx, res, "Get Transaction Report Successfully")
}
