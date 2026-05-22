package controller

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/kasvior-wallet-backend/internal/binder"
	"github.com/kasvior-wallet-backend/internal/dto"
	"github.com/kasvior-wallet-backend/internal/jwttoken"
	"github.com/kasvior-wallet-backend/internal/response"
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
	claims, ok := jwttoken.CheckClaims(ctx)
	if !ok {
		return
	}

	res, err := uc.userService.GetProfile(ctx.Request.Context(), claims.UserId)
	if err != nil {
		log.Println("Error: ", err.Error())
		response.JSONInternalServerError(ctx)
		return
	}

	response.JSONSuccess(ctx, res, "Get Profile Successfully")
}

func (uc *UserController) UpdateProfile(ctx *gin.Context) {
	claims, ok := jwttoken.CheckClaims(ctx)
	if !ok {
		return
	}

	var body dto.UserUpdateProfileRequest
	if !binder.BindFormat(ctx, &body, binding.JSON) {
		return
	}

	res, err := uc.userService.UpdateProfile(ctx.Request.Context(), claims.UserId, body)
	if err != nil {
		log.Println("Error: ", err.Error())
		response.JSONInternalServerError(ctx)
		return
	}

	response.JSONSuccess(ctx, res, "Update Profile Successfully")
}

func (uc *UserController) UpdatePassword(ctx *gin.Context) {
	claims, ok := jwttoken.CheckClaims(ctx)
	if !ok {
		return
	}

	var body dto.UserUpdatePasswordRequest
	if !binder.BindFormat(ctx, &body, binding.JSON) {
		return
	}

	if err := uc.userService.UpdatePassword(ctx.Request.Context(), claims.UserId, body); err != nil {
		if errors.Is(err, service.ErrInvalidPassword) {
			response.JSONUnauthorized(ctx, "Invalid current password")
			return
		}

		log.Println("Error: ", err.Error())
		response.JSONInternalServerError(ctx)
		return
	}

	response.JSONSuccess(ctx, nil, "Update Password Successfully")
}

func (uc *UserController) UpdatePin(ctx *gin.Context) {
	claims, ok := jwttoken.CheckClaims(ctx)
	if !ok {
		return
	}

	var body dto.UserUpdatePinRequest
	if !binder.BindFormat(ctx, &body, binding.JSON) {
		return
	}

	if err := uc.userService.UpdatePin(ctx.Request.Context(), claims.UserId, body); err != nil {
		log.Println("Error: ", err.Error())
		response.JSONInternalServerError(ctx)
		return
	}

	response.JSONSuccess(ctx, nil, "Update PIN Successfully")
}

func (uc *UserController) CheckPin(ctx *gin.Context) {
	claims, ok := jwttoken.CheckClaims(ctx)
	if !ok {
		return
	}

	var body dto.UserCheckPinRequest
	if !binder.BindFormat(ctx, &body, binding.JSON) {
		return
	}

	res, err := uc.userService.CheckPin(ctx.Request.Context(), claims.UserId, body.Pin)
	if err != nil {
		if errors.Is(err, service.ErrInvalidPin) {
			log.Println("Invalid PIN: ", err.Error())
			response.JSONUnauthorized(ctx, "Invalid PIN")
			return
		}
		if errors.Is(err, service.ErrPinNotSet) {
			log.Println("PIN not set: ", err.Error())
			response.JSONUnauthorized(ctx, "PIN not set")
			return
		}

		log.Println("Error: ", err.Error())
		response.JSONInternalServerError(ctx)
		return
	}

	response.JSONSuccess(ctx, res, "PIN Valid")
}

func (uc *UserController) GetDashboardInformation(ctx *gin.Context) {
	claims, ok := jwttoken.CheckClaims(ctx)
	if !ok {
		return
	}

	res, err := uc.userService.GetDashboardInformation(ctx.Request.Context(), claims.UserId)
	if err != nil {
		log.Println("Error: ", err.Error())
		response.JSONInternalServerError(ctx)
		return
	}

	response.JSONSuccess(ctx, res, "Get Dashboard Information Successfully")
}

func (uc *UserController) GetTransactionReport(ctx *gin.Context) {
	duration := ctx.DefaultQuery("duration", "7d")
	if duration != "7d" {
		response.JSONBadRequest(ctx)
		return
	}

	reportType := ctx.DefaultQuery("type", "all")
	if reportType != "all" && reportType != "income" && reportType != "expense" {
		response.JSONBadRequest(ctx)
		return
	}

	claims, ok := jwttoken.CheckClaims(ctx)
	if !ok {
		return
	}

	res, err := uc.userService.GetTransactionReport(ctx.Request.Context(), claims.UserId, reportType)
	if err != nil {
		log.Println("Error: ", err.Error())
		response.JSONInternalServerError(ctx)
		return
	}

	response.JSONSuccess(ctx, res, "Get Transaction Report Successfully")
}
