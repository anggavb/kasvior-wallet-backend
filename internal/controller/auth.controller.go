package controller

import (
	"errors"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/kasvior-wallet-backend/internal/apperrors"
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

// Register godoc
// @Summary		Register new user
// @Description	Create a new user account with email and password.
// @Tags		Auth
// @Accept		json
// @Produce		json
// @Param		request	body		dto.RegisterRequest	true	"Register request body"
// @Success		201		{object}	dto.Response		"Register Successfully"
// @Failure		400		{object}	dto.Response		"Bad request"
// @Failure		409		{object}	dto.Response		"Email Already Used"
// @Failure		422		{object}	dto.Response		"Validation error"
// @Failure		500		{object}	dto.Response		"Internal server error"
// @Router			/auth/register [post]
func (ac *AuthController) Register(ctx *gin.Context) {
	var body dto.RegisterRequest
	if err := binder.BindFormat(ctx, &body, binding.JSON); err != nil {
		errorMessages := binder.FormatValidationError(err)
		if len(errorMessages) > 0 && errorMessages["error"] != "" {
			response.JSONBadRequest(ctx)
			return
		}
		response.JSONUnprocessableEntity(ctx, errorMessages)
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

// Login godoc
// @Summary		Login user
// @Description	Authenticate a user and return an access token.
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			request	body		dto.LoginRequest	true	"Login request body"
// @Success		200		{object}	dto.Response		"Login Successfully"
// @Failure		400		{object}	dto.Response		"Bad request"
// @Failure		401		{object}	dto.Response		"Invalid email or password"
// @Failure		422		{object}	dto.Response		"Validation error"
// @Failure		500		{object}	dto.Response		"Internal server error"
// @Router			/auth [post]
func (ac *AuthController) Login(ctx *gin.Context) {
	var body dto.LoginRequest
	if err := binder.BindFormat(ctx, &body, binding.JSON); err != nil {
		errorMessages := binder.FormatValidationError(err)
		if len(errorMessages) > 0 && errorMessages["error"] != "" {
			response.JSONBadRequest(ctx)
			return
		}
		response.JSONUnprocessableEntity(ctx, errorMessages)
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

// ForgotPassword godoc
// @Summary		Request password reset
// @Description	Send password reset instructions when the email is registered.
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			request	body		dto.ForgotPasswordRequest	true	"Forgot password request body"
// @Success		200		{object}	dto.Response				"If the email is registered, reset instructions have been sent"
// @Failure		400		{object}	dto.Response				"Bad request"
// @Failure		422		{object}	dto.Response				"Validation error"
// @Failure		500		{object}	dto.Response				"Internal server error"
// @Router			/auth/forgot-password [post]
func (ac *AuthController) ForgotPassword(ctx *gin.Context) {
	var body dto.ForgotPasswordRequest
	if err := binder.BindFormat(ctx, &body, binding.JSON); err != nil {
		errorMessages := binder.FormatValidationError(err)
		if len(errorMessages) > 0 && errorMessages["error"] != "" {
			response.JSONBadRequest(ctx)
			return
		}
		response.JSONUnprocessableEntity(ctx, errorMessages)
		return
	}

	if err := ac.authService.ForgotPassword(ctx.Request.Context(), body); err != nil {
		log.Println("Error: ", err.Error())
		response.JSONInternalServerError(ctx)
		return
	}

	response.JSONSuccess(ctx, nil, "Reset password have been sent")
}

// ResetPassword godoc
// @Summary		Reset password
// @Description	Reset a user password using a valid password reset token.
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			request	body		dto.ResetPasswordRequest	true	"Reset password request body"
// @Success		200		{object}	dto.Response				"Reset Password Successfully"
// @Failure		400		{object}	dto.Response				"Bad request"
// @Failure		401		{object}	dto.Response				"Invalid or expired reset token"
// @Failure		422		{object}	dto.Response				"Validation error"
// @Failure		500		{object}	dto.Response				"Internal server error"
// @Router			/auth/reset-password [post]
func (ac *AuthController) ResetPassword(ctx *gin.Context) {
	var body dto.ResetPasswordRequest
	if err := binder.BindFormat(ctx, &body, binding.JSON); err != nil {
		errorMessages := binder.FormatValidationError(err)
		if len(errorMessages) > 0 && errorMessages["error"] != "" {
			response.JSONBadRequest(ctx)
			return
		}
		response.JSONUnprocessableEntity(ctx, errorMessages)
		return
	}

	if err := ac.authService.ResetPassword(ctx.Request.Context(), body); err != nil {
		log.Println("Error: ", err.Error())
		if errors.Is(err, apperrors.ErrInvalidPasswordResetToken) {
			response.JSONUnauthorized(ctx, "Invalid or expired reset token")
			return
		}
		response.JSONInternalServerError(ctx)
		return
	}

	response.JSONSuccess(ctx, nil, "Reset Password Successfully")
}

// Logout godoc
// @Summary		Logout user
// @Description	Invalidate the current bearer token.
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Security		ApiKeyAuth
// @Param			X-Swagger	header		string			false	"Set true when using a raw token from Swagger UI"
// @Success		200	{object}	dto.Response	"Logout Successfully"
// @Failure		401	{object}	dto.Response	"Unauthorized"
// @Failure		500	{object}	dto.Response	"Internal server error"
// @Router			/auth/logout [delete]
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
		if errors.Is(err, apperrors.ErrTokenAlreadyExpired) {
			response.JSONUnauthorized(ctx, "Token already expired")
			return
		}
		response.JSONInternalServerError(ctx)
		return
	}

	response.JSONNoContent(ctx)
}
