package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/kasvior-wallet-backend/internal/dto"
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
	if err := ctx.ShouldBindWith(&body, binding.JSON); err != nil {
		log.Println("Error", err.Error())
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Message: "Invalid Request Payload",
			Error:   "Bad Request",
		})
		return
	}

	res, err := ac.authService.RegisterUser(ctx.Request.Context(), body)
	if err != nil {
		log.Println("Error: ", err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Message: "Error",
			Error:   "Internal Server Error",
		})
		return
	}
	ctx.JSON(http.StatusCreated, dto.Response{
		Data:    res,
		Message: "Register Successfully",
		Success: true,
	})
}
