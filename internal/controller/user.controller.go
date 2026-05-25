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

// GetProfile godoc
// @Summary		Get current user profile
// @Description	Get profile information for the authenticated user.
// @Tags			Users
// @Produce		json
// @Security		ApiKeyAuth
// @Param			X-Swagger	header		string			false	"Set true when using a raw token from Swagger UI"
// @Success		200	{object}	dto.Response	"Get Profile Successfully"
// @Failure		401	{object}	dto.Response	"Unauthorized"
// @Failure		500	{object}	dto.Response	"Internal server error"
// @Router			/users/me/ [get]
func (uc *UserController) GetProfile(ctx *gin.Context) {
	claims, ok := jwttoken.GetClaims(ctx)
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

// UpdateProfile godoc
// @Summary		Update current user profile
// @Description	Update at least one profile field for the authenticated user. Profile fields are limited to 255 characters.
// @Tags			Users
// @Accept			json
// @Produce		json
// @Security		ApiKeyAuth
// @Param			X-Swagger	header		string							false	"Set true when using a raw token from Swagger UI"
// @Param			request	body		dto.UserUpdateProfileRequest	true	"Update profile request body"
// @Success		200		{object}	dto.Response					"Update Profile Successfully"
// @Failure		400		{object}	dto.Response					"Bad request"
// @Failure		401		{object}	dto.Response					"Unauthorized"
// @Failure		422		{object}	dto.Response					"Validation error"
// @Failure		500		{object}	dto.Response					"Internal server error"
// @Router			/users/me/ [patch]
func (uc *UserController) UpdateProfile(ctx *gin.Context) {
	claims, ok := jwttoken.GetClaims(ctx)
	if !ok {
		return
	}

	var body dto.UserUpdateProfileRequest
	if err := binder.BindFormat(ctx, &body, binding.JSON); err != nil {
		errorMessages := binder.FormatValidationError(err)
		if len(errorMessages) > 0 && errorMessages["error"] != "" {
			response.JSONBadRequest(ctx)
			return
		}
		response.JSONUnprocessableEntity(ctx, errorMessages)
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

// UpdatePassword godoc
// @Summary		Update current user password
// @Description	Update the authenticated user's password.
// @Tags			Users
// @Accept			json
// @Produce		json
// @Security		ApiKeyAuth
// @Param			X-Swagger	header		string								false	"Set true when using a raw token from Swagger UI"
// @Param			request	body		dto.UserUpdatePasswordRequest	true	"Update password request body"
// @Success		200		{object}	dto.Response						"Update Password Successfully"
// @Failure		400		{object}	dto.Response						"Bad request"
// @Failure		401		{object}	dto.Response						"Invalid current password or unauthorized"
// @Failure		422		{object}	dto.Response						"Validation error"
// @Failure		500		{object}	dto.Response						"Internal server error"
// @Router			/users/me/password [patch]
func (uc *UserController) UpdatePassword(ctx *gin.Context) {
	claims, ok := jwttoken.GetClaims(ctx)
	if !ok {
		return
	}

	var body dto.UserUpdatePasswordRequest
	if err := binder.BindFormat(ctx, &body, binding.JSON); err != nil {
		errorMessages := binder.FormatValidationError(err)
		if len(errorMessages) > 0 && errorMessages["error"] != "" {
			response.JSONBadRequest(ctx)
			return
		}
		response.JSONUnprocessableEntity(ctx, errorMessages)
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

// UpdatePin godoc
// @Summary		Update current user PIN
// @Description	Update the authenticated user's 6-digit numeric PIN.
// @Tags			Users
// @Accept			json
// @Produce		json
// @Security		ApiKeyAuth
// @Param			X-Swagger	header		string						false	"Set true when using a raw token from Swagger UI"
// @Param			request	body		dto.UserUpdatePinRequest	true	"Update PIN request body"
// @Success		200		{object}	dto.Response				"Update PIN Successfully"
// @Failure		400		{object}	dto.Response				"Bad request"
// @Failure		401		{object}	dto.Response				"Unauthorized"
// @Failure		422		{object}	dto.Response				"Validation error"
// @Failure		500		{object}	dto.Response				"Internal server error"
// @Router			/users/me/pin [patch]
func (uc *UserController) UpdatePin(ctx *gin.Context) {
	claims, ok := jwttoken.GetClaims(ctx)
	if !ok {
		return
	}

	var body dto.UserUpdatePinRequest
	if err := binder.BindFormat(ctx, &body, binding.JSON); err != nil {
		errorMessages := binder.FormatValidationError(err)
		if len(errorMessages) > 0 && errorMessages["error"] != "" {
			response.JSONBadRequest(ctx)
			return
		}
		response.JSONUnprocessableEntity(ctx, errorMessages)
		return
	}

	if err := uc.userService.UpdatePin(ctx.Request.Context(), claims.UserId, body); err != nil {
		log.Println("Error: ", err.Error())
		response.JSONInternalServerError(ctx)
		return
	}

	response.JSONSuccess(ctx, nil, "Update PIN Successfully")
}

// CheckPin godoc
// @Summary		Check current user PIN
// @Description	Validate the authenticated user's 6-digit numeric PIN.
// @Tags			Users
// @Accept			json
// @Produce		json
// @Security		ApiKeyAuth
// @Param			X-Swagger	header		string						false	"Set true when using a raw token from Swagger UI"
// @Param			request	body		dto.UserCheckPinRequest	true	"Check PIN request body"
// @Success		200		{object}	dto.Response				"PIN Valid"
// @Failure		400		{object}	dto.Response				"Bad request"
// @Failure		401		{object}	dto.Response				"Invalid PIN, PIN not set, or unauthorized"
// @Failure		422		{object}	dto.Response				"Validation error"
// @Failure		500		{object}	dto.Response				"Internal server error"
// @Router			/users/me/pin/check [get]
func (uc *UserController) CheckPin(ctx *gin.Context) {
	claims, ok := jwttoken.GetClaims(ctx)
	if !ok {
		return
	}

	var body dto.UserCheckPinRequest
	if err := binder.BindFormat(ctx, &body, binding.JSON); err != nil {
		errorMessages := binder.FormatValidationError(err)
		if len(errorMessages) > 0 && errorMessages["error"] != "" {
			response.JSONBadRequest(ctx)
			return
		}
		response.JSONUnprocessableEntity(ctx, errorMessages)
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

// GetDashboardInformation godoc
// @Summary		Get current user wallet dashboard
// @Description	Get balance, income, and expense information for the authenticated user.
// @Tags			Users
// @Produce		json
// @Security		ApiKeyAuth
// @Param			X-Swagger	header		string			false	"Set true when using a raw token from Swagger UI"
// @Success		200	{object}	dto.Response	"Get Dashboard Information Successfully"
// @Failure		401	{object}	dto.Response	"Unauthorized"
// @Failure		500	{object}	dto.Response	"Internal server error"
// @Router			/users/me/wallet [get]
func (uc *UserController) GetDashboardInformation(ctx *gin.Context) {
	claims, ok := jwttoken.GetClaims(ctx)
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

// GetTransactionReport godoc
// @Summary		Get current user transaction report
// @Description	Get the authenticated user's transaction report for the last 7 days.
// @Tags			Users
// @Produce		json
// @Security		ApiKeyAuth
// @Param			X-Swagger	header		string			false	"Set true when using a raw token from Swagger UI"
// @Param			duration	query		string			false	"Report duration"	Enums(7d)					default(7d)
// @Param			type		query		string			false	"Report type"		Enums(all, income, expense)	default(all)
// @Success		200			{object}	dto.Response	"Get Transaction Report Successfully"
// @Failure		400			{object}	dto.Response	"Bad request"
// @Failure		401			{object}	dto.Response	"Unauthorized"
// @Failure		500			{object}	dto.Response	"Internal server error"
// @Router			/users/me/transaction-report [get]
func (uc *UserController) GetTransactionReport(ctx *gin.Context) {
	var query dto.TransactionReportQueryRequest
	if err := binder.BindFormat(ctx, &query, binding.Query); err != nil {
		response.JSONBadRequest(ctx)
		return
	}

	reportType := query.Type
	if reportType == "" {
		reportType = "all"
	}

	claims, ok := jwttoken.GetClaims(ctx)
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
