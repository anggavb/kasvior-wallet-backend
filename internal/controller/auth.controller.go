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

// Register godoc
// @Summary		Register new user
// @Description	Create a new user account with email and password.
// @Tags		Auth
// @Accept		json
// @Produce		json
// @Param		request	body		dto.AuthRequest	true	"Register request body"
// @Success		201		{object}	dto.Response		"Register Successfully"
// @Failure		400		{object}	dto.Response		"Bad request"
// @Failure		409		{object}	dto.Response		"Email Already Used"
// @Failure		422		{object}	dto.Response		"Validation error"
// @Failure		500		{object}	dto.Response		"Internal server error"
// @Router			/auth/register [post]
func (ac *AuthController) Register(ctx *gin.Context) {
	var body dto.AuthRequest
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
// @Param			request	body		dto.AuthRequest	true	"Login request body"
// @Success		200		{object}	dto.Response		"Login Successfully"
// @Failure		400		{object}	dto.Response		"Bad request"
// @Failure		401		{object}	dto.Response		"Invalid email or password"
// @Failure		422		{object}	dto.Response		"Validation error"
// @Failure		500		{object}	dto.Response		"Internal server error"
// @Router			/auth [post]
func (ac *AuthController) Login(ctx *gin.Context) {
	var body dto.AuthRequest
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

// Logout godoc
// @Summary		Logout user
// @Description	Invalidate the current bearer token.
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Security		ApiKeyAuth
// @Success		200	{object}	dto.Response	"Logout Successfully"
// @Failure		401	{object}	dto.Response	"Unauthorized"
// @Failure		500	{object}	dto.Response	"Internal server error"
// @Router			/auth/logout [post]
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
